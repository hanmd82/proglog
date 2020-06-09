package discovery

import (
	"log"
	"net"

	"github.com/hashicorp/serf/serf"
)

// Membership wraps Serf to provide discovery and cluster membership functionality.
type Membership struct {
	Config
	handler Handler
	serf    *serf.Serf
	events  chan serf.Event
}

// New() creates a Membership with the required configuration and event handler.
func New(handler Handler, config Config) (*Membership, error) {
	c := &Membership{
		Config:  config,
		handler: handler,
	}

	if err := c.setupSerf(); err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	NodeName       string
	BindAddr       *net.TCPAddr
	Tags           map[string]string
	StartJoinAddrs []string
}

// setupSerf() creates and configures a Serf instance
// and starts the eventsHandler() goroutine to handle Serf’s events.
func (s *Membership) setupSerf() (err error) {
	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = s.BindAddr.IP.String()
	config.MemberlistConfig.BindPort = s.BindAddr.Port

	s.events = make(chan serf.Event)
	config.EventCh = s.events
	config.Tags = s.Tags
	config.NodeName = s.Config.NodeName

	s.serf, err = serf.Create(config)
	if err != nil {
		return err
	}

	go s.eventHandler()
	if s.StartJoinAddrs != nil {
		_, err = s.serf.Join(s.StartJoinAddrs, true)
		if err != nil {
			return err
		}
	}
	return nil
}

type Handler interface {
	Join(name, addr string) error
	Leave(name, addr string) error
}

func (s *Membership) eventHandler() {
	for e := range s.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, m := range e.(serf.MemberEvent).Members {
				if s.isLocal(m) {
					continue
				}
				s.handleJoin(m)
			}
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, m := range e.(serf.MemberEvent).Members {
				if s.isLocal(m) {
					return
				}
				s.handleLeave(m)
			}
		}
	}
}

func (s *Membership) handleJoin(m serf.Member) {
	if err := s.handler.Join(m.Name, m.Tags["rpc_addr"]); err != nil {
		log.Printf("[ERROR] proglog: failed to join: %s %s", m.Name, m.Tags["rpc_addr"])
	}
}

func (s *Membership) handleLeave(m serf.Member) {
	if err := s.handler.Leave(m.Name, m.Tags["rpc_addr"]); err != nil {
		log.Printf("[ERROR] proglog: failed to leave: %s", m.Name)
	}
}

// isLocal() checks whether the given Serf member is the local member.
func (s *Membership) isLocal(m serf.Member) bool {
	return s.serf.LocalMember().Name == m.Name
}

// Members() returns a point-in-time snapshot of the cluster's Serf members.
func (s *Membership) Members() []serf.Member {
	return s.serf.Members()
}

// Leave() tells this member to leave the Serf cluster.
func (s *Membership) Leave() error {
	return s.serf.Leave()
}

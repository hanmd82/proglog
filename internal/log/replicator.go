package log

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"

	api "github.com/hanmd82/proglog/api/v1"
)

// Replicator connects to other servers with the gRPC client.
// Configure the client with DialOptions, to authenticate with the servers.
type Replicator struct {
	DialOptions []grpc.DialOption
	LocalServer api.LogClient

	mu      sync.Mutex
	servers map[string]chan struct{}
	closed  bool
	close   chan struct{}
}

// Join(string, string) adds the given server address to the list of servers to replicate,
// and kicks off the 'replicate' goroutine to run the actual replication logic.
func (r *Replicator) Join(name, addr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.init()

	if r.closed {
		return nil
	}

	if _, ok := r.servers[addr]; ok {
		// replication in progress
		return nil
	}
	r.servers[addr] = make(chan struct{})

	go r.replicate(addr, r.servers[addr])

	return nil
}

// replicate(string, chan struct{}) replicates from server at given 'addr',
// and listens on the 'leave' signalling channel to know when to stop replication
func (r *Replicator) replicate(addr string, leave chan struct{}) {
	cc, err := grpc.Dial(addr, r.DialOptions...)
	if err != nil {
		r.err(err)
		return
	}
	defer cc.Close()

	// create a client and open up a stream to consume all logs on the server.
	client := api.NewLogClient(cc)
	ctx := context.Background()
	stream, err := client.ConsumeStream(
		ctx,
		&api.ConsumeRequest{
			Offset: 0,
		},
	)
	if err != nil {
		r.err(err)
		return
	}

	records := make(chan *api.Record)
	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				r.err(err)
				return
			}
			records <- recv.Record
		}
	}()

	// consume logs from discovered server in a stream,
	// and produce to the local server to save a copy.
	for {
		select {
		case <-r.close:
			return
		case <-leave:
			return
		case record := <-records:
			_, err = r.LocalServer.Produce(
				ctx,
				&api.ProduceRequest{
					Record: record,
				},
			)
			if err != nil {
				r.err(err)
				return
			}
		}
	}
}

// Leave(string, string) handles the server leaving the cluster by removing the server
// from the list of servers to replicate and closes the server’s associated channel
func (r *Replicator) Leave(name, addr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.init()

	if _, ok := r.servers[addr]; !ok {
		return nil
	}

	close(r.servers[addr])
	delete(r.servers, addr)
	return nil
}

func (r *Replicator) init() {
	if r.servers == nil {
		r.servers = make(map[string]chan struct{})
	}
	if r.close == nil {
		r.close = make(chan struct{})
	}
}

func (r *Replicator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.init()

	if r.closed {
		return nil
	}

	r.closed = true
	close(r.close)
	return nil
}

func (r *Replicator) err(err error) {
	log.Printf("[ERROR] proglog: %v", err)
}

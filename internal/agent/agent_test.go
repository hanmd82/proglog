package agent_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/hanmd82/proglog/api/v1"
	"github.com/hanmd82/proglog/internal/agent"
	"github.com/hanmd82/proglog/internal/config"
)

func TestAgent(t *testing.T) {
	// serverTLSConfig defines the configuration of the certificate that is served to clients.
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	// peerTLSConfig defines the configuration of the certificate that is served between servers in the cluster.
	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	var agents []*agent.Agent
	for i := 0; i < 3; i++ {
		ports := dynaport.Get(2)
		bindAddr := &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: ports[0],
		}
		rpcPort := ports[1]

		dataDir, err := ioutil.TempDir("", "server-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(
				startJoinAddrs,
				agents[0].Config.BindAddr.String(),
			)
		}

		agent, err := agent.New(agent.Config{
			NodeName:        fmt.Sprintf("%d", i),
			StartJoinAddrs:  startJoinAddrs,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			DataDir:         dataDir,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
		})
		require.NoError(t, err)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir))
		}
	}()

	time.Sleep(3 * time.Second)

	// test produce to and consume from a single node
	leaderClient := client(t, agents[0], peerTLSConfig)
	produceResponse, err := leaderClient.Produce(
		context.Background(),
		&api.ProduceRequest{
			Record: &api.Record{
				Value: []byte("foo"),
			},
		},
	)
	require.NoError(t, err)

	consumeResponse, err := leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Value, []byte("foo"))

	// wait until replication has finished
	time.Sleep(3 * time.Second)

	// test consume from another node which has replicated the record
	followerClient := client(t, agents[1], peerTLSConfig)
	consumeResponse, err = followerClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Value, []byte("foo"))
}

func client(t *testing.T, agent *agent.Agent, tlsConfig *tls.Config) api.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCreds),
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", agent.Config.BindAddr.IP.String(), agent.Config.RPCPort),
		opts...,
	)
	require.NoError(t, err)

	client := api.NewLogClient(conn)
	return client
}

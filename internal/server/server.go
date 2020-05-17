package server

import (
	"context"

	api "github.com/hanmd82/proglog/api/v1"
	"google.golang.org/grpc"
)

// Ensure that the *grpcServer type satisfies the api.LogServer interface
var _ api.LogServer = (*grpcServer)(nil)

type Config struct {
	CommitLog CommitLog
}

type grpcServer struct {
	*Config
}

func newgrpcServer(config *Config) (srv *grpcServer, err error) {
	srv = &grpcServer{
		Config: config,
	}
	return srv, nil
}

// NewGRPCServer(*Config, ...grpc.ServerOption) instantiates the Log service
// creates a gRPC server with given gRPC server options,
// and registers the Log service to the gRPC server
func NewGRPCServer(config *Config, opts ...grpc.ServerOption) (*grpc.Server, error) {
	gsrv := grpc.NewServer(opts...)
	srv, err := newgrpcServer(config)
	if err != nil {
		return nil, err
	}
	api.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}

// // Implement the LogServer interface
// type LogServer interface {
// 	Produce(context.Context, *ProduceRequest) (*ProduceResponse, error)
// 	Consume(context.Context, *ConsumeRequest) (*ConsumeResponse, error)
// 	ConsumeStream(*ConsumeRequest, Log_ConsumeStreamServer) error
// 	ProduceStream(Log_ProduceStreamServer) error
// }

func (s *grpcServer) Produce(
	ctx context.Context,
	req *api.ProduceRequest,
) (*api.ProduceResponse, error) {
	offset, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &api.ProduceResponse{Offset: offset}, nil
}

func (s *grpcServer) Consume(
	ctx context.Context,
	req *api.ConsumeRequest,
) (*api.ConsumeResponse, error) {
	record, err := s.CommitLog.Read(req.Offset)
	if err != nil {
		return nil, err
	}
	return &api.ConsumeResponse{Record: record}, nil
}

// ProduceStream(api.Log_ProduceStreamServer) implements bidirectional streaming RPC
func (s *grpcServer) ProduceStream(
	stream api.Log_ProduceStreamServer,
) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		res, err := s.Produce(stream.Context(), req)
		if err != nil {
			return err
		}

		if err = stream.Send(res); err != nil {
			return err
		}
	}
}

// ConsumeStream(*api.ConsumeRequest, api.Log_ConsumeStreamServer) implements server-side streaming RPC
func (s *grpcServer) ConsumeStream(
	req *api.ConsumeRequest,
	stream api.Log_ConsumeStreamServer,
) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			res, err := s.Consume(stream.Context(), req)
			switch err.(type) {
			case nil:
			case api.ErrOffsetOutOfRange:
				continue
			default:
				return err
			}
			if err = stream.Send(res); err != nil {
				return err
			}
			req.Offset++
		}
	}
}

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}

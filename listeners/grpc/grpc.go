package grpc

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	name    string // Name of the service
	server  *grpc.Server
	network string
	address string
}

func NewGRPCServer(interceptStreaming func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error, interceptUnary func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)) (*GRPCServer, error) {
	var grpcs *grpc.Server

	if interceptStreaming != nil && interceptUnary != nil {
		grpcs = grpc.NewServer(
			grpc.StreamInterceptor(interceptStreaming),
			grpc.UnaryInterceptor(interceptUnary),
		)

		return &GRPCServer{server: grpcs}, nil
	}

	if interceptStreaming != nil {
		grpcs = grpc.NewServer(
			grpc.StreamInterceptor(interceptStreaming),
		)

		return &GRPCServer{server: grpcs}, nil
	}

	if interceptUnary != nil {
		grpcs = grpc.NewServer(
			grpc.UnaryInterceptor(interceptUnary),
		)

		return &GRPCServer{server: grpcs}, nil
	}

	return nil, errors.New("No streaming or unary function provided")
}

func (grpc *GRPCServer) Run() {
	l, err := net.Listen(grpc.network, grpc.address)
	if err != nil {
		panic("")
	}

	grpc.server.Serve(l)
}

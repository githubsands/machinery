package grpc

import (
	"context"
	"errors"
	"net"

	"github.com/githubsands/machinery/listeners/chat"
	"google.golang.org/grpc"
)

// interceptStreaming is used when needing a polling, i.e notifications, listening for events, and sending data followed by deltas
type interceptStreaming func(interface{}, grpc.ServerStream, *grpc.StreamServerInfo, grpc.StreamHandler) error

// inteceptUnary is generally used when expecting to send a single event back
type interceptUnary func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error)

type GRPCServer struct {
	name    string // Name of the service
	server  *grpc.Server
	network string
	address string

	chats chan<- chat.ChatMsg
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

func (grpc *GRPCServer) AddChat(chats chan<- chat.ChatMsg) {
	grpc.chats = chats
}

func (grpc *GRPCServer) SendToChat(msg chat.ChatMsg) {
	grpc.chats <- msg
}

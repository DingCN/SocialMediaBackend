package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"github.com/DingCN/SocialMediaBackend/pkg/backend"
	pb "github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	backend, _ := backend.New()
	s := grpc.NewServer()
	pb.RegisterTwitterRPCServer(s, backend)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

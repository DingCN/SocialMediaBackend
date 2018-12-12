package main

import (
	"log"
	"net"

	"github.com/DingCN/SocialMediaBackend/pkg/backendraft"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// pb "google.golang.org/grpc/examples/helloworld/helloworld"

	pb "github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

// const (
// 	port = ":50051"
// )

// Start backend server
func main() {
	backend, _ := backendraft.New()

	lis, err := net.Listen("tcp", backend.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTwitterRPCServer(s, backend)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

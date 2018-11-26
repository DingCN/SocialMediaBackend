package main

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/DingCN/SocialMediaBackend/pkg/backend"
	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
	"github.com/DingCN/SocialMediaBackend/pkg/web"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

// pb "google.golang.org/grpc/examples/helloworld/helloworld"

var webSrv = &web.Web{}

func startBackend() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	backend, _ := backend.New()
	s := grpc.NewServer()
	protocol.RegisterTwitterRPCServer(s, backend)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func startWeb() {

	backendAddr := "localhost:50051"
	conn, err := grpc.Dial(backendAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("web adn backend did not connect: %v", err)

	}
	webSrv.C = protocol.NewTwitterRPCClient(conn)

}
func TestStartServer(t *testing.T) {
	startWeb()
	// Starting Backend
	go startBackend()
}
func TestSignupRPC(t *testing.T) {
	//starting Web

	//calling RPC
	_, err := webSrv.SignupRPCSend("user1", "password1")
	if err != nil {
		t.Fatalf("Signup Incorrect")
	}
	_, err = webSrv.SignupRPCSend("user1", "password1")
	if err.Error() != errorcode.ErrUsernameTaken.Error() {
		t.Fatalf("Signup Incorrect")
	}
}

func TestMomentsRPC(t *testing.T) {

	// env set
	_, err := webSrv.SignupRPCSend("Test_MomentsAlice", "Test_MomentsAlice")
	if err != nil {
		t.Fatalf("Moments Incorrect")
	}
	_, err = webSrv.SignupRPCSend("Test_MomentsBob", "Test_MomentsBob")
	if err != nil {
		t.Fatalf("Moments Incorrect")
	}
	_, err = webSrv.SignupRPCSend("Test_MomentsCain", "Test_MomentsCain")
	if err != nil {
		t.Fatalf("Moments Incorrect")
	}
	_, err = webSrv.AddTweetRPCSend("Test_MomentsBob", "Test_MomentsBob's post")
	if err != nil {
		t.Fatalf("Moments Incorrect")
	}
	_, err = webSrv.AddTweetRPCSend("Test_MomentsCain", "Test_MomentsCain's post")
	if err != nil {
		t.Fatalf("Moments Incorrect")
	}
	reply, err := webSrv.MomentRandomFeedsRPCSend()
	if err != nil {
		t.Fatalf("Moments Incorrect")
	}
	tweets := reply.TweetList
	if len(tweets) != 2 || tweets[1].UserName != "Test_MomentsBob" || tweets[1].Body != "Test_MomentsBob's post" || tweets[0].UserName != "Test_MomentsCain" || tweets[0].Body != "Test_MomentsCain's post" {
		t.Fatalf("Moments incorrect")
		fmt.Printf("%+v\n", tweets)
	}
}

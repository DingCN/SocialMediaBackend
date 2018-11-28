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

// run whole file to test, testing single test case won't work since we cannot start server in every test function

var webSrv = &web.Web{}

func startBackend() {
	lis, err := net.Listen("tcp", ":50051") // TODO add to config
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
	// Starting Web
	startWeb()
	// Starting Backend
	go startBackend()
}

// TestSignupRPC tests a successful signup, and then signup again to see if "Username taken error" is triggered correctly
func TestSignupRPC(t *testing.T) {
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

// TestMomentsRPC tests the moments feature
// Moments:
// When a user registers, he isn't following any other users.
// We provide a moment page so that it can get the newest posts even he is not following their owner
func TestMomentsRPC(t *testing.T) {
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

// TestGetFollowingTweetsRPC tests if a user only gets tweets from those users he/she follows
func TestGetFollowingTweetsRPC(t *testing.T) {
	webSrv.SignupRPCSend("TestGetFollowingTweetsRPC_Alice", "TestGetFollowingTweetsRPC_Alice")
	webSrv.SignupRPCSend("TestGetFollowingTweetsRPC_Bob", "TestGetFollowingTweetsRPC_Bob")
	webSrv.SignupRPCSend("TestGetFollowingTweetsRPC_Cain", "TestGetFollowingTweetsRPC_Cain")
	webSrv.SignupRPCSend("TestGetFollowingTweetsRPC_Doge", "TestGetFollowingTweetsRPC_Doge")
	webSrv.FollowUnFollowRPCSend("TestGetFollowingTweetsRPC_Alice", "TestGetFollowingTweetsRPC_Bob")
	webSrv.FollowUnFollowRPCSend("TestGetFollowingTweetsRPC_Alice", "TestGetFollowingTweetsRPC_Cain")
	// Alice is following Bob
	_, err := webSrv.AddTweetRPCSend("TestGetFollowingTweetsRPC_Bob", "TestGetFollowingTweetsRPC_Bob's post")
	if err != nil {
		t.Fatalf("AddTweetRPC Incorrect")
	}
	_, err = webSrv.AddTweetRPCSend("TestGetFollowingTweetsRPC_Cain", "TestGetFollowingTweetsRPC_Cain's post")
	if err != nil {
		t.Fatalf("AddTweetRPC Incorrect")
	}
	_, err = webSrv.AddTweetRPCSend("TestGetFollowingTweetsRPC_Bob", "TestGetFollowingTweetsRPC_Bob's post2")
	if err != nil {
		t.Fatalf("AddTweetRPC Incorrect")
	}
	_, err = webSrv.AddTweetRPCSend("TestGetFollowingTweetsRPC_Doge", "TestGetFollowingTweetsRPC_Doge's post")
	if err != nil {
		t.Fatalf("AddTweetRPC Incorrect")
	}

	//module test
	reply, err := webSrv.GetFollowingTweetsRPCSend("TestGetFollowingTweetsRPC_Alice")
	if err != nil {
		t.Fatalf("TestGetFollowingTweetsRPC incorrect")
	}
	tweets := reply.Tweet
	if len(tweets) != 3 || tweets[2].UserName != "TestGetFollowingTweetsRPC_Bob" || tweets[2].Body != "TestGetFollowingTweetsRPC_Bob's post" || tweets[1].UserName != "TestGetFollowingTweetsRPC_Cain" || tweets[1].Body != "TestGetFollowingTweetsRPC_Cain's post" || tweets[0].UserName != "TestGetFollowingTweetsRPC_Bob" || tweets[0].Body != "TestGetFollowingTweetsRPC_Bob's post2" {
		t.Fatalf("TestGetFollowingTweetsRPC incorrect")
	}
}

// TestUserProfileRPC tests the userprofile page, in which number of following, number of follower, feeds is displayed to user
func TestUserProfileRPC(t *testing.T) {
	webSrv.SignupRPCSend("TestUserProfileRPCAlice", "TestUserProfileRPCAlice")
	webSrv.SignupRPCSend("TestUserProfileRPC_Bob", "TestUserProfileRPC_Bob")
	webSrv.SignupRPCSend("TestUserProfileRPC_Cain", "TestUserProfileRPC_Cain")
	webSrv.SignupRPCSend("TestUserProfileRPC_Doge", "TestUserProfileRPC_Doge")
	webSrv.FollowUnFollowRPCSend("TestUserProfileRPC_Alice", "TestUserProfileRPC_Bob")
	webSrv.FollowUnFollowRPCSend("TestUserProfileRPC_Alice", "TestUserProfileRPC_Cain")
	// Alice is following Bob
	webSrv.AddTweetRPCSend("TestUserProfileRPC_Bob", "TestUserProfileRPC_Bob's post")

	username := "TestUserProfileRPC_Bob"
	// Query()["key"] will return an array of items,
	// we only want the single item.
	reply, err := webSrv.GetUserProfileRPCSend(username)
	if err != nil {
		t.Fatalf("TestUserProfileRPC incorrect")
	}
	if reply.Username != "TestUserProfileRPC_Bob" || len(reply.TweetList) != 1 || reply.TweetList[0].Body != "TestUserProfileRPC_Bob's post" {
		t.Fatalf("UserProfile incorrect")
	}
}

package main

import (
	"testing"

	"github.com/DingCN/SocialMediaBackend/pkg/backend"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
)

// pb "google.golang.org/grpc/examples/helloworld/helloworld"

func TestAddUser(t *testing.T) {
	backend, _ := backend.New()
	success, err := backend.Storage.AddUser("", "12345678")
	if success != false || err.Error() != errorcode.ErrInvalidUsername.Error() {
		t.Fatalf("AddUser incorrect")
	}
	success, err = backend.Storage.AddUser("user1", "")
	if success != false || err.Error() != errorcode.ErrInvalidPassword.Error() {
		t.Fatalf("AddUser incorrect")
	}
	success, err = backend.Storage.AddUser("user1", "12345678")
	if success != true || err != nil {
		t.Fatalf("AddUser incorrect")
	}
	success, err = backend.Storage.AddUser("user1", "12345678")
	if success != false || err.Error() != errorcode.ErrUsernameTaken.Error() {
		t.Fatalf("AddUser incorrect")
	}
}

func TestGetUser(t *testing.T) {
	backend, _ := backend.New()
	success, err := backend.Storage.AddUser("user1", "12345678")
	if success != true || err != nil {
		t.Fatalf("GetUser incorrect")
	}
	pUser, err := backend.Storage.GetUser("user1")
	if pUser == nil || err != nil {
		t.Fatalf("GetUser incorrect")
	}
	pUser, err = backend.Storage.GetUser("aaaaaaaaa")
	if pUser != nil || err.Error() != errorcode.ErrUserNotExist.Error() {
		t.Fatalf("GetUser incorrect")
	}
}

func TestAddTweet(t *testing.T) {
	backend, _ := backend.New()
	success, err := backend.Storage.AddUser("user1", "12345678")
	success, err = backend.Storage.AddTweet("user1", "tweet1")
	if success != true || err != nil {
		t.Fatalf("TestAddTweet incorrect")
	}
	pUser, err := backend.Storage.GetUser("user1")
	if pUser == nil || err != nil {
		t.Fatalf("TestAddTweet incorrect")
	}
	if pUser.TweetList[0].Body != "tweet1" || pUser.TweetList[0].UserName != "user1" {
		t.Fatalf("TestAddTweet incorrect")
	}
}

func TestGetTweetByUsername(t *testing.T) {
	backend, _ := backend.New()
	backend.Storage.AddUser("user1", "12345678")
	backend.Storage.AddTweet("user1", "tweet1")
	tweets, err := backend.Storage.GetTweetByUsername("user1")
	if err != nil {
		t.Fatalf("TestGetTweetByUsername incorrect")
	}
	if tweets[0].Body != "tweet1" || tweets[0].UserName != "user1" {
		t.Fatalf("TestGetTweetByUsername incorrect")
	}
}

func TestGetRandomTweet(t *testing.T) {
	backend, _ := backend.New()
	backend.Storage.AddUser("user1", "12345678")
	backend.Storage.AddTweet("user1", "tweet1")
	tweets, _ := backend.Storage.GetRandomTweet()
	if tweets[0].UserName != "user1" || tweets[0].Body != "tweet1" {
		t.Fatalf("TestGetTweetByUsername incorrect")
	}
}
func TestFollowUnFollow(t *testing.T) {

}
func TestCheckIfFollowing(t *testing.T) {

}
func TestGetAllFollowing(t *testing.T) {

}
func TestGetFollowingTweets(t *testing.T) {

}

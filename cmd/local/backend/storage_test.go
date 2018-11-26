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
	backend, _ := backend.New()
	backend.Storage.AddUser("user1", "123456")
	backend.Storage.AddUser("user2", "12345678")
	backend.Storage.AddUser("user3", "12345678")
	backend.Storage.FollowUnFollow("user1", "user2")
	backend.Storage.FollowUnFollow("user1", "user3")
	// Test1 - normal case
	following, err := backend.Storage.CheckIfFollowing("user1", "user2")
	if !following || err != nil {
		t.Fatalf("TestFollowUnFollow - test1 incorrect")
	}
	backend.Storage.FollowUnFollow("user1", "user3")
	following, err = backend.Storage.CheckIfFollowing("user1", "user3")
	if following || err != nil {
		t.Fatalf("TestFollowUnFollow - test1 incorrect")
	}

	// Test2 - Follow non-existing user
	following, err = backend.Storage.CheckIfFollowing("user1", "nilUser")
	if following == true {
		t.Fatalf("Test FollowUnFollow - test2 wrong return value")
	}
	if err == nil {
		t.Fatalf("Test FollowUnFollow - test2 error type incorrect")
	}
}
func TestCheckIfFollowing(t *testing.T) {
	backend, _ := backend.New()
	backend.Storage.AddUser("user1", "123456")
	backend.Storage.AddUser("user2", "1234567")
	backend.Storage.FollowUnFollow("user1", "user2")
	// Test1 - normal case
	following, err := backend.Storage.CheckIfFollowing("user1", "user2")
	if !following || err != nil {
		t.Fatalf("TestCheckIfFollowing - test1 incorrect")
	}
	backend.Storage.FollowUnFollow("user1", "user2")
	following, err = backend.Storage.CheckIfFollowing("user1", "user2")
	if following || err != nil {
		t.Fatalf("TestCheckIfFollowing - test1 incorrect")
	}
}
func TestGetAllFollowing(t *testing.T) {
	backend, _ := backend.New()
	backend.Storage.AddUser("user1", "123456")
	backend.Storage.AddUser("user2", "1234567")
	backend.Storage.AddUser("user3", "12345678")
	backend.Storage.FollowUnFollow("user1", "user2")
	backend.Storage.FollowUnFollow("user1", "user3")

	followingList, err := backend.Storage.GetAllFollowing("user1")
	// Test1 - normal case
	if len(followingList) != 2 || err != nil {
		t.Fatalf("TestGetAllFollowing - test1 incorrect")
	}

}
func TestGetFollowingTweets(t *testing.T) {
	backend, _ := backend.New()
	backend.Storage.AddUser("user1", "123456")
	backend.Storage.AddUser("user2", "1234567")
	backend.Storage.AddUser("user3", "12345678")
	backend.Storage.FollowUnFollow("user1", "user2")
	backend.Storage.FollowUnFollow("user1", "user3")
	backend.Storage.AddTweet("user2", "user2's tweet")
	backend.Storage.AddTweet("user3", "user3's tweet1")
	backend.Storage.AddTweet("user2", "user2's tweet2")
	backend.Storage.AddTweet("user3", "user3's tweet2")

	tweetList, _ := backend.Storage.GetFollowingTweets("user1")
	if len(tweetList) != 4 {
		t.Fatalf("TestGetFollowingTweets: incorrect number")
	}

	if tweetList[0].UserName != "user3" {
		t.Fatalf("TestGetFollowingTweets: incorrect order")
	}
}

// Package backend stores all of our state for the application
package backend

import (
	"context"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.

// Backend server
type backend struct {
	Storage storage
	//srv *http.Server
	// //client handle when comm with backend
	// c protocol.TwitterRPCClient
}

// New config
func New() (*backend, error) {
	return &backend{
		Storage: storage{
			UserList:         userlist{Users: map[string]*User{}},
			CentralTweetList: centraltweetlist{Tweets: []*Tweet{}},
		},

		// srv: &http.Server{
		// 	Addr: cfg.Addr,
		// },
	}, nil
}

// func (s *backend) Start() error {

// 	pb.RegisterTwitterRPCServer(s, &backend{})
// 	// Register reflection service on gRPC server.
// 	reflection.Register(s)
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }

func (s *backend) SignupRPC(ctx context.Context, in *protocol.SignupRequest) (*protocol.SignupReply, error) {
	username := in.GetUsername()
	password := in.GetPassword()
	ok, err := s.Storage.AddUser(username, password)

	reply := protocol.SignupReply{}
	reply.Username = username
	reply.Success = ok
	return &reply, err
}

func (s *backend) LoginRPC(ctx context.Context, in *protocol.LoginRequest) (*protocol.LoginReply, error) {
	username := in.GetUsername()
	password := in.GetPassword()
	reply := protocol.LoginReply{}
	reply.Username = username

	pUser, err := s.Storage.GetUser(username)
	if err != nil {

		reply.Success = false
		return &reply, err
	}
	if pUser.Password != password {
		err := errorcode.ErrIncorrectPassword
		reply.Success = false
		return &reply, err
	}

	reply.Success = true
	return &reply, nil
}

// OPAddTweet(username, post)
func (s *backend) AddTweetRPC(ctx context.Context, in *protocol.AddTweetRequest) (*protocol.AddTweetReply, error) {
	username := in.GetUsername()
	post := in.GetPost()
	reply := protocol.AddTweetReply{}
	reply.Username = username
	ok, err := s.Storage.AddTweet(username, post)
	reply.Success = ok
	return &reply, err
}

// OPFollowUnFollow(username, target[0])
func (s *backend) FollowUnFollowRPC(ctx context.Context, in *protocol.FollowUnFollowRequest) (*protocol.FollowUnFollowReply, error) {
	username := in.GetUsername()
	targetname := in.GetTargetname()
	reply := protocol.FollowUnFollowReply{}
	reply.Username = username
	ok, err := s.Storage.FollowUnFollow(username, targetname)
	reply.Success = ok
	return &reply, err
}

// OPGetFollowingTweets(pUser.UserName)
// func (s *backend) GetFollowingTweetsRPC(ctx context.Context, in *protocol.GetFollowingTweetsRequest) (*protocol.GetFollowingTweetsReply, error) {
// 	username := in.GetUsername()
// 	password := in.GetPassword()
// 	reply := protocol.GetFollowingTweetsReply{}
// 	reply.Username = username

// 	pUser, err := s.Storage.GetUser(username)
// 	if err != nil {

// 		reply.Success = false
// 		return &reply, err
// 	}
// 	if pUser.Password != password {
// 		err := errorcode.ErrIncorrectPassword
// 		reply.Success = false
// 		return &reply, err
// 	}

// 	reply.Success = true
// 	return &reply, nil
// }

func (s *backend) GetFollowingTweetsRPC(ctx context.Context, in *protocol.GetFollowingTweetsRequest) (*protocol.GetFollowingTweetsReply, error) {
	username := in.GetUsername()
	reply := protocol.GetFollowingTweetsReply{}
	reply.Username = username
	tweets, err := s.Storage.GetFollowingTweets(username)
	reply.Tweet, err = s.ConvertTweetListToProtoTweetList(tweets)
	reply.Success = true
	return &reply, err
}

func (s *backend) GetUserProfileRPC(ctx context.Context, in *protocol.GetUserProfileRequest) (*protocol.GetUserProfileReply, error) {
	username := in.GetUsername()
	reply := &protocol.GetUserProfileReply{}
	reply.Username = username
	pUser, err := s.Storage.GetUserProfile(username)
	reply.Username = pUser.UserName
	reply.TweetList, err = s.ConvertTweetListToProtoTweetList(pUser.TweetList)
	if err != nil {
		return nil, err
	}
	reply.FollowerList, err = s.ConvertFollowListToProtoFollowList(pUser.FollowerList)
	if err != nil {
		return nil, err
	}
	reply.FollowingList, err = s.ConvertFollowListToProtoFollowList(pUser.FollowingList)
	if err != nil {
		return nil, err
	}
	reply.Success = true
	return reply, nil

}

func (s *backend) ConvertTweetListToProtoTweetList(tweets []Tweet) ([]*protocol.Tweet, error) {
	res := []*protocol.Tweet{}
	for _, tweet := range tweets {
		stProtoTweet := protocol.Tweet{}
		stProtoTweet.UserName = tweet.UserName
		stProtoTweet.Body = tweet.Body
		stProtoTweet.Timestamp = &tweet.Timestamp
		// tweetTime := tweet.Timestamp
		// s := int64(tweetTime.Seconds())     // from 'int'
		// n := int32(tweetTime.Nanoseconds()) // from 'int'
		// ts := &timestamp.Timestamp{Seconds: s, Nanos: n}
		res = append(res, &stProtoTweet)
	}
	return res, nil
}

// Both Following and Follower list
func (s *backend) ConvertFollowListToProtoFollowList(followList map[string]bool) ([]string, error) {
	res := []string{}
	for user, _ := range followList {
		res = append(res, user)
	}
	return res, nil
}

// OPGetRandomTweet()

// Package backend stores all of our state for the application
package backendraft

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"strings"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
	"go.etcd.io/etcd/raft/raftpb"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.

// Backend server
type backend struct {
	kvStore kvstore
	//srv *http.Server
	// //client handle when comm with backend
	// c protocol.TwitterRPCClient
}

// New config
func New() (*backend, error) {
	// raft
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()
	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)
	// raft provides a commit stream for the proposals from the http api
	var kvs *kvstore
	getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }
	commitC, errorC, snapshotterReady := newRaftNode(*id, strings.Split(*cluster, ","), *join, getSnapshot, proposeC, confChangeC)
	kvs = newKVStore(<-snapshotterReady, proposeC, commitC, errorC)
	return &backend{*kvs}, nil
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
	type st struct {
		username string
		password string
	}
	byte, _ := json.Marshal(st{username, password})
	ok, err := s.kvStore.Propose(protocol.Functions_FunctionName_value["SignupRPC"], byte)
	// ok, err := s.kvStore.Propose(protocol.Functions_FunctionName_value["SignupRPC"], json.Marshal(st{username, password}))

	// ok, err := s.kvStore.kvStore.AddUser(username, password)

	reply := protocol.SignupReply{}
	reply.Username = username
	// reply.Success = true
	// return &reply, nil
	reply.Success = ok
	return &reply, err
}

func (s *backend) LoginRPC(ctx context.Context, in *protocol.LoginRequest) (*protocol.LoginReply, error) {
	username := in.GetUsername()
	password := in.GetPassword()
	reply := protocol.LoginReply{}
	reply.Username = username

	pUser, err := s.kvStore.GetUser(username)
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
	ok, err := s.kvStore.kvStore.AddTweet(username, post)
	reply.Success = ok
	return &reply, err
}

// OPFollowUnFollow(username, target[0])
func (s *backend) FollowUnFollowRPC(ctx context.Context, in *protocol.FollowUnFollowRequest) (*protocol.FollowUnFollowReply, error) {
	username := in.GetUsername()
	targetname := in.GetTargetname()
	reply := protocol.FollowUnFollowReply{}
	reply.Username = username
	//Propose(protonum, {username, targetname})
	type st struct {
		username   string
		targetname string
	}
	byte, _ := json.Marshal(st{username, targetname})
	ok, err := s.kvStore.Propose(protocol.Functions_FunctionName_value["FollowUnFollowRPC"], byte)
	// ok, err := s.kvStore.Propose(protocol.Functions_FunctionName_value["FollowUnFollowRPC"], byte)
	// ok, err := s.kvStore.kvStore.FollowUnFollow(username, targetname)
	reply.Success = ok
	return &reply, err
}

func (s *backend) GetFollowingTweetsRPC(ctx context.Context, in *protocol.GetFollowingTweetsRequest) (*protocol.GetFollowingTweetsReply, error) {
	username := in.GetUsername()
	reply := protocol.GetFollowingTweetsReply{}
	reply.Username = username
	tweets, err := s.kvStore.GetFollowingTweets(username)
	reply.Tweet, err = s.ConvertTweetListToProtoTweetList(tweets)
	reply.Success = true
	return &reply, err
}

func (s *backend) GetUserProfileRPC(ctx context.Context, in *protocol.GetUserProfileRequest) (*protocol.GetUserProfileReply, error) {
	username := in.GetUsername()
	reply := &protocol.GetUserProfileReply{}
	reply.Username = username
	pUser, err := s.kvStore.GetUserProfile(username)
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

// We store the following list and follower list of a user in map[string]bool to ensure O(1) for checking if user is following another user
// ConvertFollowListToProtoFollowList convert a map struct of follower list to a []string struct for displaying by front-end
// This function converts both Following list and Follower list
func (s *backend) ConvertFollowListToProtoFollowList(followList map[string]bool) ([]string, error) {
	res := []string{}
	for user, _ := range followList {
		res = append(res, user)
	}
	return res, nil
}

// MomentRandomFeedsRPC is used for Moments feature
func (s *backend) MomentRandomFeedsRPC(ctx context.Context, in *protocol.MomentRandomFeedsRequest) (*protocol.MomentRandomFeedsReply, error) {

	reply := &protocol.MomentRandomFeedsReply{}
	tweetlist := s.kvStore.MomentRandomFeeds()
	protoTweetList, err := s.ConvertTweetListToProtoTweetList(tweetlist)
	if err != nil {
		log.Println(err)
	}
	reply.TweetList = protoTweetList

	reply.Success = true
	return reply, nil

}

// CheckIfFollowingRPC checks if user is following a target user
func (s *backend) CheckIfFollowingRPC(ctx context.Context, in *protocol.CheckIfFollowingRequest) (*protocol.CheckIfFollowingReply, error) {
	username := in.Username
	targetname := in.Targetname
	reply := &protocol.CheckIfFollowingReply{}
	reply.IsFollowing, _ = s.kvStore.CheckIfFollowing(username, targetname)
	return reply, nil
}

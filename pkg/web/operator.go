package web

import (
	"context"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

// func RPCAddUser(ctx context.Context, in protocol.SignupRequest) (protocol.SignupReply, error) {
// 	username := in.GetUsername()
// 	password := in.GetPassword()
// 	UserList.mutex.Lock()

// 	defer UserList.mutex.Unlock()
// 	UserList.Users[username] = &User{UserName: username, Password: password, FollowingList: map[string]bool{}, FollowerList: map[string]bool{}}

// 	reply := protocol.SignupReply{}
// 	reply.Username = username
// 	reply.Success = true
// 	return reply, nil
// }

func (web *Web) SignupRPCSend(username string, password string) (*protocol.SignupReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.SignupRequest{}
	request.Username = username
	request.Password = password
	reply, err := web.C.SignupRPC(ctx, &request)
	return reply, err
}

func (web *Web) LoginRPCSend(username string, password string) (*protocol.LoginReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.LoginRequest{}
	request.Username = username
	request.Password = password
	reply, err := web.C.LoginRPC(ctx, &request)
	return reply, err
}

func (web *Web) AddTweetRPCSend(username string, post string) (*protocol.AddTweetReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.AddTweetRequest{}
	request.Username = username
	request.Post = post
	reply, err := web.C.AddTweetRPC(ctx, &request)
	return reply, err
}

func (web *Web) FollowUnFollowRPCSend(username string, targetname string) (*protocol.FollowUnFollowReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.FollowUnFollowRequest{}
	request.Username = username
	request.Targetname = targetname
	reply, err := web.C.FollowUnFollowRPC(ctx, &request)
	return reply, err
}

func (web *Web) GetFollowingTweetsRPCSend(username string) (*protocol.GetFollowingTweetsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.GetFollowingTweetsRequest{}
	request.Username = username
	reply, err := web.C.GetFollowingTweetsRPC(ctx, &request)

	return reply, err
}

func (web *Web) GetUserProfileRPCSend(username string) (*protocol.GetUserProfileReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.GetUserProfileRequest{}
	request.Username = username
	reply, err := web.C.GetUserProfileRPC(ctx, &request)

	return reply, err
}

func (web *Web) MomentRandomFeedsRPCSend() (*protocol.MomentRandomFeedsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.MomentRandomFeedsRequest{}
	reply, err := web.C.MomentRandomFeedsRPC(ctx, &request)
	return reply, err
}

func (web *Web) CheckIfFollowingRPCSend(username string, targetname string) (*protocol.CheckIfFollowingReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.CheckIfFollowingRequest{}
	request.Username = username
	request.Targetname = targetname
	reply, err := web.C.CheckIfFollowingRPC(ctx, &request)
	return reply, err
}

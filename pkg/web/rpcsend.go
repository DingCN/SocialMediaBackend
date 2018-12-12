package web

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
	"google.golang.org/grpc/status"

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
	for i, _ := range web.C {
		fmt.Println(web.C[i])
		reply, err := web.C[i].SignupRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err
	}
	return nil, errorcode.ErrRPCConnectionLost

}

func (web *Web) LoginRPCSend(username string, password string) (*protocol.LoginReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.LoginRequest{}
	request.Username = username
	request.Password = password
	for i, _ := range web.C {
		reply, err := web.C[i].LoginRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

func (web *Web) AddTweetRPCSend(username string, post string) (*protocol.AddTweetReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.AddTweetRequest{}
	request.Username = username
	request.Post = post
	for i, _ := range web.C {
		reply, err := web.C[i].AddTweetRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

func (web *Web) FollowUnFollowRPCSend(username string, targetname string) (*protocol.FollowUnFollowReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.FollowUnFollowRequest{}
	request.Username = username
	request.Targetname = targetname
	for i, _ := range web.C {
		reply, err := web.C[i].FollowUnFollowRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

func (web *Web) GetFollowingTweetsRPCSend(username string) (*protocol.GetFollowingTweetsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.GetFollowingTweetsRequest{}
	request.Username = username
	for i, _ := range web.C {
		reply, err := web.C[i].GetFollowingTweetsRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

func (web *Web) GetUserProfileRPCSend(username string) (*protocol.GetUserProfileReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.GetUserProfileRequest{}
	request.Username = username
	for i, _ := range web.C {
		reply, err := web.C[i].GetUserProfileRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

func (web *Web) MomentRandomFeedsRPCSend() (*protocol.MomentRandomFeedsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.MomentRandomFeedsRequest{}
	for i, _ := range web.C {
		reply, err := web.C[i].MomentRandomFeedsRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

func (web *Web) CheckIfFollowingRPCSend(username string, targetname string) (*protocol.CheckIfFollowingReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.CheckIfFollowingRequest{}
	request.Username = username
	request.Targetname = targetname
	for i, _ := range web.C {
		reply, err := web.C[i].CheckIfFollowingRPC(ctx, &request)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unavailable:
					continue
				}
			}
		}
		return reply, err

	}
	return nil, errorcode.ErrRPCConnectionLost
}

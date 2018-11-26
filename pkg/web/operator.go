package web

import (
	"context"
	"fmt"
	"sort"
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
func OPAddUser(username string, password string) bool {
	UserList.mutex.Lock()
	defer UserList.mutex.Unlock()
	UserList.Users[username] = &User{UserName: username, Password: password, FollowingList: map[string]bool{}, FollowerList: map[string]bool{}}
	return true
}

func OPGetUser(username string) *User {
	return UserList.Users[username]
}

func OPAddTweet(username string, post string) bool {
	CentralTweetList.mutex.Lock()
	defer CentralTweetList.mutex.Unlock()
	UserList.mutex.Lock()
	defer UserList.mutex.Unlock()
	timestamp := time.Now()
	tweet := Tweet{UserName: username, Timestamp: timestamp, Body: post}
	CentralTweetList.Tweets = append(CentralTweetList.Tweets, &tweet)
	UserList.Users[username].TweetList = append(UserList.Users[username].TweetList, tweet)
	fmt.Printf("post: %s successfully created by user:%s\n", post, username)
	return true
}

func OPGetTweetByUsername(username string) []Tweet {
	return UserList.Users[username].TweetList
}
func OPGetRandomTweet() []Tweet {
	var count int = 0
	tweets := []Tweet{}
	for i := len(CentralTweetList.Tweets) - 1; i >= 0; i-- {
		tweets = append(tweets, *CentralTweetList.Tweets[i])
		count++
		if count >= MaxFeedsNum { ///////////////////////////////////////////////////////////////////////TODO add to config
			return tweets
		}
	}
	return tweets
}
func OPGetFollowingTweets(username string) []Tweet {
	res := []Tweet{}
	followings := OPGetAllFollowing(username)
	for _, username := range followings {
		res = append(res, UserList.Users[username].TweetList...) // ... lets you pass multiple arguments to a variadic function from a slice
	}
	// log
	sortedTweets := OPSortTweets(res)
	return sortedTweets
}

// OPSortTweets - sort tweets in descending order (LIFO)
func OPSortTweets(tweets []Tweet) []Tweet {
	res := make(timeSlice, 0, len(tweets))
	for _, d := range tweets {
		res = append(res, d)
	}

	sort.Sort(res)
	return res
}

func OPGetAllFollowing(username string) []string {
	followings := UserList.Users[username].FollowingList
	returnList := []string{}
	for followingname, isFollowing := range followings {
		if isFollowing == true {
			returnList = append(returnList, followingname)
			fmt.Printf("user:%s 's following found: %s\n", username, followingname)

		}
	}
	return returnList
}

func OPFollowUnFollow(username string, targetname string) bool {
	res, ok := UserList.Users[username].FollowingList[targetname]
	if ok == true && res == true {
		//already following, set UnFollow by deleting it instead
		delete(UserList.Users[username].FollowingList, targetname)
		delete(UserList.Users[targetname].FollowerList, username)

		// UserList.Users[username].FollowingList[targetname] = false
		// UserList.Users[targetname].FollowerList[username] = false
		fmt.Printf("%s just unfollowed %s\n", username, targetname)
	} else {
		//set Follow
		UserList.Users[username].FollowingList[targetname] = true
		UserList.Users[targetname].FollowerList[username] = true
		fmt.Printf("%s just followed %s\n", username, targetname)

	}
	return true
}

func OPCheckIfFollowing(username string, targetname string) bool {
	res, ok := UserList.Users[username].FollowingList[targetname]
	if ok == true && res == true {
		return true
	}
	return false
}

func OPGetAllUsers() []string {
	var res []string
	for username, _ := range UserList.Users {
		res = append(res, username)
	}
	return res
}

func (web *Web) SignupRPCSend(username string, password string) (*protocol.SignupReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.SignupRequest{}
	request.Username = username
	request.Password = password
	reply, err := web.c.SignupRPC(ctx, &request)
	return reply, err
}

func (web *Web) LoginRPCSend(username string, password string) (*protocol.LoginReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.LoginRequest{}
	request.Username = username
	request.Password = password
	reply, err := web.c.LoginRPC(ctx, &request)
	return reply, err
}

func (web *Web) AddTweetRPCSend(username string, post string) (*protocol.AddTweetReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.AddTweetRequest{}
	request.Username = username
	request.Post = post
	reply, err := web.c.AddTweetRPC(ctx, &request)
	return reply, err
}

func (web *Web) FollowUnFollowRPCSend(username string, targetname string) (*protocol.FollowUnFollowReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.FollowUnFollowRequest{}
	request.Username = username
	request.Targetname = targetname
	reply, err := web.c.FollowUnFollowRPC(ctx, &request)
	return reply, err
}

func (web *Web) GetFollowingTweetsRPCSend(username string) (*protocol.GetFollowingTweetsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.GetFollowingTweetsRequest{}
	request.Username = username
	reply, err := web.c.GetFollowingTweetsRPC(ctx, &request)

	return reply, err
}

// convert protoTimestamp to time.Time
func Timestamp(ts *protocol.Timestamp) time.Time {
	var t time.Time
	if ts == nil {
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}
	return t
}

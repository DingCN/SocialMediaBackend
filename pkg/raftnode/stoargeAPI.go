package raftnode

import (
	"fmt"
	"sort"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
	"github.com/DingCN/SocialMediaBackend/pkg/twitterTimestamp"
)

// AddUser register a user and his password to database
func (Storage *storage) AddUser(username string, password string) (bool, error) {
	Storage.UserList.mutex.Lock()
	defer Storage.UserList.mutex.Unlock()
	// username cannot be empty
	if len(username) < 1 {
		err := errorcode.ErrInvalidUsername
		return false, err
	}
	// password length cannot be less than 6
	if len(password) < 6 {
		err := errorcode.ErrInvalidPassword
		return false, err
	}
	_, ok := Storage.UserList.Users[username]
	if ok == true {
		err := errorcode.ErrUsernameTaken
		return false, err
	}

	Storage.UserList.Users[username] = &User{UserName: username, Password: password, FollowingList: map[string]bool{}, FollowerList: map[string]bool{}}
	return true, nil
}

// GetUser takes a username and returns a User struct, containing:
// UserID           int
// UserName         string
// Password         string
// Auth             string
// TweetList        []Tweet
// FollowingList    map[string]bool
// FollowerList     map[string]bool

func (Storage *storage) GetUser(username string) (*User, error) {
	pUser, ok := Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return nil, err
	}
	return pUser, nil
}

func (Storage *storage) GetUserProfile(username string) (*User, error) {
	pUser, ok := Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return nil, err
	}
	return pUser, nil
}

func (Storage *storage) AddTweet(username string, post string) (bool, error) {
	Storage.CentralTweetList.mutex.Lock()
	defer Storage.CentralTweetList.mutex.Unlock()
	Storage.UserList.mutex.Lock()
	defer Storage.UserList.mutex.Unlock()
	go_time := time.Now()
	timestamp := *twitterTimestamp.TimestampProto(go_time)
	tweet := Tweet{UserName: username, Timestamp: timestamp, Body: post}
	Storage.CentralTweetList.Tweets = append(Storage.CentralTweetList.Tweets, &tweet)
	pUser, ok := Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return false, err
	}
	pUser.TweetList = append(pUser.TweetList, tweet)
	fmt.Printf("post: %s successfully created by user:%s\n", post, username)
	return true, nil
}

func (Storage *storage) GetTweetByUsername(username string) ([]Tweet, error) {
	pUser, ok := Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return nil, err
	}
	return pUser.TweetList, nil
}
func (Storage *storage) GetRandomTweet() ([]Tweet, error) {
	count := 0
	tweets := []Tweet{}
	for i := len(Storage.CentralTweetList.Tweets) - 1; i >= 0; i-- {
		tweets = append(tweets, *Storage.CentralTweetList.Tweets[i])
		count++
		if count >= MaxFeedsNum { ///////////////////////////////////////////////////////////////////////TODO add to config
			return tweets, nil
		}
	}
	return tweets, nil
}

// GetFollowingTweets - return an array of sorted tweets posted by users in input user's following list
func (Storage *storage) GetFollowingTweets(username string) ([]Tweet, error) {
	res := []Tweet{}

	followings, err := Storage.GetAllFollowing(username)
	if err != nil {
		return nil, err
	}
	for _, username := range followings {
		pUser, ok := Storage.UserList.Users[username]
		if ok == false {
			err := errorcode.ErrUserNotExist
			return nil, err
		}
		res = append(res, pUser.TweetList...) // ... lets you pass multiple arguments to a variadic function from a slice
	}
	// log
	sortedTweets := Storage.SortTweets(res)
	return sortedTweets, nil
}

// SortTweets - input: []Tweet, output: []Tweet sorted by
func (Storage *storage) SortTweets(tweets []Tweet) []Tweet {
	res := make(timeSlice, 0, len(tweets))
	for _, d := range tweets {
		res = append(res, d)
	}

	sort.Sort(res)
	return res
}

func (Storage *storage) GetAllFollowing(username string) ([]string, error) {
	pUser, ok := Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return nil, err
	}
	followings := pUser.FollowingList
	returnList := []string{}
	for followingname, isFollowing := range followings {
		if isFollowing == true {
			returnList = append(returnList, followingname)
			fmt.Printf("user:%s 's following found: %s\n", username, followingname)

		}
	}
	return returnList, nil
}

func (Storage *storage) FollowUnFollow(username string, targetname string) (bool, error) {
	pUser, ok := Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return false, err
	}
	res, ok := pUser.FollowingList[targetname]
	if ok == true && res == true {
		//already following, set UnFollow by deleting it instead
		delete(Storage.UserList.Users[username].FollowingList, targetname)
		delete(Storage.UserList.Users[targetname].FollowerList, username)

		// UserList.Users[username].FollowingList[targetname] = false
		// UserList.Users[targetname].FollowerList[username] = false
		fmt.Printf("%s just unfollowed %s\n", username, targetname)
	} else {
		//set Follow
		Storage.UserList.Users[username].FollowingList[targetname] = true
		Storage.UserList.Users[targetname].FollowerList[username] = true
		fmt.Printf("%s just followed %s\n", username, targetname)

	}
	return true, nil
}

func (Storage *storage) CheckIfFollowing(username string, targetname string) (bool, error) {
	pUser, ok := Storage.UserList.Users[targetname]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return false, err
	}
	pUser, ok = Storage.UserList.Users[username]
	if ok == false {
		err := errorcode.ErrUserNotExist
		return false, err
	}
	res, ok := pUser.FollowingList[targetname]
	if ok == true && res == true {
		return true, nil
	}
	return false, nil
}

func (Storage *storage) MomentRandomFeeds() []Tweet {
	var count int = 0
	tweets := []Tweet{}
	for i := len(Storage.CentralTweetList.Tweets) - 1; i >= 0; i-- {
		tweets = append(tweets, *Storage.CentralTweetList.Tweets[i])
		count++
		if count >= 20 { ///////////////////////////////////////////////////////////////////////TODO add to config
			return tweets
		}
	}
	return tweets
}

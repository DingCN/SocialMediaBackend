package web

import (
	"fmt"
	"sort"
	"time"
)

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
	//sort

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
		if count >= 3 { ///////////////////////////////////////////////////////////////////////TODO add to config
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
	return res
}

func OPSortTweets(tweets []Tweet) []Tweet {
	res := make(timeSlice, 0, len(tweets))
	for _, d := range res {
		res = append(res, d)
	}
	sort.Sort(res)
	return res
}

func OPGetAllFollowing(username string) []string {
	followings := UserList.Users[username].FollowingList
	returnList := []string{}
	for username, isFollowing := range followings {
		if isFollowing == true {
			returnList = append(returnList, username)

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

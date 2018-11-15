package web

import "sort"

// User type definition
type User struct {
	UserID        int
	UserName      string
	Password      string
	Auth          string //currently use "true/false"; TODO change to auth token later
	TweetList     []Tweet
	FollowingList map[string]bool //map username to true/false
}

// Tweet type definition
type Tweet struct {
	//userID    string
	userName  string
	timestamp string
	body      string
}

// Global Variables
// TODO: make thread safe later
// UserList - List storing all users
var UserList map[string]User

// TweetList - List storing all tweets post by all users
var CentralTweetList []Tweet

// Use temporarily to generate unique ID
var userIDCounter int

// Helper function to sort Tweets by timestamp
func getTimestamp() string {
	//TODO: generate timestamp
}

type By func(t1, t2 *Tweet) bool

func (by By) Sort(tweets []Tweet) {
	ts := &tweetSorter{
		tweets: tweets,
		by:     by,
	}
	sort.Sort(ts)
}

type tweetSorter struct {
	tweets []Tweet
	by     func(t1, t2 *Tweet) bool
}

// IsFollowing : check if an user is followed by the current user
func (u User) IsFollowing(name string) bool {
	elem, ok := u.FollowingList[name]
	if ok == true {
		return true
	}
	return false
}

// GetTimeline : return list of tweets post by all users followed by current User
// func (u *User) GetTimeline() []Tweet {
// 	var tweetListTmp []Tweet
// 	for _, user := range u.FollowingList {
// 		for _, t := range user.TweetList {
// 			tweetListTmp = append(tweetListTmp, t)
// 		}
// 	}
// 	return tweetListTmp
// }

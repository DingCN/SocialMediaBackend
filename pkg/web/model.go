package web

import (
	"sync"
	"time"
)

////////////////////////////////////////////////////////////TODO delete this

// // global variables for storing user data goes here, will be replaced with database later
// var username_password = map[string]string{}

// // Alice: [Bob, Cain]
// // Alice is following Bob and Cain
// // hope this definition is correct
// var followingList = map[string][]string{}

// //Alice: [[Post1, timestamp], [Post2, timestamp]]
// var Posts = map[string][]TimedPost{}

// type TimedPost struct {
// 	Post      string
// 	timestamp time.Time
// }

//////////////////////////////////////////////////////////////////////

// User type definition
type User struct {
	UserID        int
	UserName      string
	Password      string
	Auth          string //currently use "true/false"; TODO change to auth token later
	TweetList     []Tweet
	FollowingList map[string]bool //map username to true/false
	//FollowingList []string        // list of username
	FollowerList map[string]bool //map username to true/false
}

// Tweet type definition
type Tweet struct {
	//userID    string
	UserName  string
	Timestamp time.Time
	Body      string
}

// Global Variables
// TODO: make thread safe later
// UserList - List storing all users
// [username: User]
type userlist struct {
	Users map[string]*User
	mutex sync.Mutex
}

var UserList = userlist{Users: map[string]*User{}}

type centraltweetlist struct {
	Tweets []*Tweet
	mutex  sync.Mutex
}

// TweetList - List storing all tweets post by all users
var CentralTweetList = centraltweetlist{Tweets: []*Tweet{}}

// Use temporarily to generate unique ID
var userIDCounter int

// Helper function to sort Tweets by timestamp
// func getTimestamp() time.Time {
// 	return time.Now()
// }

// type By func(t1, t2 *Tweet) bool

// func (by By) Sort(tweets []Tweet) {
// 	ts := &tweetSorter{
// 		tweets: tweets,
// 		by:     by,
// 	}
// 	sort.Sort(ts)
// }

// type tweetSorter struct {
// 	tweets []Tweet
// 	by     func(t1, t2 *Tweet) bool
// }

/// https://stackoverflow.com/questions/23121026/sorting-by-time-time-in-golang
type timeSlice []Tweet

// Forward request for length
func (p timeSlice) Len() int {
	return len(p)
}

// Define compare
func (p timeSlice) Less(i, j int) bool {
	return p[i].Timestamp.After(p[j].Timestamp)
}

// Define swap over an array
func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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

// Render Template Helpers
// func renderTemplate(w http.ResponseWriter, tmplname string, te)
type UserTmpl struct {
	UserName     string
	NumTweets    int
	NumFollowing int
	NumFollowers int
	TweetList    []Tweet
}

type UserListTmpl struct {
	AlreadyFollowed bool
	Following       bool
	UserName        string
	UserList        map[string]bool
}

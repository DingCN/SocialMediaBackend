package raftnode

import (
	"sync"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

var MaxFeedsNum int = 20

// var UserList = userlist{Users: map[string]*User{}}
// var CentralTweetList = centraltweetlist{Tweets: []*Tweet{}}

type storage struct {
	UserList         userlist
	CentralTweetList centraltweetlist
}

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
	Timestamp protocol.Timestamp
	Body      string
}

///////////////////////////////////////////////////////////////////////
//////////////// global variables//////////////////////////////////////
///////////////////////////////////////////////////////////////////////

// UserList - List storing all users
// [username: User]
type userlist struct {
	Users map[string]*User
	mutex sync.Mutex
}

type centraltweetlist struct {
	Tweets []*Tweet
	mutex  sync.Mutex
}

// TweetList - List storing all tweets post by all users

// Use temporarily to generate unique ID
var userIDCounter int

// sort time
/// https://stackoverflow.com/questions/23121026/sorting-by-time-time-in-golang
type timeSlice []Tweet

// Forward request for length
func (p timeSlice) Len() int {
	return len(p)
}

// Define compare
func (p timeSlice) Less(i, j int) bool {
	if p[i].Timestamp.Seconds > p[j].Timestamp.Seconds || p[i].Timestamp.Nanos > p[j].Timestamp.Nanos {
		return true
	}
	return false
}

// Define swap over an array
func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

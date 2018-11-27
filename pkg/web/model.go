package web

import (
	"sync"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

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

var MaxFeedsNum int = 20

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

// Use temporarily to generate unique ID
var userIDCounter int

// Render Template Helpers
// func renderTemplate(w http.ResponseWriter, tmplname string, te)
type UserTmpl struct {
	UserName     string
	NumTweets    int
	NumFollowing int
	NumFollowers int
	TweetList    []*protocol.Tweet
	IsFollowing  bool
	NotSelf      bool
}

type UserListTmpl struct {
	AlreadyFollowed bool
	Following       bool
	UserName        string
	UserList        []string
}

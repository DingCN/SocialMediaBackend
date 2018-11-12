package web

import "time"

// global variables for storing user data goes here, will be replaced with database later
var username_password map[string]string

// Alice: [Bob, Cain]
// Alice is following Bob and Cain
// hope this definition is correct
var followingList map[string][]string

//Alice: [[Post1, timestamp], [Post2, timestamp]]
var Posts map[string][]TimedPost

type TimedPost struct {
	Post      string
	timestamp time.Time
}

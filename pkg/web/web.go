package web

import (
	"context"
	"fmt"
	"net/http"
)

// global variables for storing user data goes here, will be replaced with database later
//var username_password map[string]string

// Alice: [Bob, Cain]
// Alice is following Bob and Cain
// hope this definition is correct
//var followingList map[string][]string

type Web struct {
	srv *http.Server
}

func New(cfg *Config) (*Web, error) {
	return &Web{
		srv: &http.Server{
			Addr: cfg.Addr,
		},
	}, nil
}

func (w *Web) Start() error {

	http.HandleFunc("/", index)              // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	return err
}

func (w *Web) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

//Business Logic
func login(w http.ResponseWriter, r *http.Request) {
	//use POST method in front end

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	//Test------
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	//--------

	//value, _ := username_password[username]
	currentUser := UserList[username]
	if password == currentUser.Password {
		//TODO return login success

	} else {
		//TODO return login failed
	}
}

// CreateAccount :
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	//use POST method in front end
	r.ParseForm()
	// logic part of log in
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	// Test -------------
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	//-------------------
	// Check if user exist
	user, ok := UserList[username]

	if ok == true {
		// user already exist
		// ask user to login instead
		http.Redirect(w, r, "/", http.StatusFound)
		//TODO if password match, redirect to login success page
	} else {
		// create new user
		newUser := User{
			UserID:   userIDCounter,
			UserName: username,
			Password: password,
			Auth:     "true",
		}

		// TODO: make thread-safe later
		userIDCounter++
		UserList[username] = newUser

		//TODO redirect to login home page
	}
	if len(password) < 6 {
		//TODO return login failed, password length less than 6
	}
}

// FollowOrUnfollow :
func FollowOrUnfollow(w http.ResponseWriter, r *http.Request) {
	targetUserName := r.FormValue("username")
	currUserName := r.FormValue("currentuser")

	// TODO check if logged in

	//
	if r.FormValue("follow") == "1" {
		UserList[currUserName].FollowingList[targetUserName] = true
	} else {
		// follow == "0" unfollow
		UserList[currUserName].FollowingList[targetUserName] = false
	}
	// TODO redirect to homepage
}

// TODO create post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	currUserName := r.URL.Path[len("/home/"):]
	tweetBody := r.FormValue("body")
	newTweet := Tweet{userName: currUserName,
		timestamp: getTimeStamp(),
		body:      tweetBody,
	}
	CentralTweetList = append(CentralTweetList, newTweet)
	// TODO: consider if every User should have a copy of their own tweets

	// TODO: redirect to /home
}

// TODO view feeds
func ViewFeeds(w http.ResponseWriter, r *http.Request) {

	currUserName := r.URL.Path[len("/home/"):]
	// iterate through user's following list, then iterate throuth their posts list, O(n*m)

	var feeds []Tweet
	// One way to do this is to filter from CentralTweetList

	for _, t := range CentralTweetList {
		if UserList[currUserName].IsFollowing(t.userName) {
			feeds = append(feeds, t)
		}
	}

	// TODO: render timeline

}

// TweetGround: display feeds posted by all users
func TweetGround(w http.ResponseWriter, r *http.Request) {

}

// func login(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method) //get request method
// 	if r.Method == "GET" {
// 		t, _ := template.ParseFiles("login.gtpl")
// 		t.Execute(w, nil)
// 	} else {
// 		r.ParseForm()
// 		// logic part of log in
// 		username := r.Form["username"][0]
// 		password := r.Form["password"][0]
// 		fmt.Println("username:", username)
// 		fmt.Println("password:", password)
// 		value, ok := username_password[username]
// 		if password == value {
// 			//TODO return login success
// 		} else {
// 			//TODO return login failed
// 		}
// 	}
// }

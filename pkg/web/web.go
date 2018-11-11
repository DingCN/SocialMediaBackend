package web

import (
	"context"
	"fmt"
	"net/http"
)

// global variables for storing user data goes here, will be replaced with database later
var username_password map[string]string

// Alice: [Bob, Cain]
// Alice is following Bob and Cain
// hope this definition is correct
var followingList map[string][]string

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
	r.ParseForm()
	// logic part of log in
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	value, _ := username_password[username]
	if password == value {
		//TODO return login success
	} else {
		//TODO return login failed
	}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	//use POST method in front end
	r.ParseForm()
	// logic part of log in
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	value, _ := username_password[username]
	if len(password) < 6 {
		//TODO return login failed, password length less than 6
	}
	if value == "" {
		username_password[username] = password
		//TODO return login success
	} else {
		//TODO return login failed, user already exists
	}
}

// TODO Post request to follow or unfollow a target user
func FollowOrUnfollow(w http.ResponseWriter, r *http.Request) {

}

// TODO Get request to check if user is following a target user
func GetFollowingStatus(w http.ResponseWriter, r *http.Request) {

}

// TODO create post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// store in a dict, user:[post1, post2, post3]
	// remember to append timestamp
}

// TODO view feeds
func ViewFeeds(w http.ResponseWriter, r *http.Request) {
	// iterate through user's following list, then iterate throuth their posts list, O(n*m)

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

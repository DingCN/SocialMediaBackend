package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// global variables for storing user data goes here, will be replaced with database later
var username_password = map[string]string{}

// Alice: [Bob, Cain]
// Alice is following Bob and Cain
// hope this definition is correct
var followingList = map[string][]string{}

//Alice: [[Post1, timestamp], [Post2, timestamp]]
var Posts = map[string][]TimedPost{}

type TimedPost struct {
	Post      string
	timestamp time.Time
}

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

	http.HandleFunc("/", index) // set router
	http.HandleFunc("/login", Login)
	http.HandleFunc("/createAccount", CreateAccount)
	http.HandleFunc("/getAllFollowing", GetAllFollowing)
	http.HandleFunc("/createPost", CreatePost)

	err := http.ListenAndServe(":9090", nil) // set listen port
	return err
}

func (w *Web) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

//Business Logic
func Login(w http.ResponseWriter, r *http.Request) {
	//tested
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		value, _ := username_password[username]
		if password == value {
			expiration := time.Now().Add(30 * time.Second)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			_, err := w.Write([]byte("login success"))
			if err != nil {
				panic(err)
			}
		} else {

			_, err := w.Write([]byte("username or passwd incorrect"))
			if err != nil {
				panic(err)
			}

		}
	}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	//tested
	if r.Method == "GET" {
		t, _ := template.ParseFiles("createAccount.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		value, _ := username_password[username]
		if len(password) < 1 {
			_, err := w.Write([]byte("password length less than 1"))
			if err != nil {
				panic(err)
			}
		}
		if value == "" {
			username_password[username] = password
			_, err := w.Write([]byte("create account success"))
			if err != nil {
				panic(err)
			}
		} else {
			_, err := w.Write([]byte("user already exists"))
			if err != nil {
				panic(err)
			}
		}
	}
}

// TODO Post request to follow or unfollow a target user
func FollowOrUnfollow(w http.ResponseWriter, r *http.Request) {
	//input username to follow

}

// TODO Get request to check the following list of a user
func GetAllFollowing(w http.ResponseWriter, r *http.Request) {
	//TODO not finished yet
	r.ParseForm()
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		_, err := w.Write([]byte("login first to get following list"))
		if err != nil {
			panic(err)
		}

	}
	username := cookie.Value
	followings := followingList[username]
	returnList := ""
	for _, user := range followings {
		returnList = returnList + user
		returnList = returnList + " "

	}
	_, err := w.Write([]byte("Users you are following: " + returnList))
	if err != nil {
		panic(err)
	}
}

// TODO Get request to check if user is following a target user
func IfFollowing(w http.ResponseWriter, r *http.Request) {

}

// TODO Post request create post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	//tested
	if r.Method == "GET" {
		t, _ := template.ParseFiles("createPost.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		cookie, _ := r.Cookie("username")
		if cookie == nil {
			_, err := w.Write([]byte("login first to create post"))
			if err != nil {
				panic(err)
			}

		}
		username := cookie.Value
		post := r.Form["Post"][0]
		if post == "" {
			_, err := w.Write([]byte("input is empty"))
			if err != nil {
				panic(err)
			}
		}
		timedPost := TimedPost{post, time.Now()}

		Posts[username] = append(Posts[username], timedPost)
		_, err := w.Write([]byte("create post success"))
		if err != nil {
			panic(err)
		}
	}
}

// TODO view feeds
func ViewFeeds(w http.ResponseWriter, r *http.Request) {
	// iterate through user's following list, then iterate throuth their posts list, O(n*m)
	// sort them by timestamp

}

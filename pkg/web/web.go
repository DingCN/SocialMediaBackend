package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Web server
type Web struct {
	srv *http.Server
}

// New config
func New(cfg *Config) (*Web, error) {
	return &Web{
		srv: &http.Server{
			Addr: cfg.Addr,
		},
	}, nil
}

// Start server
func (w *Web) Start() error {

	http.HandleFunc("/", index) // set router
	http.HandleFunc("/login", Login)
	http.HandleFunc("/createAccount", CreateAccount)
	http.HandleFunc("/getAllFollowing", GetAllFollowing)
	http.HandleFunc("/getAllFollower", GetAllFollower)
	http.HandleFunc("/createPost", CreatePost)
	http.HandleFunc("/ViewFeeds", ViewFeeds)
	http.HandleFunc("/i/moments", MomentRandomFeeds)
	//http.HandleFunc("/ListUser", ListUser)

	err := http.ListenAndServe(":8080", nil) // set listen port
	return err
}

func (w *Web) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

// Business Logic
// func ListUser(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		userList := OPGetAllUsers()
// 		json.NewEncoder(w).Encode(userList)
// 		return
// 	}
// 	//POST request, follow or unfollow
// 	cookie, _ := r.Cookie("username")
// 	if cookie == nil {
// 		json.NewEncoder(w).Encode("login first to follow")
// 		return
// 	}
// 	r.ParseForm()
// 	target, ok := r.URL.Query()["username"]

// 	if !ok || len(target[0]) < 1 {
// 		log.Println("Url Param 'username' is missing")
// 		return
// 	}

// 	OPFollowUnFollow(username, target[0])

// }

// func UserHome(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 	} else {
// 	}
// }

// Login .
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		pUser, ok := UserList.Users[username]
		if ok == false { // not found
			json.NewEncoder(w).Encode("username or passwd incorrect")
			return
		}
		if password == pUser.Password {
			expiration := time.Now().Add(30 * time.Minute)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			json.NewEncoder(w).Encode("login success")
			return
		}
	}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("createAccount.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		if len(password) < 1 {
			json.NewEncoder(w).Encode("password length less than 1")
			return
		}
		_, ok := UserList.Users[username]
		if ok == false { // record not found, creating account...
			OPAddUser(username, password)
			json.NewEncoder(w).Encode("create account success")
			return
		} else {
			json.NewEncoder(w).Encode("user already exists")
			return
		}
	}
}

// TODO Post request to follow or unfollow a target user
func FollowOrUnfollow(w http.ResponseWriter, r *http.Request) {
	//input username to follow
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		json.NewEncoder(w).Encode("login first to follow")
		return
	}
	username := cookie.Value
	target, ok := r.URL.Query()["username"]
	if !ok || len(target[0]) < 1 {
		json.NewEncoder(w).Encode("url parameter incorrect")
		return
	}

	OPFollowUnFollow(username, target[0])

}

func GetAllFollowing(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		json.NewEncoder(w).Encode("login first to get following list")
		return
	}
	username := cookie.Value

	followings := UserList.Users[username].FollowingList
	//TODO render
	json.NewEncoder(w).Encode(followings)
	return

	// returnList := ""
	// for user, isFollowing := range followings {
	// 	if isFollowing == true {
	// 		returnList = returnList + user
	// 		returnList = returnList + " "
	// 	}

	// }
	// _, err := w.Write([]byte("Users you are following: " + returnList))
	// if err != nil {
	// 	panic(err)
	// }
}

func GetAllFollower(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		json.NewEncoder(w).Encode("login first to get follower list")
		return
	}
	username := cookie.Value

	followers := UserList.Users[username].FollowerList
	//TODO render
	json.NewEncoder(w).Encode(followers)
	return

}

//  Get request to check if user is following a target user
func IfFollowing(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		json.NewEncoder(w).Encode(false)
		return
	}

	username := cookie.Value
	target, ok := r.URL.Query()["username"]
	if !ok || len(target[0]) < 1 {
		json.NewEncoder(w).Encode(false)
		return
	}
	res := OPCheckIfFollowing(username, target[0])
	json.NewEncoder(w).Encode(res)
	return
}

// Post request create post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	//tested
	if r.Method == "GET" {
		t, _ := template.ParseFiles("createPost.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		cookie, _ := r.Cookie("username")
		if cookie == nil {
			json.NewEncoder(w).Encode("login first to create post")
			return

		}
		username := cookie.Value
		post := r.PostFormValue("Post")
		if post == "" {
			json.NewEncoder(w).Encode("input is empty")
			return
		}
		OPAddTweet(username, post)
		json.NewEncoder(w).Encode("create post success")
		return
	}
}

// ViewFeeds ..
func ViewFeeds(w http.ResponseWriter, r *http.Request) {
	// iterate through user's following list, then iterate throuth their posts list, O(n*m)
	// sort them by timestamp

	//https://golangcode.com/get-a-url-parameter-from-a-request/
	usernames, ok := r.URL.Query()["username"]

	if !ok || len(usernames[0]) < 1 {
		log.Println("Url Param 'username' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	username := usernames[0]
	unsortedTweets := OPGetFollowingTweets(username)
	sortedTweets := OPSortTweets(unsortedTweets)
	json.NewEncoder(w).Encode(sortedTweets)
	return
	//TODO render
	// var res string
	// for _, tweet := range sortedTweets {
	// 	res += tweet.Body

	// }

	// _, err := w.Write([]byte(res))
	// if err != nil {
	// 	panic(err)
	// }
}

func MomentRandomFeeds(w http.ResponseWriter, r *http.Request) {
	tweets := OPGetRandomTweet()
	json.NewEncoder(w).Encode(tweets)
	return
}

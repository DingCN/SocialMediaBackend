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

	http.HandleFunc("/", Index) // set router
	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/createAccount", CreateAccount)
	http.HandleFunc("/getAllFollowing", GetAllFollowing)
	http.HandleFunc("/getAllFollower", GetAllFollower)
	http.HandleFunc("/createPost", CreatePost)
	http.HandleFunc("/userProfile", UserProfile)
	http.HandleFunc("/i/moments", MomentRandomFeeds)
	//http.HandleFunc("/ListUser", ListUser)

	err := http.ListenAndServe(":8080", nil) // set listen port
	return err
}

func (w *Web) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

// Index .
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

// Login .
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, _ := template.ParseFiles("index.html")

	loginResult := ""
	// logic part of log in
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	//Check is login or sign up

	if r.PostFormValue("login") != "" {
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		pUser, ok := UserList.Users[username]
		if ok == false { // not found
			json.NewEncoder(w).Encode("username or passwd incorrect")
			loginResult = "User not found. Please try again."
			t.Execute(w, loginResult)
		} else if password != pUser.Password {
			json.NewEncoder(w).Encode("username or passwd incorrect")
			loginResult := "Incorrect password. Please try again."
			t.Execute(w, loginResult)
		} else { // login success, redirect to home
			newURL := fmt.Sprintf("/home?username=%s", username)
			http.Redirect(w, r, newURL, 302)
		}
	} else if r.PostFormValue("signup") != "" {
		CreateAccount(w, r)
	}
}

// Home .
func Home(w http.ResponseWriter, r *http.Request) {
	usernames, ok := r.URL.Query()["username"]
	if !ok {
		return
	}

	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "username", Value: usernames[0], Expires: expiration}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode("login success")

	// Render
	pUser, ok := UserList.Users[usernames[0]]
	h, _ := template.ParseFiles("home.html")
	unsortedTweets := OPGetFollowingTweets(pUser.UserName)
	sortedTweets := OPSortTweets(unsortedTweets)
	userHome := UserTmpl{
		username:     usernames[0],
		numTweets:    len(pUser.TweetList),
		numFollowing: len(pUser.FollowingList),
		numFollowers: len(pUser.FollowerList),
		tweetList:    sortedTweets,
	}
	h.Execute(w, userHome)
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
			t, _ := template.ParseFiles("index.html")
			signUpResult := "password length less than 1"
			t, err := template.ParseFiles("index.html")
			if err != nil {
				panic(err)
			}
			t.Execute(w, signUpResult)
		}
		_, ok := UserList.Users[username]
		if ok == false { // record not found, creating account...
			OPAddUser(username, password)
			json.NewEncoder(w).Encode("create account success")
			newURL := fmt.Sprintf("/home?username=%s", username)
			http.Redirect(w, r, newURL, 302)

		} else {
			t, _ := template.ParseFiles("index.html")
			signUpResult := "user already exists"
			json.NewEncoder(w).Encode("user already exists")
			t.Execute(w, signUpResult)
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
	usernames, ok := r.URL.Query()["username"]
	if !ok || len(usernames[0]) < 1 {
		json.NewEncoder(w).Encode("url parameter incorrect")
		return
	}
	username := usernames[0]

	followings := UserList.Users[username].FollowingList
	//TODO render
	json.NewEncoder(w).Encode(followings)
	t, err := template.ParseFiles("userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowingTmpl := UserListTmpl{
		listType: "following",
		username: username,
		userlist: followings,
	}
	json.NewEncoder(w).Encode(followings)
	t.Execute(w, newFollowingTmpl)

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
	t, err := template.ParseFiles("userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowerTmpl := UserListTmpl{
		listType: "follower",
		username: username,
		userlist: followers,
	}
	json.NewEncoder(w).Encode(followers)
	t.Execute(w, newFollowerTmpl)
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

		newURL := fmt.Sprintf("/home?username=%s", username)
		http.Redirect(w, r, newURL, 302)
		return
	}
}

// ViewFeeds ..
// Usage modification: change to view another user's profile page
func UserProfile(w http.ResponseWriter, r *http.Request) {
	//https://golangcode.com/get-a-url-parameter-from-a-request/
	usernames, ok := r.URL.Query()["username"]

	if !ok || len(usernames[0]) < 1 {
		log.Println("Url Param 'username' is missing")
		return
	}
	// Query()["key"] will return an array of items,
	// we only want the single item.
	pUser, ok := UserList.Users[usernames[0]]
	h, _ := template.ParseFiles("userprofile.html")
	userProfile := UserTmpl{
		username:     usernames[0],
		numTweets:    len(pUser.TweetList),
		numFollowing: len(pUser.FollowingList),
		numFollowers: len(pUser.FollowerList),
		tweetList:    pUser.TweetList,
	}
	h.Execute(w, userProfile)
}

func MomentRandomFeeds(w http.ResponseWriter, r *http.Request) {
	tweets := OPGetRandomTweet()
	json.NewEncoder(w).Encode(tweets)
	t, _ := template.ParseFiles("moments.html")
	t.Execute(w, tweets)
}

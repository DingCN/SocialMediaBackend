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
	http.HandleFunc("/home", Home) // view feeds
	http.HandleFunc("/createAccount", CreateAccount)
	http.HandleFunc("/getAllFollowing", GetAllFollowing)
	http.HandleFunc("/getAllFollower", GetAllFollower)
	http.HandleFunc("/createPost", CreatePost)
	http.HandleFunc("/userProfile", UserProfile) //tweet for a single user
	http.HandleFunc("/i/moments", MomentRandomFeeds)
	http.HandleFunc("/FollowOrUnfollow", FollowOrUnfollow)
	//http.HandleFunc("/ListUser", ListUser)

	err := http.ListenAndServe(":8080", nil) // set listen port
	return err
}

func (w *Web) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

// Index .
func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

// Login .
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, _ := template.ParseFiles("frontend/index.html")

	//loginResult := ""
	// logic part of log in
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	//Check is login or sign up

	if r.PostFormValue("login") != "" {
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		// pUser, ok := UserList.Users[username]
		_, ok := UserList.Users[username]
		if ok == false { // not found
			// 	loginResult = "User not found. Please try again."
			// 	t.Execute(w, loginResult)
			// } else if password != pUser.Password {
			loginResult := "Incorrect password. Please try again."
			t.Execute(w, loginResult)
			json.NewEncoder(w).Encode("username or passwd incorrect")
		} else { // login success, redirect to home

			expiration := time.Now().Add(30 * time.Minute)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", 302)
			//Test
			json.NewEncoder(w).Encode("login success")
		}
	} else if r.PostFormValue("signup") != "" {
		CreateAccount(w, r)
	}
}

// Home .
func Home(w http.ResponseWriter, r *http.Request) {
	// usernames, ok := r.URL.Query()["username"]
	cookie, _ := r.Cookie("username")
	username := cookie.Value

	// Render
	pUser, ok := UserList.Users[username]
	if !ok {
		log.Println("login first to get follower list")
	}
	h, err := template.ParseFiles("frontend/home.html")
	if err != nil {
		panic(err)
	}
	unsortedTweets := OPGetFollowingTweets(pUser.UserName)
	fmt.Printf("Following post for user: %s found: ", username)
	for _, tweet := range unsortedTweets {
		fmt.Printf("%s; ", tweet.Body)
	}
	fmt.Printf("\n")
	sortedTweets := OPSortTweets(unsortedTweets)
	userHome := UserTmpl{
		UserName:     username,
		NumTweets:    len(pUser.TweetList),
		NumFollowing: len(pUser.FollowingList),
		NumFollowers: len(pUser.FollowerList),
		TweetList:    sortedTweets,
	}
	err = h.Execute(w, userHome)
	if err != nil {
		panic(err)
	}
	// log
	fmt.Printf("Following post for user: %s found: ", username)
	for _, tweet := range sortedTweets {
		fmt.Printf("%s; ", tweet.Body)
	}
	fmt.Printf("\n")
	json.NewEncoder(w).Encode(sortedTweets)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	//var signUpResult string

	if len(password) < 1 {

		t, _ := template.ParseFiles("frontend/index.html")
		t, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, "password length less than 1")
		json.NewEncoder(w).Encode("password length less than 1")
	}
	_, ok := UserList.Users[username]
	if ok == false { // record not found, creating account...
		OPAddUser(username, password)

		// redirect to login
		// TODO: redirect to Login to save repeating code
		expiration := time.Now().Add(30 * time.Minute)
		cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/home", 302)

		json.NewEncoder(w).Encode("create account success")

	} else {
		t, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, "user already exists")
		// json.NewEncoder(w).Encode("user already exists")

	}
}

// TODO Post request to follow or unfollow a target user
func FollowOrUnfollow(w http.ResponseWriter, r *http.Request) {
	//input username to follow
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		//json.NewEncoder(w).Encode("login first to follow")
		log.Println("login first to get follower list")
		return
	}
	username := cookie.Value
	target, ok := r.URL.Query()["username"]
	if !ok || len(target[0]) < 1 {
		//json.NewEncoder(w).Encode("url parameter incorrect")
		return
	}

	OPFollowUnFollow(username, target[0])

	newURL := fmt.Sprintf("/userprofile?username=%s", target)
	http.Redirect(w, r, newURL, 302)
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
	//json.NewEncoder(w).Encode(followings)
	t, err := template.ParseFiles("frontend/userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowingTmpl := UserListTmpl{
		AlreadyFollowed: false,
		Following:       true,
		UserName:        username,
		UserList:        followings,
	}
	//json.NewEncoder(w).Encode(followings)
	t.Execute(w, newFollowingTmpl)
	json.NewEncoder(w).Encode(followings)

}

func GetAllFollower(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	usernames, ok := r.URL.Query()["username"]
	if !ok || len(usernames[0]) < 1 {
		json.NewEncoder(w).Encode("url parameter incorrect")
		return
	}
	username := usernames[0]

	followers := UserList.Users[username].FollowerList
	//TODO render
	t, err := template.ParseFiles("frontend/userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowerTmpl := UserListTmpl{
		AlreadyFollowed: false,
		Following:       false,
		UserName:        username,
		UserList:        followers,
	}

	t.Execute(w, newFollowerTmpl)
	//Test
	json.NewEncoder(w).Encode(followers)
	return

}

//  Get request to check if user is following a target user
func IfFollowing(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, _ := r.Cookie("username")
	if cookie == nil {
		//json.NewEncoder(w).Encode(false)
		return
	}

	//username := cookie.Value
	target, ok := r.URL.Query()["username"]
	if !ok || len(target[0]) < 1 {
		//json.NewEncoder(w).Encode(false)
		return
	}
	//res := OPCheckIfFollowing(username, target[0])
	//json.NewEncoder(w).Encode(res)
	return
}

// Post request create post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	//tested
	if r.Method == "GET" {
		t, _ := template.ParseFiles("frontend/createPost.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		cookie, _ := r.Cookie("username")
		if cookie == nil {
			//json.NewEncoder(w).Encode("login first to create post")
			return

		}
		username := cookie.Value
		post := r.PostFormValue("body")
		if post == "" {
			//json.NewEncoder(w).Encode("input is empty")
			//
			return
		}
		OPAddTweet(username, post)

		newURL := fmt.Sprintf("/home?username=%s", username)
		http.Redirect(w, r, newURL, 302)
		json.NewEncoder(w).Encode("create post success")
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
	username := usernames[0]
	// Query()["key"] will return an array of items,
	// we only want the single item.
	pUser, ok := UserList.Users[usernames[0]]
	h, err := template.ParseFiles("frontend/userprofile.html")
	if err != nil {
		panic(err)
	}
	userProfile := UserTmpl{
		UserName:     usernames[0],
		NumTweets:    len(pUser.TweetList),
		NumFollowing: len(pUser.FollowingList),
		NumFollowers: len(pUser.FollowerList),
		TweetList:    pUser.TweetList,
	}
	h.Execute(w, userProfile)
	//Test
	tweets := OPGetTweetByUsername(username)
	json.NewEncoder(w).Encode(tweets)
}

func MomentRandomFeeds(w http.ResponseWriter, r *http.Request) {
	tweets := OPGetRandomTweet()
	t, err := template.ParseFiles("frontend/moments.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, tweets)
	json.NewEncoder(w).Encode(tweets)
}

package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
	"google.golang.org/grpc"
)

// Web server
type Web struct {
	srv *http.Server
	//client handle when comm with backend
	c protocol.TwitterRPCClient
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
func (web *Web) Start() error {
	backendAddr := "localhost:50051"
	conn, err := grpc.Dial(backendAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("web adn backend did not connect: %v", err)
		return err

	}
	defer conn.Close()
	web.c = protocol.NewTwitterRPCClient(conn)

	// Contact the backend and print out its response.

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// r, err := c.SignupRPC(ctx, &pb.HelloRequest{Name: name})
	// if err != nil {
	// 	log.Fatalf("could not signup: %v", err)
	// }

	http.HandleFunc("/", web.Index) // set router
	http.HandleFunc("/login", web.Login)
	http.HandleFunc("/home", web.Home) // view feeds
	http.HandleFunc("/createAccount", web.CreateAccount)
	http.HandleFunc("/getAllFollowing", web.GetAllFollowing)
	http.HandleFunc("/getAllFollower", web.GetAllFollower)
	http.HandleFunc("/createPost", web.CreatePost)
	http.HandleFunc("/userProfile", web.UserProfile) //tweet for a single user
	http.HandleFunc("/i/moments", web.MomentRandomFeeds)
	http.HandleFunc("/FollowOrUnfollow", web.FollowOrUnfollow)
	//http.HandleFunc("/ListUser", ListUser)

	err = http.ListenAndServe(":8080", nil) // set listen port
	return err
}

func (web *Web) Shutdown(ctx context.Context) error {
	return web.srv.Shutdown(ctx)
}

// Index .
func (web *Web) Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

// Login .
func (web *Web) Login(w http.ResponseWriter, r *http.Request) {
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

		loginReply, err := web.LoginRPCSend(username, password)
		if err == nil && loginReply.Success == true {
			// success
			expiration := time.Now().Add(30 * time.Minute)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", 302)
			//Test
			json.NewEncoder(w).Encode("login success")
		} else if err.Error() == errorcode.ErrUserNotExist.Error() {
			loginResult := "User not found. Please try again."
			t.Execute(w, loginResult)
		} else if err.Error() == errorcode.ErrIncorrectPassword.Error() {
			loginResult := "Incorrect password. Please try again."
			t.Execute(w, loginResult)
		}
	} else if r.PostFormValue("signup") != "" {
		web.CreateAccount(w, r)
	}
}

// Home .
func (web *Web) Home(w http.ResponseWriter, r *http.Request) {
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
	sortedTweets := OPGetFollowingTweets(pUser.UserName)
	fmt.Printf("Following post for user: %s found: ", username)
	for _, tweet := range sortedTweets {
		fmt.Printf("%s; ", tweet.Body)
	}
	fmt.Printf("\n")
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
}

func (web *Web) CreateAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	//var signUpResult string
	_, err := web.SignupRPCSend(username, password)

	if err == nil { // success

		// redirect to login
		// TODO: redirect to Login to save repeating code
		expiration := time.Now().Add(30 * time.Minute)
		cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/home", 302)

		json.NewEncoder(w).Encode("create account success")
	} else if err.Error() == errorcode.ErrInvalidUsername.Error() {
		t, _ := template.ParseFiles("frontend/index.html")
		t, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, "Cannot signup, username invalid")
	} else if err.Error() == errorcode.ErrInvalidPassword.Error() {
		t, _ := template.ParseFiles("frontend/index.html")
		t, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, "Cannot signup, password length not enough")
	} else if err.Error() == errorcode.ErrUsernameTaken.Error() {
		t, _ := template.ParseFiles("frontend/index.html")
		t, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, "Cannot signup, username already taken")
	}
}

// TODO Post request to follow or unfollow a target user
func (web *Web) FollowOrUnfollow(w http.ResponseWriter, r *http.Request) {
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
	_, err := web.FollowUnFollowRPCSend(username, target[0])
	if err == nil { // success
		newURL := fmt.Sprintf("/userProfile?username=%s", target[0])
		http.Redirect(w, r, newURL, 302)
	} else {
		json.NewEncoder(w).Encode(err.Error())
	}
}

func (web *Web) GetAllFollowing(w http.ResponseWriter, r *http.Request) {
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
	t.Execute(w, newFollowingTmpl)
	//json.NewEncoder(w).Encode(followings)

}

func (web *Web) GetAllFollower(w http.ResponseWriter, r *http.Request) {
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
	//json.NewEncoder(w).Encode(followers)
	return

}

//  Get request to check if user is following a target user
func (web *Web) IfFollowing(w http.ResponseWriter, r *http.Request) {
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
func (web *Web) CreatePost(w http.ResponseWriter, r *http.Request) {
	//tested
	if r.Method == "GET" {
		t, _ := template.ParseFiles("frontend/createPost.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		cookie, _ := r.Cookie("username")
		if cookie == nil {
			json.NewEncoder(w).Encode("login first to create post")
			return

		}
		username := cookie.Value
		post := r.PostFormValue("body")
		if post == "" {
			json.NewEncoder(w).Encode("input is empty")
			//
			return
		}

		_, err := web.AddTweetRPCSend(username, post)

		if err == nil { // success

			newURL := fmt.Sprintf("/home?username=%s", username)
			http.Redirect(w, r, newURL, 302)
			json.NewEncoder(w).Encode("create post success")
		} else {
			t, _ := template.ParseFiles("frontend/createPost.html")
			t.Execute(w, nil)
		}
		return
	}
}

// ViewFeeds ..
// Usage modification: change to view another user's profile page
func (web *Web) UserProfile(w http.ResponseWriter, r *http.Request) {
	//https://golangcode.com/get-a-url-parameter-from-a-request/
	usernames, ok := r.URL.Query()["username"]

	if !ok || len(usernames[0]) < 1 {
		log.Println("Url Param 'username' is missing")
		return
	}
	username := usernames[0]
	// Query()["key"] will return an array of items,
	// we only want the single item.
	pUser, ok := UserList.Users[username]
	h, err := template.ParseFiles("frontend/userprofile.html")
	if err != nil {
		panic(err)
	}
	userProfile := UserTmpl{
		UserName:     username,
		NumTweets:    len(pUser.TweetList),
		NumFollowing: len(pUser.FollowingList),
		NumFollowers: len(pUser.FollowerList),
		TweetList:    pUser.TweetList,
	}
	h.Execute(w, userProfile)

}

func (web *Web) MomentRandomFeeds(w http.ResponseWriter, r *http.Request) {
	tweets := OPGetRandomTweet()
	t, err := template.ParseFiles("frontend/moments.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, tweets)
}

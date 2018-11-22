package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/web"
)

var addr = "//127.0.0.1:8080"

///////////////////////////////////////////////////////////////////////
//////////////// module tests starts///////////////////////////////////
///////////////////////////////////////////////////////////////////////

// When a user registers, he isn't following any other users.
// We provide a moment page so that it can get the newest posts even he is not following their owner
func Test_Moments(t *testing.T) {
	// env set
	ForTestCreateAccount(t, "Test_MomentsAlice", "Test_MomentsAlice")
	ForTestCreateAccount(t, "Test_MomentsBob", "Test_MomentsBob")
	ForTestCreateAccount(t, "Test_MomentsCain", "Test_MomentsCain")
	ForTestCreatePost(t, "Test_MomentsBob", "Test_MomentsBob's post")
	ForTestCreatePost(t, "Test_MomentsCain", "Test_MomentsCain's post")

	tweets := web.OPGetRandomTweet()
	if len(tweets) != 2 || tweets[1].UserName != "Test_MomentsBob" || tweets[1].Body != "Test_MomentsBob's post" || tweets[0].UserName != "Test_MomentsCain" || tweets[0].Body != "Test_MomentsCain's post" {
		t.Fatalf("Moments incorrect")
		fmt.Printf("%+v\n", tweets)
	}
}

//Test for view feeds
// user only get feeds for those he/she follows
// ordered by timestamp, new to old
// Alice follows Bob and Cain, and she only gets feed by these two
func Test_Home(t *testing.T) {
	// env set
	ForTestCreateAccount(t, "Test_HomeAlice", "Test_HomeAlice")
	ForTestCreateAccount(t, "Test_HomeBob", "Test_HomeBob")
	ForTestCreateAccount(t, "Test_HomeCain", "Test_HomeCain")
	ForTestCreateAccount(t, "Test_HomeDoge", "Test_HomeDoge")
	ForTestFollowUnFollow(t, "Test_HomeAlice", "Test_HomeBob")
	ForTestFollowUnFollow(t, "Test_HomeAlice", "Test_HomeCain")
	// Alice is following Bob
	ForTestCreatePost(t, "Test_HomeBob", "Test_HomeBob's post")
	ForTestCreatePost(t, "Test_HomeCain", "Test_HomeCain's post")
	ForTestCreatePost(t, "Test_HomeBob", "Test_HomeBob's post2")
	ForTestCreatePost(t, "Test_HomeDoge", "Test_HomeDoge's post")

	//module test
	username := "Test_HomeAlice"
	pUser, ok := web.UserList.Users[username]
	if !ok {
		log.Println("Test home fails")
	}
	unsortedTweets := web.OPGetFollowingTweets(pUser.UserName)
	fmt.Printf("Following post for user: %s found: ", username)
	for _, tweet := range unsortedTweets {
		fmt.Printf("%s; ", tweet.Body)
	}
	fmt.Printf("\n")
	sortedTweets := web.OPSortTweets(unsortedTweets)
	if len(sortedTweets) != 3 || sortedTweets[2].UserName != "Test_HomeBob" || sortedTweets[2].Body != "Test_HomeBob's post" || sortedTweets[1].UserName != "Test_HomeCain" || sortedTweets[1].Body != "Test_HomeCain's post" || sortedTweets[0].UserName != "Test_HomeBob" || sortedTweets[0].Body != "Test_HomeBob's post2" {
		t.Fatalf("Home(ViewFeeds) incorrect")
	}
}
func Test_UserProfile(t *testing.T) {
	// env set
	ForTestCreateAccount(t, "Test_UserProfileAlice", "Test_UserProfileAlice")
	ForTestCreateAccount(t, "Test_UserProfileBob", "Test_UserProfileBob")
	ForTestCreateAccount(t, "Test_UserProfileCain", "Test_UserProfileCain")
	// Alice is following Bob
	ForTestCreatePost(t, "Test_UserProfileBob", "Test_UserProfileBob's post")
	ForTestCreatePost(t, "Test_UserProfileCain", "Test_UserProfileCain's post")

	username := "Test_UserProfileBob"
	// Query()["key"] will return an array of items,
	// we only want the single item.
	pUser, _ := web.UserList.Users[username]
	userProfile := web.UserTmpl{
		UserName:     username,
		NumTweets:    len(pUser.TweetList),
		NumFollowing: len(pUser.FollowingList),
		NumFollowers: len(pUser.FollowerList),
		TweetList:    pUser.TweetList,
	}
	if userProfile.UserName != "Test_UserProfileBob" || userProfile.NumTweets != 1 || userProfile.TweetList[0].Body != "Test_UserProfileBob's post" {
		t.Fatalf("UserProfile incorrect")
	}
}

///////////////////////////////////////////////////////////////////////
//////////////// module tests ends/////////////////////////////////////
///////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////
//////////////////// End to End tests//////////////////////////////////
///////////////////////////////////////////////////////////////////////
func Test_CreateAccount(t *testing.T) {
	var expected = "\"create account success\"\n"
	actual := ForTestCreateAccount(t, "test1", "test1")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}
func Test_CreateTwoAccount(t *testing.T) {
	var expected = "\"create account success\"\n"
	actual := ForTestCreateAccount(t, "test2-1", "test2-1")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	expected = "\"create account success\"\n"
	actual = ForTestCreateAccount(t, "test2-2", "test2-2")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// First time create should return success, second time should fail since user already exists
func Test_CreateAccountTwice(t *testing.T) {
	var expected = "\"create account success\"\n"
	actual := ForTestCreateAccount(t, "test3", "test3")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

	res := httptest.NewRecorder()
	tmpl, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(res, "user already exists")
	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	expected = string(body)

	actual = ForTestCreateAccount(t, "test3", "test3")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func Test_Login(t *testing.T) {
	ForTestCreateAccount(t, "test4", "test4")
	var expected = "\"login success\"\n"
	actual := ForTestLogin(t, "test4", "test4")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

	res := httptest.NewRecorder()
	tmpl, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		panic(err)
	}
	loginResult := "Incorrect username or password. Please try again."
	tmpl.Execute(res, loginResult)
	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	expected = string(body)

	actual = ForTestLogin(t, "4tset", "4tset")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func Test_CreatePost(t *testing.T) {
	ForTestCreateAccount(t, "Test_CreatePost", "Test_CreatePost")

	ForTestLogin(t, "Test_CreatePost", "Test_CreatePost")
	/////TODO Fatal bug if forged a non-exist username

	actual := ForTestCreatePost(t, "Test_CreatePost", "Test_CreatePost")
	var expected = "create post success"

	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func Test_FollowUnFollow(t *testing.T) {
	ForTestCreateAccount(t, "Test_FollowAlice", "Test_FollowAlice")
	ForTestCreateAccount(t, "Test_FollowBob", "Test_FollowBob")
	ForTestFollowUnFollow(t, "Test_FollowAlice", "Test_FollowBob")

	// test following list of Alice
	actual := ForTestFollowingList(t, "Test_FollowAlice")
	var followings = map[string]bool{}
	followings["Test_FollowBob"] = true
	res := httptest.NewRecorder()
	tmpl, err := template.ParseFiles("frontend/userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowingTmpl := web.UserListTmpl{
		AlreadyFollowed: false,
		Following:       true,
		UserName:        "Test_FollowAlice",
		UserList:        followings,
	}
	tmpl.Execute(res, newFollowingTmpl)
	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	expected := string(body)
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	// eq := reflect.DeepEqual(actual, expected)
	// if !eq {
	// 	t.Fatalf("FollowingList incorrect")
	// }

	// test follower list of Bob
	actual = ForTestFollowerList(t, "Test_FollowBob")
	var followers = map[string]bool{}
	followers["Test_FollowAlice"] = true
	res = httptest.NewRecorder()
	tmpl, err = template.ParseFiles("frontend/userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowerTmpl := web.UserListTmpl{
		AlreadyFollowed: false,
		Following:       false,
		UserName:        "Test_FollowBob",
		UserList:        followers,
	}
	tmpl.Execute(res, newFollowerTmpl)
	resp = res.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	expected = string(body)
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

	//Testing Unfollow
	ForTestFollowUnFollow(t, "Test_FollowAlice", "Test_FollowBob")

	// test following list of Alice
	actual = ForTestFollowingList(t, "Test_FollowAlice")
	followings = map[string]bool{}

	res = httptest.NewRecorder()
	tmpl, err = template.ParseFiles("frontend/userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowingTmpl = web.UserListTmpl{
		AlreadyFollowed: false,
		Following:       true,
		UserName:        "Test_FollowAlice",
		UserList:        followings,
	}
	tmpl.Execute(res, newFollowingTmpl)
	resp = res.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	expected = string(body)
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

	// test follower list of Bob
	actual = ForTestFollowerList(t, "Test_FollowBob")
	followers = map[string]bool{}
	res = httptest.NewRecorder()
	tmpl, err = template.ParseFiles("frontend/userlist.html")
	if err != nil {
		panic(err)
	}
	newFollowerTmpl = web.UserListTmpl{
		AlreadyFollowed: false,
		Following:       false,
		UserName:        "Test_FollowBob",
		UserList:        followers,
	}
	tmpl.Execute(res, newFollowerTmpl)
	resp = res.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	expected = string(body)
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

}

///////////////////////////////////////////////////////////////////////
//////////////// End to End tests ends/////////////////////////////////
///////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////
//////////////// tests helpers/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

func ForTestCreateAccount(t *testing.T, username string, password string) string {
	var path = "/createAccount.html"
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = form
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.CreateAccount)
	handler.ServeHTTP(res, req)
	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	return string(body)
}

func ForTestLogin(t *testing.T, username string, password string) string {
	var path = "/login.html"
	form := url.Values{}
	form.Add("login", "true")
	form.Add("username", username)
	form.Add("password", password)
	//resp, err = http.PostForm(addr+path, form)
	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = form
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.Login)
	handler.ServeHTTP(res, req)

	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	return string(body)
}

func ForTestCreatePost(t *testing.T, username string, post string) string {
	var path = "/createPost"
	form := url.Values{}
	form.Add("body", post)
	//resp, err = http.PostForm(addr+path, form)
	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = form
	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
	req.AddCookie(&cookie)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.CreatePost)
	handler.ServeHTTP(res, req)
	//CreateAccount(res, req)
	var actual string
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return ""
	}

	return actual
}

func ForTestFollowUnFollow(t *testing.T, username string, targetname string) {
	var path = "/FollowUnfollow"
	var urlparameter = "?username=" + targetname
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = form
	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
	req.AddCookie(&cookie)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.FollowOrUnfollow)
	handler.ServeHTTP(res, req)
	return
}

func ForTestFollowingList(t *testing.T, username string) string {
	var path = "getAllFollowing.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.GetAllFollowing)
	handler.ServeHTTP(res, req)

	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	return string(body)
}

func ForTestFollowerList(t *testing.T, username string) string {
	var path = "getAllFollower.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.GetAllFollower)
	handler.ServeHTTP(res, req)

	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	return string(body)
}

func ForTestHome(t *testing.T, username string) string {

	var path = "home.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
	req.AddCookie(&cookie)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.Home)
	handler.ServeHTTP(res, req)

	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	return string(body)
}

func ForTestUserProfile(t *testing.T, username string) []web.Tweet {
	var path = "userProfile.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.UserProfile)
	handler.ServeHTTP(res, req)

	var actual []web.Tweet
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

func ForTestMoments(t *testing.T) []web.Tweet {
	var path = "i/moments.html"
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(web.MomentRandomFeeds)
	handler.ServeHTTP(res, req)

	var actual []web.Tweet
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		panic(err)
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

var addr = "//127.0.0.1:8080"

// When a user registers, he isn't following any other users.
// We provide a moment page so that it can get the newest posts even he is not following their owner
func Test_Moments(t *testing.T) {
	ForTestCreateAccount(t, "Test_MomentsAlice", "Test_MomentsAlice")
	ForTestCreateAccount(t, "Test_MomentsBob", "Test_MomentsBob")
	ForTestCreateAccount(t, "Test_MomentsCain", "Test_MomentsCain")

	ForTestCreatePost(t, "Test_MomentsBob", "Test_MomentsBob's post")
	ForTestCreatePost(t, "Test_MomentsCain", "Test_MomentsCain's post")
	actual := ForTestMoments(t)
	// unable to test TimeStamp since it's set on server side
	if len(actual) != 2 || actual[1].UserName != "Test_MomentsBob" || actual[1].Body != "Test_MomentsBob's post" || actual[0].UserName != "Test_MomentsCain" || actual[0].Body != "Test_MomentsCain's post" {
		t.Fatalf("Moments incorrect")
		fmt.Printf("%+v\n", actual)
	}
}
func Test_CreateAccount(t *testing.T) {
	var expected = "create account success"
	actual := ForTestCreateAccount(t, "test1", "test1")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}
func Test_CreateTwoAccount(t *testing.T) {
	var expected = "create account success"
	actual := ForTestCreateAccount(t, "test2-1", "test2-1")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	expected = "create account success"
	actual = ForTestCreateAccount(t, "test2-2", "test2-2")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}
func Test_CreateAccountTwice(t *testing.T) {
	var expected = "create account success"
	actual := ForTestCreateAccount(t, "test3", "test3")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	expected = "user already exists"
	actual = ForTestCreateAccount(t, "test3", "test3")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func Test_Login(t *testing.T) {
	ForTestCreateAccount(t, "test4", "test4")
	var expected = "login success"
	actual := ForTestLogin(t, "test4", "test4")
	if actual != expected {
		t.Fatalf("Expected %s got %s", expected, actual)
	}

	expected = "username or passwd incorrect"
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
	var expected = map[string]bool{}
	expected["Test_FollowBob"] = true
	eq := reflect.DeepEqual(actual, expected)
	if !eq {
		t.Fatalf("FollowingList incorrect")
	}

	// test follower list of Bob
	actual = ForTestFollowerList(t, "Test_FollowBob")
	expected = map[string]bool{}
	expected["Test_FollowAlice"] = true
	eq = reflect.DeepEqual(actual, expected)
	if !eq {
		t.Fatalf("FollowerList incorrect")
	}

	//Testing Unfollow
	ForTestFollowUnFollow(t, "Test_FollowAlice", "Test_FollowBob")

	// test following list of Alice
	actual = ForTestFollowingList(t, "Test_FollowAlice")
	expected = map[string]bool{}
	eq = reflect.DeepEqual(actual, expected)
	if !eq {
		t.Fatalf("FollowingList incorrect")
	}

	// test follower list of Bob
	actual = ForTestFollowerList(t, "Test_FollowBob")
	expected = map[string]bool{}
	eq = reflect.DeepEqual(actual, expected)
	if !eq {
		t.Fatalf("FollowerList incorrect")
	}

}

//Test for view feeds
// user only get feeds for those he/she follows
// ordered by timestamp, new to old
// Alice follows Bob and Cain, and she only gets feed by these two
func Test_Home(t *testing.T) {
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

	actual := ForTestHome(t, "Test_HomeAlice")
	// unable to test TimeStamp since it's set on server side
	if len(actual) != 3 || actual[2].UserName != "Test_HomeBob" || actual[2].Body != "Test_HomeBob's post" || actual[1].UserName != "Test_HomeCain" || actual[1].Body != "Test_HomeCain's post" || actual[0].UserName != "Test_HomeBob" || actual[0].Body != "Test_HomeBob's post2" {
		t.Fatalf("Home(ViewFeeds) incorrect")
	}
}
func Test_UserProfile(t *testing.T) {
	ForTestCreateAccount(t, "Test_UserProfileAlice", "Test_UserProfileAlice")
	ForTestCreateAccount(t, "Test_UserProfileBob", "Test_UserProfileBob")
	ForTestCreateAccount(t, "Test_UserProfileCain", "Test_UserProfileCain")
	// Alice is following Bob
	ForTestCreatePost(t, "Test_UserProfileBob", "Test_UserProfileBob's post")
	ForTestCreatePost(t, "Test_UserProfileCain", "Test_UserProfileCain's post")
	actual := ForTestUserProfile(t, "Test_UserProfileBob")
	// unable to test TimeStamp since it's set on server side
	if len(actual) != 1 || actual[0].UserName != "Test_UserProfileBob" || actual[0].Body != "Test_UserProfileBob's post" {
		t.Fatalf("UserProfile incorrect")
	}
}

// func Test_UserProfile(t *testing.T) {
// 	ForTestCreateAccount(t, "Test_UserProfileAlice", "Test_UserProfileAlice")
// 	ForTestCreateAccount(t, "Test_UserProfileBob", "Test_UserProfileBob")
// 	ForTestFollow(t, "Test_UserProfileAlice", "Test_UserProfileBob")
// 	// Alice is following Bob
// 	ForTestCreatePost(t, "Test_UserProfileBob", "Test_UserProfileBob's post")
// 	ForTestUserProfile(t, "Test_UserProfileAlice")
// }

func ForTestCreateAccount(t *testing.T, username string, password string) string {
	var path = "/createAccount.html"
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	//resp, err = http.PostForm(addr+path, form)
	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = form
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateAccount)
	handler.ServeHTTP(res, req)

	var actual string
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return ""
	}
	return actual

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
	handler := http.HandlerFunc(Login)
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

func ForTestCreatePost(t *testing.T, username string, post string) string {
	var path = "/createPost"
	form := url.Values{}
	form.Add("Post", post)
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
	handler := http.HandlerFunc(CreatePost)
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
	handler := http.HandlerFunc(FollowOrUnfollow)
	handler.ServeHTTP(res, req)
	return
}

func ForTestFollowingList(t *testing.T, username string) map[string]bool {
	var path = "getAllFollowing.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllFollowing)
	handler.ServeHTTP(res, req)

	var actual = map[string]bool{}
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

func ForTestFollowerList(t *testing.T, username string) map[string]bool {
	var path = "getAllFollower.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllFollower)
	handler.ServeHTTP(res, req)

	var actual = map[string]bool{}
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

func ForTestHome(t *testing.T, username string) []Tweet {
	var path = "home.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(Home)
	handler.ServeHTTP(res, req)

	var actual []Tweet
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

func ForTestUserProfile(t *testing.T, username string) []Tweet {
	var path = "userProfile.html"
	var urlparameter = "?username=" + username
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(UserProfile)
	handler.ServeHTTP(res, req)

	var actual []Tweet
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

func ForTestMoments(t *testing.T) []Tweet {
	var path = "i/moments.html"
	form := url.Values{}
	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(MomentRandomFeeds)
	handler.ServeHTTP(res, req)

	var actual []Tweet
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		t.Fatalf("HTTP error")
		return nil
	}
	return actual
}

// func ForTestUserProfile(t *testing.T, username string) {
// 	var path = "userProfile.html"
// 	var urlparameter = "?username=" + username
// 	form := url.Values{}
// 	req, err := http.NewRequest("POST", addr+path+urlparameter, strings.NewReader(form.Encode()))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	res := httptest.NewRecorder()
// 	handler := http.HandlerFunc(GetAllFollower)
// 	handler.ServeHTTP(res, req)

// 	var actual = map[string]bool{}
// 	err = json.NewDecoder(res.Body).Decode(&actual)
// 	if err != nil {
// 		http.Error(res, err.Error(), 400)
// 		t.Fatalf("HTTP error")
// 		return nil
// 	}
// 	return actual
// }

/////// Original version
// func Test_Login(t *testing.T) {
// 	var path = "/login.html"
// 	form := url.Values{}
// 	form.Add("username", "asdf")
// 	form.Add("password", "asdf")
// 	//resp, err = http.PostForm(addr+path, form)
// 	req, err := http.NewRequest("POST", addr+path, strings.NewReader(form.Encode()))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req.PostForm = form
// 	res := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Login)
// 	handler.ServeHTTP(res, req)
// 	//CreateAccount(res, req)

// 	var expected = "login success"
// 	var actual string
// 	err = json.NewDecoder(res.Body).Decode(&actual)
// 	if err != nil {
// 		http.Error(res, err.Error(), 400)
// 		t.Fatalf("HTTP error")
// 		return
// 	}
// 	if expected != actual {
// 		t.Fatalf("Expected %s got %s", expected, actual)
// 	}
// }

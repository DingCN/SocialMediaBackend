package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/backend"
	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
	"github.com/DingCN/SocialMediaBackend/pkg/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var addr = "//127.0.0.1:8080"

func startBackend() {
	lis, err := net.Listen("tcp", backend.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	backend, _ := backend.New()
	s := grpc.NewServer()
	protocol.RegisterTwitterRPCServer(s, backend)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func startWeb() {
	backendAddr := "localhost:50051"
	conn, err := grpc.Dial(backendAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("web adn backend did not connect: %v", err)

	}
	webSrv.C = protocol.NewTwitterRPCClient(conn)
}

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

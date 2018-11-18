package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var addr = "//127.0.0.1:9090"

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
	var path = "/createPost.html"
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

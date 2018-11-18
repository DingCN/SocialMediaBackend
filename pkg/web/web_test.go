package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var addr = "//127.0.0.1:9090/createAccount.html"

func Test_CreateAccount(t *testing.T) {
	//var path = "/createAccount.html"
	form := url.Values{}
	form.Add("username", "asdf")
	form.Add("password", "asdf")
	//resp, err = http.PostForm(addr+path, form)
	req, err := http.NewRequest("POST", addr, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = form
	res := httptest.NewRecorder()

	CreateAccount(res, req)

	var expected = "create account success"
	var actual string
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

// func Test_CreateAccount(t *testing.T) {

// 	data := url.Values{}
// 	data.Set("username", "asdf")
// 	data.Add("password", "asdf")

// 	req, err := http.NewRequest("POST", "//127.0.0.1:9090/createAccount.html", strings.NewReader(data.Encode()))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	res := httptest.NewRecorder()

// 	CreateAccount(res, req)

// 	exp := "Hello World"
// 	act := res.Body.String()
// 	if exp != act {
// 		t.Fatalf("Expected %s gog %s", exp, act)
// 	}
// }

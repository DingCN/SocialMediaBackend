package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_CreateAccount(t *testing.T) {

	data := url.Values{}
	data.Set("username", "asdf")
	data.Add("password", "asdf")

	req, err := http.NewRequest("POST", "//127.0.0.1:9090/createAccount.html", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	CreateAccount(res, req)

	exp := "Hello World"
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected %s gog %s", exp, act)
	}
}

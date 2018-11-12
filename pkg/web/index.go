package web

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("username")
	fmt.Fprint(w, cookie)
	_, err := w.Write([]byte("UnLoggedInPage"))
	if err != nil {
		panic(err)
	}
}

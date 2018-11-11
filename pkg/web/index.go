package web

import (
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("Hello welcome to my website"))
	if err != nil {
		panic(err)
	}
}

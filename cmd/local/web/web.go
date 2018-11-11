package main

import (
	"final_project/pkg/web"
	"os"
)

func main() {
	cfg := &web.Config{
		Addr: os.Getenv("HOST"),
	}

	webSrv, err := web.New(cfg)
	if err != nil {
		panic(err)
	}

	err = webSrv.Start()
	if err != nil {
		panic(err)
	}
}

package main

import (
	"os"

	"github.com/DingCN/SocialMediaBackend/pkg/web"
)

func main() {
	cfg := &web.Config{
		Addr: os.Getenv("HOST"),
		// MaxFeedsNum: 3,
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

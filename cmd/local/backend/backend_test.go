package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/web"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
)

// pb "google.golang.org/grpc/examples/helloworld/helloworld"
func TestSignupRPC(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := protocol.SignupRequest{}
	request.Username = "user1"
	request.Password = "password1"
	reply, err := web.c.SignupRPC(ctx, &request)
	if err != nil {
		t.Fatalf("Signup incorrect")
	}
	reply, err := web.c.SignupRPC(ctx, &request)
	if err.Error() != errorcode.ErrUsernameTaken {
		t.Fatalf("Signup incorrect")
	}
}

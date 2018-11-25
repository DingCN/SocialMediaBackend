package main

import (
	"testing"

	"github.com/DingCN/SocialMediaBackend/pkg/backend"

	"github.com/DingCN/SocialMediaBackend/pkg/errorcode"
)

// pb "google.golang.org/grpc/examples/helloworld/helloworld"

func TestAddUser(t *testing.T) {
	backend, _ := backend.New()
	success, err := backend.Storage.AddUser("", "12345678")
	if success != false || err.Error() != errorcode.ErrInvalidUsername.Error() {
		t.Fatalf("AddUser incorrect")
	}
	success, err = backend.Storage.AddUser("asdfghjkl", "")
	if success != false || err.Error() != errorcode.ErrInvalidPassword.Error() {
		t.Fatalf("AddUser incorrect")
	}
	success, err = backend.Storage.AddUser("asdfghjkl", "12345678")
	if success != true || err != nil {
		t.Fatalf("AddUser incorrect")
	}
	success, err = backend.Storage.AddUser("asdfghjkl", "12345678")
	if success != false || err.Error() != errorcode.ErrUsernameTaken.Error() {
		t.Fatalf("AddUser incorrect")
	}
}

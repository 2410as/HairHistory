package request

import "net/http"

// CreateUser is an empty request for POST /api/users.
// (Anonymous user id generation has no input fields)
type CreateUser struct{}

func NewCreateUser(_ *http.Request) (*CreateUser, error) {
	return &CreateUser{}, nil
}


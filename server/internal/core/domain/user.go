package domain

import "github.com/google/uuid"

type User struct {
	Id          string
	Username    string
	Language    string
	CurrentRoom *Room
}

func NewUser(username, language string) *User {
	return &User{
		Id:          uuid.NewString(),
		CurrentRoom: nil,
		Username:    username,
		Language:    language,
	}
}

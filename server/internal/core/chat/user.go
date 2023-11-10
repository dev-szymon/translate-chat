package chat

import "github.com/google/uuid"

type User struct {
	Id          string
	Username    string
	Language    string
	CurrentRoom *Room
}

func NewUser() *User {
	return &User{
		Id: uuid.NewString(),
	}
}

func (u *User) UpdateUserDetails(username, language string) {
	u.Username = username
	u.Language = language
}

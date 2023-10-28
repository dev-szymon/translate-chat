package domain

import "sync"

type Chat struct {
	roomMu sync.Mutex
	Rooms  map[string]*Room

	usersMu sync.Mutex
	Users   map[*User]bool
}

func NewChat() *Chat {
	return &Chat{}
}

func (c *Chat) AddUser(u *User) {
	c.usersMu.Lock()
	defer c.usersMu.Unlock()
	c.Users[u] = true
}

func (c *Chat) AddRoom(r *Room) {
	c.roomMu.Lock()
	defer c.roomMu.Unlock()
	c.Rooms[r.Id] = r
}

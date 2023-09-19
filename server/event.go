package main

import "encoding/json"

const (
	errorEvent      = "error"
	joinRoom        = "join-room"
	userJoinedEvent = "user-joined"
	newMessageEvent = "new-message"
	leaveRoom       = "leave-room"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Message struct {
	Id          string  `json:"id"`
	Transcript  string  `json:"transcript"`
	Confidence  float32 `json:"confidence"`
	Translation *string `json:"translation"`
	SenderId    string  `json:"senderId"`
}

type CurrentRoom struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Users []*User `json:"users"`
}

type UserJoinedPayload struct {
	NewUser *User        `json:"newUser"`
	Room    *CurrentRoom `json:"room"`
}

type NewMessagePayload struct {
	Message *Message `json:"message"`
}

type LeaveRoomPayload struct {
	RoomId string `json:"roomId"`
}
type UserLeftRoomMessage struct {
	UserId string `json:"userId"`
}

type ErrorPayload struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type JoinRoomPayload struct {
	Username string `json:"username"`
	Language string `json:"language"`
	RoomId   string `json:"roomId"`
}

package websocket

import (
	"encoding/json"
)

const (
	ERROR_EVENT       = "error-event"
	JOIN_ROOM_EVENT   = "join-room-event"
	USER_JOINED_EVENT = "user-joined-event"
	NEW_MESSAGE_EVENT = "new-message-event"
	LEAVE_ROOM_EVENT  = "leave-room-event"
)

type InboundEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OutboundEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Message struct {
	Id          string    `json:"id"`
	Transcript  string    `json:"transcript"`
	Confidence  float32   `json:"confidence"`
	Translation *string   `json:"translation"`
	Sender      *RoomUser `json:"sender"`
}

type RoomUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Language string `json:"language"`
}
type CurrentRoom struct {
	Id    string      `json:"id"`
	Name  string      `json:"name"`
	Users []*RoomUser `json:"users"`
}

type UserJoinedPayload struct {
	NewUser *RoomUser    `json:"newUser"`
	Room    *CurrentRoom `json:"room"`
}

type NewMessagePayload struct {
	Message *Message `json:"message"`
}

type UserLeftRoomMessage struct {
	UserId string `json:"userId"`
}

type ErrorPayload struct {
	Error string `json:"error"`
}

type JoinRoomPayload struct {
	Username string `json:"username"`
	Language string `json:"language"`
	RoomId   string `json:"roomId"`
}

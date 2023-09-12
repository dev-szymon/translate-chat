package main

import "encoding/json"

const (
	errorMessage         = "error"
	joinRoom             = "join-room"
	userJoined           = "user-joined"
	joinRoomConfirmation = "join-room-confirmation"
	leaveRoom            = "leave-room"
	updateLanguage       = "update-language"
	translatedMessage    = "translated-message"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type JoinRoomPayload struct {
	Username string `json:"username"`
	Language string `json:"language"`
	RoomId   string `json:"roomId"`
}

type UserJoinedRoomPayload struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Language string `json:"language"`
	RoomId   string `json:"roomId"`
	RoomName string `json:"roomName"`
	Users    []User `json:"users"`
}

type LeaveRoomPayload struct {
	RoomId string `json:"roomId"`
}
type UserLeftRoomMessage struct {
	UserId string `json:"userId"`
}

type UpdateLanguagePayload struct {
	UserId   string `json:"userId"`
	Language string `json:"language"`
}

type TranslatedMessagePayload struct {
	Transcript  string  `json:"transcript"`
	Confidence  float32 `json:"confidence"`
	Translation *string `json:"translation"`
	UserId      string  `json:"userId"`
	RoomId      string  `json:"roomId"`
}

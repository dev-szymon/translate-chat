package main

import "github.com/dev-szymon/translate-chat/server/lib"

const updateLanguage = "update-language"
const joinRoom = "join-room"
const leaveRoom = "leave-room"

type Message struct {
	Type     string `json:"type"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
	Language string `json:"language"`
	Value    string `json:"value"`
}

type BroadcastTrasncript struct {
	FromUser   *User
	Transcript *lib.TranscriptResponse
}

package main

const updateLang = "update-lang"
const joinRoom = "join-room"
const leaveRoom = "leave-room"
const message = "message"

type Message struct {
	Type     string `json:"type"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
	Lang     string `json:"lang"`
	Value    string `json:"value"`
}

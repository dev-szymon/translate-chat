package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-szymon/translate-chat/server/lib"
	"golang.org/x/net/websocket"
)

type BroadcastMessage struct {
	FromUser    string  `json:"from_user"`
	Translation string  `json:"translation"`
	Transcript  string  `json:"transcript"`
	Confidence  float32 `json:"condifence"`
}

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	currentRoom *Room
	Language    string `json:"language"`
	conn        *websocket.Conn
	messageCh   chan []byte
}

func (usr *User) sendTranslation(broadcastTranscript *BroadcastTrasncript) {
	ctx := context.Background()
	br := &BroadcastMessage{
		Transcript: broadcastTranscript.Transcript.Transcript,
		Confidence: broadcastTranscript.Transcript.Confidence,
	}

	if broadcastTranscript.FromUser.Language != usr.Language {
		translation, err := lib.TranslateTranscript(ctx, broadcastTranscript.FromUser.Language, usr.Language, broadcastTranscript.Transcript.Transcript)
		if err != nil {
			fmt.Println("Error translating file.")
		}
		br.Translation = translation.Translation
	}

	bytes, err := json.Marshal(br)
	if err != nil {
		usr.conn.Write([]byte("Something went wrong"))
	}
	usr.conn.Write(bytes)
}

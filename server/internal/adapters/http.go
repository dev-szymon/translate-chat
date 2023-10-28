package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os/exec"

	"github.com/dev-szymon/translate-chat/server/internal/core/domain"
	"github.com/dev-szymon/translate-chat/server/internal/ports"
	"github.com/gorilla/websocket"
)

type httpResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) HandlePostFile(chat *domain.Chat, ts ports.TranslateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		if r.Method == "POST" {
			encoder := json.NewEncoder(w)
			roomId := r.FormValue("roomId")
			if roomId == "" {
				encoder.Encode(&httpResponse{Error: "room if missing"})
				return
			}
			userId := r.FormValue("userId")
			if userId == "" {
				encoder.Encode(&httpResponse{Error: "user not found"})
				return
			}
			var user *domain.User

			room, ok := chat.Rooms[roomId]
			if !ok {
				encoder.Encode(&httpResponse{Error: "room not found"})
				return
			}

			for id, u := range room.Users {
				if id == userId {
					user = u
					break
				}
			}
			if user == nil {
				encoder.Encode(&httpResponse{Error: "user not found"})
				return
			}

			file, _, err := r.FormFile("file")
			if err != nil {
				fmt.Println("Failed to parse form file:", err)
				encoder.Encode(&httpResponse{Error: "error parsing file"})
				return
			}
			defer file.Close()

			flacFile, err := convertFileToFlac(file)
			if err != nil {
				fmt.Println("Failed to convert audio file", err)
				encoder.Encode(&httpResponse{Error: "error"})
				return
			}

			transcript, err := ts.TranscribeAudio(ctx, user.Language, flacFile)
			if err != nil {
				fmt.Println("failed to translate audio file", err)
				encoder.Encode(&httpResponse{Error: "error"})
				return
			}
			user.CurrentRoom.BroadcastMessage(transcript, user)

			encoder.Encode(&httpResponse{Message: "ok"})
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// CheckOrigin:     func(r *http.Request) bool { return true },
}

func (s *Server) HandleWS(chat *domain.Chat, e ports.EventService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		encoder := json.NewEncoder(w)

		if err != nil {
			encoder.Encode(&httpResponse{Error: "room if missing"})
		}

		fmt.Println("New connection from: ", ws.RemoteAddr())

		for {
			_, msgBytes, err := ws.ReadMessage()
			if err != nil {
				if err == io.EOF {
					break
				}
				b := e.EncodeEvent(errorEvent, &ErrorPayload{Message: "Unable to read from buffer", Error: err.Error()})
				ws.WriteMessage(websocket.BinaryMessage, b)
				continue
			}

			var event Event
			err = json.Unmarshal(msgBytes, &event)
			if err != nil {
				b := e.EncodeEvent(errorEvent, &ErrorPayload{Message: "Unable to unmarshal event payload", Error: err.Error()})
				ws.WriteMessage(websocket.BinaryMessage, b)
			}

			switch event.Type {
			case joinRoom:
				var payload JoinRoomPayload
				var room *domain.Room
				err := json.Unmarshal(event.Payload, &payload)
				if err != nil {
					b := e.EncodeEvent(errorEvent, &ErrorPayload{Message: "Unable to unmarshal join-room payload", Error: err.Error()})
					ws.WriteMessage(websocket.BinaryMessage, b)
				}

				room = chat.Rooms[payload.RoomId]
				if room == nil {
					for _, r := range chat.Rooms {
						if r.Name == payload.RoomId {
							room = r
							break
						}
					}
				}
				if room == nil {
					room = domain.NewRoom()
				}
				user := domain.NewUser(payload.Username, payload.Language)
				room.AddUser(user)
			case leaveRoom:
				// user.currentRoom.leaveCh <- user
			}
		}
	}
}

func convertFileToFlac(file multipart.File) ([]byte, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
		return nil, err
	}

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-c:a", "flac", "-f", "flac", "-")

	var (
		output bytes.Buffer
		errors bytes.Buffer
	)

	cmd.Stdin = bytes.NewReader(b)
	cmd.Stdout = &output
	cmd.Stderr = &errors

	err = cmd.Run()
	if err != nil {
		log.Printf("FFmpeg stderr: %s", errors.String())
		return nil, err
	}

	return output.Bytes(), nil
}

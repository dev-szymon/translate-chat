package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dev-szymon/translate-chat/server/handlers"
	"github.com/dev-szymon/translate-chat/server/lib"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func enableCors(handler http.HandlerFunc) http.HandlerFunc {
	webUrl := os.Getenv("WEB_URL")
	if webUrl == "" {
		log.Fatal("WEB_URL environmental variable missing.")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		handler(w, r)
	}
}

type googleAppCredentials struct {
	ProjectID string `json:"project_id"`
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("not able to load env vile")
	}
	credsFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFilePath == "" {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS enviromental variable missing")
	}
	credsFile, err := os.Open(credsFilePath)
	if err != nil {
		return err

	}
	c, err := io.ReadAll(credsFile)
	if err != nil {
		return err
	}
	var googleAppCreds googleAppCredentials
	err = json.Unmarshal(c, &googleAppCreds)
	if err != nil {
		return err
	}
	os.Setenv("GOOGLE_PROJECT_ID", googleAppCreds.ProjectID)
	return nil
}

type Room struct {
	id          string
	users       map[*User]bool
	broadcastCh chan []byte
	joinCh      chan *User
	leaveCh     chan *User
}

func (r *Room) run() {
	fmt.Println("new room started")
	for {
		select {
		case user := <-r.joinCh:
			r.users[user] = true
			user.currentRoom = r
			fmt.Println("new user joined the room: ", user.conn.RemoteAddr())
		case msg := <-r.broadcastCh:
			fmt.Println(string(msg))
			signature := msg[:1]
			fmt.Println(signature)
			fmt.Println(string(signature))
			for usr := range r.users {
				usr.conn.Write(msg)
			}
		case usr := <-r.leaveCh:
			delete(r.users, usr)
			usr.currentRoom = nil
		}
	}

}

type Server struct {
	rooms      map[string]*Room
	users      map[*User]bool
	register   chan *User
	unregister chan *User
}

func (s *Server) newRoom() *Room {
	id := lib.GeneratePseudoID()

	r := &Room{
		id:          id,
		users:       make(map[*User]bool),
		broadcastCh: make(chan []byte),
		joinCh:      make(chan *User),
		leaveCh:     make(chan *User),
	}

	s.rooms[id] = r
	go r.run()
	return r
}

func (s *Server) start() {
	for {
		select {
		case user := <-s.register:
			s.users[user] = true
		case user := <-s.unregister:
			if _, ok := s.users[user]; ok {
				for _, r := range s.rooms {
					delete(r.users, user)
				}
				delete(s.users, user)
				close(user.messageCh)
			}
		}
	}
}

func (s *Server) serveWS(conn *websocket.Conn) {
	user := &User{
		server:      s,
		conn:        conn,
		currentRoom: nil,
		messageCh:   make(chan []byte, 256),
	}
	s.register <- user
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			conn.Write([]byte("wrong message format"))
			continue
		}
		msgBytes := buf[:n]
		msg := string(msgBytes)

		switch true {
		case msg == "new-room":
			room := s.newRoom()
			room.joinCh <- user
			conn.Write([]byte(fmt.Sprintf("Created new room: %s", room.id)))
		case strings.HasPrefix(msg, "join-room/"):
			roomId, ok := strings.CutPrefix(msg, "join-room/")
			if !ok {
				conn.Write([]byte("wrong message format"))
			}
			var foundRoom *Room
			for id := range s.rooms {
				if id == roomId {
					foundRoom = s.rooms[id]
					break
				}
			}
			if foundRoom == nil {
				conn.Write([]byte(fmt.Sprintf("Created new room: %s", foundRoom.id)))
			} else {
				foundRoom.joinCh <- user
			}
		default:
			user.currentRoom.broadcastCh <- msgBytes
		}
	}

}

type User struct {
	server      *Server
	currentRoom *Room
	conn        *websocket.Conn
	messageCh   chan []byte
}

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal("Not able to load necessary environmental variables: ", err)
	}

	s := &Server{rooms: make(map[string]*Room),
		users:      make(map[*User]bool),
		register:   make(chan *User),
		unregister: make(chan *User),
	}

	go s.start()
	http.Handle("/ws", websocket.Handler(s.serveWS))
	http.HandleFunc("/translate-file", enableCors(handlers.HandleTranslateFile))

	fmt.Println("Translation service starting on: ", "http://localhost:8055")
	log.Fatal(http.ListenAndServe(":8055", nil))
}

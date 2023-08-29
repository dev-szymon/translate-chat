package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func enableCors(handler http.HandlerFunc) http.HandlerFunc {
	webUrl := os.Getenv("WEB_URL")
	if webUrl == "" {
		log.Fatal("WEB_URL environmental variable missing.")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", webUrl)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		handler(w, r)
	}
}

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal("Not able to load necessary environmental variables: ", err)
	}

	s := NewServer()

	http.Handle("/ws", websocket.Handler(s.serveWS))
	http.HandleFunc("/translate-file", enableCors(s.handleTranslateFile()))

	fmt.Println("Server starting on: ", "http://localhost:8055")
	log.Fatal(http.ListenAndServe(":8055", nil))
}

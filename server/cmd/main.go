package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dev-szymon/translate-chat/server/internal/adapters"
	"github.com/dev-szymon/translate-chat/server/internal/core"
	"github.com/dev-szymon/translate-chat/server/internal/core/domain"
)

func main() {
	err := core.LoadEnv()
	if err != nil {
		log.Fatal("Not able to load necessary environmental variables: ", err)
	}

	ts := adapters.NewTranslateService()
	chat := domain.NewChat()
	server := adapters.NewServer()
	es := adapters.NewEventService()

	http.HandleFunc("/ws", server.HandleWS(chat, es))
	http.HandleFunc("/translate-file", adapters.EnableCors(server.HandlePostFile(chat, ts)))

	fmt.Println("Server starting on: ", "http://localhost:8055")
	log.Fatal(http.ListenAndServe(":8055", nil))
}

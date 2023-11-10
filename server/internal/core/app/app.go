package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dev-szymon/translate-chat/server/internal/ports"
)

type App struct {
	ts  ports.TranslateServicePort
	wss ports.WebsocketServerPort
}

func NewApp(ts ports.TranslateServicePort, wss ports.WebsocketServerPort) *App {
	return &App{
		ts:  ts,
		wss: wss,
	}
}

func (a *App) Run() {
	http.HandleFunc("/ws", a.wss.HandleWS(a.ts))

	fmt.Println("Server starting on: ", "http://localhost:8055")
	log.Fatal(http.ListenAndServe(":8055", nil))
}

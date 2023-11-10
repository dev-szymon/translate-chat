package main

import (
	"fmt"

	"github.com/dev-szymon/translate-chat/server/internal/adapters/translate"
	"github.com/dev-szymon/translate-chat/server/internal/adapters/websocket"
	"github.com/dev-szymon/translate-chat/server/internal/core/app"
	"github.com/dev-szymon/translate-chat/server/internal/core/config"
	"github.com/dev-szymon/translate-chat/server/internal/ports"
)

func main() {
	c := config.MustLoadEnv()

	var ts ports.TranslateServicePort
	if c.Environment == "TEST" {
		fmt.Println("Translation service running in debug mode...")
		ts = translate.NewDebugTranslateService()
	} else {
		ts = translate.NewTranslateService()
	}
	wss := websocket.NewServer()

	app := app.NewApp(ts, wss)

	app.Run()
}

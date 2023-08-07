package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dev-szymon/translate-chat/server/handlers"
	"github.com/joho/godotenv"
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

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal("Not able to load necessary environmental variables: ", err)
	}

	http.HandleFunc("/translate-file", enableCors(handlers.HandleTranslateFile))

	fmt.Println("Translation service starting on: ", "http://localhost:8055")
	log.Fatal(http.ListenAndServe(":8055", nil))
}

package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
)

type googleAppCredentials struct {
	ProjectID string `json:"project_id"`
}

const (
	DEFAULT_WEB_URL                        = "http://localhost:5173"
	DEFAULT_GOOGLE_APPLICATION_CREDENTIALS = "./google_application_credentials.json"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("not able to load env vile")
	}

	webUrl := os.Getenv("WEB_URL")
	if webUrl == "" {
		fmt.Printf("WEB_URL environmental variable missing. Using default: %s\n", DEFAULT_WEB_URL)
		webUrl = DEFAULT_WEB_URL
	}

	credsFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFilePath == "" {
		fmt.Printf("GOOGLE_APPLICATION_CREDENTIALS enviromental variable missing. Using default %s\n", DEFAULT_GOOGLE_APPLICATION_CREDENTIALS)
		credsFilePath = DEFAULT_GOOGLE_APPLICATION_CREDENTIALS
	}
	credsFile, err := os.Open(credsFilePath)
	if err != nil {
		return err
	}
	defer credsFile.Close()

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

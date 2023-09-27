package main

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

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("not able to load env vile")
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

package config

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
	DEFAULT_ENV                            = "TEST"
)

type Config struct {
	Environment string
}

func MustLoadEnv() *Config {
	config := &Config{}

	err := godotenv.Load()
	if err != nil {
		panic("not able to load env file")
	}

	environmentEnv := os.Getenv("ENV")
	if environmentEnv == "" {
		fmt.Printf("ENV environmental variable missing. Using default %s\n", DEFAULT_ENV)
		os.Setenv("ENV", DEFAULT_ENV)
		config.Environment = DEFAULT_ENV
	}
	config.Environment = environmentEnv

	webUrl := os.Getenv("WEB_URL")
	if webUrl == "" {
		fmt.Printf("WEB_URL environmental variable missing. Using default: %s\n", DEFAULT_WEB_URL)
		webUrl = DEFAULT_WEB_URL
	}

	if config.Environment != "TEST" {
		credsFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if credsFilePath == "" {
			fmt.Printf("GOOGLE_APPLICATION_CREDENTIALS environmental variable missing. Using default %s\n", DEFAULT_GOOGLE_APPLICATION_CREDENTIALS)
			credsFilePath = DEFAULT_GOOGLE_APPLICATION_CREDENTIALS
		}
		credsFile, err := os.Open(credsFilePath)
		if err != nil {
			panic("google app credentials file not found")
		}
		defer credsFile.Close()

		c, err := io.ReadAll(credsFile)
		if err != nil {
			panic("could not read google app credentials file")
		}
		var googleAppCreds googleAppCredentials
		err = json.Unmarshal(c, &googleAppCreds)
		if err != nil {
			panic("could not unmarshal google app credentials")
		}
		os.Setenv("GOOGLE_PROJECT_ID", googleAppCreds.ProjectID)
	}

	return config
}

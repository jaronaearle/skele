package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func GetConfig() (config Config, err error) {
	bt, err := getEnvVar("BOT_TOKEN")
	if err != nil {
		fmt.Printf("GetConfig: Error loading bot token: %v", err)
		return
	}
	token := fmt.Sprintf("Bot %s", bt)
	config = Config{BotToken: token}
	return
}

func getEnvVar(k string) (v string, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("getEnvVar: Error loading env vars from .env: %v\n", err)
		return
	}
	v = os.Getenv(k)
	return
}

package config

import (
	"fmt"
	"os"

	"github.com/honeybadger-io/honeybadger-go"
	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func GetConfig() (cfg Config, err error) {
	bt, err := getEnvVar("BOT_TOKEN")
	if err != nil {
		honeybadger.Notify("GetConfig: Error loading bot token: %v", err)
		return
	}

	hb, err := getEnvVar("HONEY_BADGER_API_KEY")
	if err != nil {
		honeybadger.Notify("GetConfig: Error loading honeybadger api key: %v", err)
		return
	}

	honeybadger.Configure(honeybadger.Configuration{APIKey: hb})

	token := fmt.Sprintf("Bot %s", bt)
	cfg = Config{BotToken: token}

	return
}



func getEnvVar(k string) (v string, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		honeybadger.Notify("getEnvVar: Error loading env vars from .env: %v\n", err)
		return
	}
	v = os.Getenv(k)

	return
}


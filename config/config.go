package config

import (
	"fmt"
	"os"

	"github.com/honeybadger-io/honeybadger-go"
	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	PaperTrailHost string
}

func GetConfig() (cfg Config, err error) {
	bt, err := getEnvVar("BOT_TOKEN")
	if err != nil {
		return
		}

	hb, err := getEnvVar("HONEY_BADGER_API_KEY")
	if err != nil {
		return
		}

	pt, err := getEnvVar("PAPER_TRAIL_HOST")
	if err != nil {
		return
		}

	honeybadger.Configure(honeybadger.Configuration{APIKey: hb})

	botToken := fmt.Sprintf("Bot %s", bt)
	cfg = Config{BotToken: botToken, PaperTrailHost: pt}

	return
}

func getEnvVar(k string) (v string, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		err = fmt.Errorf("getEnvVar: Error pulling %v from .env", k)
			}
	v = os.Getenv(k)

	return
}


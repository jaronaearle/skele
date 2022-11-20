package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Session *discordgo.Session
}

func NewDiscordBot(session *discordgo.Session) *DiscordBot {
	return &DiscordBot{
		Session: session,
	}
}

func (db *DiscordBot) StartBotSession() {
	err := db.Session.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Captain Hook is running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	defer db.Session.Close()
}
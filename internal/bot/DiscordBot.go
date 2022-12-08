package bot

import (
	"fmt"

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

func (db *DiscordBot) RegisterHandlers() {
	db.Session.AddHandler(messageCreate)

	db.Session.Identify.Intents = discordgo.IntentGuildMessages
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// TODO: access channel messages for chat op things
	if m.Content == "" {
		msgs, err := s.ChannelMessages(m.ChannelID, 1, "", "", m.ID)
		if err != nil {
			fmt.Printf("messageCreate: Error getting channel messages: %v", err)
		}
		fmt.Println(msgs)
	}
}

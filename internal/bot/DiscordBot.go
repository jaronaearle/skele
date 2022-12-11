package bot

import (
	"fmt"
	"skele/internal/data"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Exp bool
	Session *discordgo.Session
}

func NewDiscordBot(session *discordgo.Session, exp bool) *DiscordBot {
	return &DiscordBot{
		Exp: exp,
		Session: session,
	}
}

func (db *DiscordBot) RegisterHandlers() {
	db.Session.AddHandler(messageCreate)

	db.Session.Identify.Intents = discordgo.IntentGuildMessages
}

func (db *DiscordBot) SendEmbedMessage(m *discordgo.MessageEmbed, id string) (err error) {
	var channel string
	if db.Exp {
		channel = data.ChannelIDs.Exp
	} else {
		channel = id
	}

	_, err = db.Session.ChannelMessageSendEmbed(channel, m)
	if err != nil {
		fmt.Printf("SendMessageEmbed: Error sending embed message: %v\n", err)
		return
	}
	
	fmt.Printf("Embed message sent to %s at %v", channel, time.Now().Local())

	return
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

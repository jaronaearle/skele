package handlers

import (
	"fmt"
	"skele/internal/bot"
	"skele/internal/data"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	WORDLE_URL     = "https://www.nytimes.com/games/wordle/index.html"
	WORDLE_IMG_URL = "https://assets-prd.ignimgs.com/2022/04/15/wordle-1650045194490.jpg"
	FHP_URL        = "https://fuckinghomepage.com/"
)

type ScheduledMessageHandler struct {
	DiscordBot *bot.DiscordBot
}

func (s *ScheduledMessageHandler) SendMessage(m *discordgo.MessageEmbed, id string) {
	err := s.DiscordBot.SendEmbedMessage(m, id)
	if err != nil {
		fmt.Printf("ScheduledMessageHandler: Error sending embed message: %v\n", err)
	}

	fmt.Printf("SendMessage: sent %v message - %v\n", m.Title, time.Now().Local())
}

func (s *ScheduledMessageHandler) PrepareWordleMessage() (m *discordgo.MessageEmbed, id string) {
	m = buildMessageEmbed(WORDLE_URL, WORDLE_IMG_URL, "Time to Wordle!")
	id = data.ChannelIDs.General

	return
}

func (s *ScheduledMessageHandler) PrepareFHPMessage() (m *discordgo.MessageEmbed, id string) {
	m = buildMessageEmbed(FHP_URL, "", "Learn some shit today")
	id = data.ChannelIDs.Exp
	
	return
}

// Will probably want to break embeds out into their own templates.
func buildMessageEmbed(url, imgUrl, title string) (m *discordgo.MessageEmbed) {
	m = &discordgo.MessageEmbed{
		URL:   url,
		Title: title,
		Image: &discordgo.MessageEmbedImage{
			URL: imgUrl,
		},
	}

	return
}

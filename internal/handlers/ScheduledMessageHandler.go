package handlers

import (
	"fmt"
	"skele/internal/bot"
	"skele/internal/data"

	"github.com/bwmarrin/discordgo"
)

const (
	WORDLE_URL        = "https://www.nytimes.com/games/wordle/index.html"
	FHP_URL           = "https://fuckinghomepage.com/"
)

type ScheduledMessageHandler struct {
	DiscordBot *bot.DiscordBot
}

func (s *ScheduledMessageHandler) SendMessage(m *discordgo.MessageEmbed, id string) {
	// TODO: add dev flag to change channel ids when testing
	_, err := s.DiscordBot.Session.ChannelMessageSendEmbed(id, m)
	if err != nil {
		fmt.Printf("ScheduledMessageHandler: Error sending embed message: %v\n", err)
	}
	fmt.Println("sent scheduled msg")
}

func (s *ScheduledMessageHandler) PrepareWordleMessage() (m *discordgo.MessageEmbed, id string) {
	m = buildMessageEmbed(WORDLE_URL, "Time to Wordle!")
	id = data.EXP_CHANNEL_ID
	return
}

func (s *ScheduledMessageHandler) PrepareFHPMessage() (m *discordgo.MessageEmbed, id string) {
	m = buildMessageEmbed(FHP_URL, "Learn some shit today")
	id = data.EXP_CHANNEL_ID
	return
}

func buildMessageEmbed(url, title string) (m *discordgo.MessageEmbed) {
	m = &discordgo.MessageEmbed{
		URL:   url,
		Title: title,
	}

	return
}

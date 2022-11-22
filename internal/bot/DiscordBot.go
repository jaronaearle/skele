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

func (db *DiscordBot) RegisterHandlers() {
	db.Session.AddHandler(messageCreate)

	db.Session.Identify.Intents = discordgo.IntentGuildMessages
}

func (db *DiscordBot) StartBotSession() {
	err := db.Session.Open()
	if err != nil {
		fmt.Println("StartBotSession: Error opening discord bot connection:", err)
		return
	}

	fmt.Println("Captain Hook is running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	defer db.Session.Close()
}

func messageCreate (s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	emb := &discordgo.MessageEmbed{
		URL: "https://utahavalanchecenter.org/sites/default/files/forecast/202211/20221119-061743-6.png",
		Title: "Avy Report: November 20, 2022",
		Description: "Generally safe avalanche conditions.\nSmall avalanches are possible in isolated areas or extreme terrain.\n Remember that risk is inherent in mountain travel.\n\nWe are transitioning back to intermittent updates.  The next update will be on Wednesday.",
		Image: &discordgo.MessageEmbedImage{
			URL: "https://utahavalanchecenter.org/sites/default/files/forecast/202211/20221119-061743-6.png",
			Width: 50,
			Height: 50,
		},
	}

	if m.Content == "" {
		msgs, err := s.ChannelMessages(m.ChannelID, 1, "", "", m.ID)
		if err != nil {
			fmt.Printf("messageCreate: Error getting channel messages: %v", err)
		}
		fmt.Println(msgs)

	}

	if m.Content == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Hello from skele")
		s.ChannelMessageSendEmbed(m.ChannelID, emb)
		if err != nil {
			fmt.Printf("messageCreate: Error sending message: %v", err)
		}
	}
}

package handlers

import (
	"fmt"
	"skele/internal/bot"
	"skele/internal/crawlers"
	"skele/internal/data"

	"github.com/bwmarrin/discordgo"
)

const (
	BASE_URL             = "https://utahavalanchecenter.org/"
	FORECAST             = "forecast/salt-lake"
	SKI_PEEPS_CHANNEL_ID = "894699071535718401"
	EXP_CHANNEL_ID       = "899363073989689394"
)

type AvyReportHandler struct {
	AvyCrawler *crawlers.AvyCrawler
	DiscordBot *bot.DiscordBot
}

func (a *AvyReportHandler) SendAvyReport() {
	rp, _ := a.getAvyReport()
	em := a.buildReportEmbed(rp)

	// TODO: add dev flag to change channel ids when testing
	_, err := a.DiscordBot.Session.ChannelMessageSendEmbed(EXP_CHANNEL_ID, em)
	if err != nil {
		fmt.Printf("SendAvyReport: Error sending embed message: %v\n", err)
	}
	fmt.Println("sent report msg")
}

func (a *AvyReportHandler) getAvyReport() (rp data.AvyReport, err error) {
	rp, err = a.AvyCrawler.GetReport()
	if err != nil {
		fmt.Println("getAvyReport: Error calling GetReport: ", err)
		return
	}
	return
}

func (a *AvyReportHandler) buildReportEmbed(rp data.AvyReport) (em *discordgo.MessageEmbed) {
	url := fmt.Sprintf("%s%s", BASE_URL, FORECAST)
	title := fmt.Sprintf("Avy Report: %s", rp.Date)
	imgUrl := fmt.Sprintf("%s%s", BASE_URL, rp.ImageUrl)

	em = &discordgo.MessageEmbed{
		URL:         url,
		Title:       title,
		Description: rp.Details,
		Image: &discordgo.MessageEmbedImage{
			URL:    imgUrl,
			Width:  50,
			Height: 50,
		},
	}
	return
}

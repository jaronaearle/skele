package handlers

import (
	"fmt"
	"skele/internal/bot"
	"skele/internal/crawlers"
	"skele/internal/data"

	"github.com/bwmarrin/discordgo"
)

const (
	BASE_URL = "https://utahavalanchecenter.org/"
	FORECAST = "forecast/salt-lake"
)

type AvyReportHandler struct {
	AvyCrawler *crawlers.AvyCrawler
	DiscordBot *bot.DiscordBot
}

func (a *AvyReportHandler) SendAvyReport() {
	rp, _ := a.getAvyReport()

	em := a.buildReportEmbed(rp)

	_, err := a.DiscordBot.Session.ChannelMessageSendEmbed("ski peeps chan id", em)
	if err != nil {
		fmt.Printf("reportMessageCreate: Error sending embed message: %v\n", err)
	}
}

func(a *AvyReportHandler) getAvyReport() (rp data.AvyReport, err error) {
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
		URL: url,
		Title: title,
		Description: rp.Details,
		Image: &discordgo.MessageEmbedImage{
			URL: imgUrl,
			Width: 50,
			Height: 50,
		},
	}
	return
}

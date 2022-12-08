package handlers

import (
	"fmt"
	"skele/internal/bot"
	"skele/internal/crawlers"
	"skele/internal/data"
	"time"

	"github.com/bwmarrin/discordgo"
)

type AvyReportHandler struct {
	AvyCrawler *crawlers.AvyCrawler
	DiscordBot *bot.DiscordBot
}

func (a *AvyReportHandler) SendAvyReport() {
	rp, _ := a.getAvyReport()
	em := buildReportEmbed(rp)

	// TODO: add dev flag to change channel ids when testing
	_, err := a.DiscordBot.Session.ChannelMessageSendEmbed(data.ChannelIDs.SkiPeeps, em)
	if err != nil {
		fmt.Printf("SendAvyReport: Error sending embed message: %v\n", err)
	}
	fmt.Printf("SendAvyReport: sent report message - %v\n", time.Now())
}

func (a *AvyReportHandler) getAvyReport() (rp data.AvyReport, err error) {
	rp, err = a.AvyCrawler.GetReport()
	if err != nil {
		fmt.Println("getAvyReport: Error calling GetReport: ", err)
		return
	}

	return
}

func buildReportEmbed(rp data.AvyReport) (em *discordgo.MessageEmbed) {
	url := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, data.AvyUrlPaths.Forecast)
	title := fmt.Sprintf("Avy Report: %s", rp.Date)
	imgUrl := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, rp.ImageUrl)

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

func (a *AvyReportHandler) SendTodaysAvyList() {
	_, err := a.AvyCrawler.GetTodaysAvyList()
	if err != nil {
		fmt.Println("SendTodaysAvyList: Error: ", err)
		return
	}
}

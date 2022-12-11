package handlers

import (
	"fmt"
	"skele/internal/bot"
	"skele/internal/crawlers"
	"skele/internal/data"
	"time"

	"github.com/bwmarrin/discordgo"
)
type AvyCrawlerHandler struct {
	AvyCrawler *crawlers.AvyCrawler
	DiscordBot *bot.DiscordBot
}

func (a *AvyCrawlerHandler) SendAvyReport() {
	rp, _ := a.getAvyReport()
	em := buildReportEmbed(rp)

	err := a.DiscordBot.SendEmbedMessage(em, data.ChannelIDs.SkiPeeps)
	if err != nil {
		fmt.Printf("SendAvyReport: Error sending embed message: %v\n", err)
		return
	}
	fmt.Printf("SendAvyReport: sent report message - %v\n", time.Now())
}

func (a *AvyCrawlerHandler) getAvyReport() (rp data.AvyReport, err error) {
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

func (a *AvyCrawlerHandler) SendTodaysAvyList() {
	avs, date, _ := a.getAvyList()
	em := buildTodaysAvyEmbed(avs, date)

	err := a.DiscordBot.SendEmbedMessage(em, data.ChannelIDs.SkiPeeps)
	if err != nil {
		fmt.Printf("SendAvyReport: Error sending embed message: %v\n", err)
		return
	}

	fmt.Printf("SendAvyReport: sent report message - %v\n", time.Now())
}

func (a *AvyCrawlerHandler) getAvyList() (avs []data.Avy, date string, err error) {
	avs, date, err = a.AvyCrawler.GetTodaysAvyList()
	if err != nil {
		fmt.Println("getAvyList: Error: ", err)
		return
	}

	return
}

func buildTodaysAvyEmbed(avs []data.Avy, date string) (em *discordgo.MessageEmbed) {
	var list string

	for _, a := range avs {
		list = fmt.Sprintf("%s\n%s\n%s\n", list, fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, a.Url), a.Title)
	}
	fmt.Println(avs)
	fmt.Println(list)

	em = &discordgo.MessageEmbed{
		// URL:         url,
		Title:       "Avys today",
		Description: list,
	}

	return
}

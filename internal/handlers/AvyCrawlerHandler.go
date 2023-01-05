package handlers

import (
	"fmt"
	"log"
	"skele/internal/bot"
	"skele/internal/crawlers"
	"skele/internal/data"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	slc =  "Salt Lake"
	ogden = "Ogden"
	provo = "Provo"
	uintas = "Uintas"
	moab = "Moab"
)
type AvyCrawlerHandler struct {
	AvyCrawler *crawlers.AvyCrawler
	DiscordBot *bot.DiscordBot
}

func (a *AvyCrawlerHandler) SendAvyReport() {
	rp, err := a.getAvyReport()
	if err != nil {
		err = fmt.Errorf("SendAvyReport: %w", err)
		log.Println(err)

		return
	}

	em := buildReportEmbed(rp)

	err = a.DiscordBot.SendEmbedMessage(em, data.ChannelIDs.SkiPeeps)
	if err != nil {
		err = fmt.Errorf("SendAvyReport: Error sending embed message: %w", err)
		log.Println(err)

		return
	}
	log.Println("SendAvyReport: sent report message")
}

func (a *AvyCrawlerHandler) getAvyReport() (rp data.AvyReport, err error) {
	rp, err = a.AvyCrawler.GetReport()
	if err != nil {
		err = fmt.Errorf("getAvyReport: Error calling GetReport: %w", err)
		log.Println(err)

		return
	}

	return
}

func buildReportEmbed(rp data.AvyReport) (em *discordgo.MessageEmbed) {
	url := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, data.AvyUrlPaths.Forecast)
	title := fmt.Sprintf("Avy Report: %s", rp.Date)
	description := fmt.Sprintf("%s\n\n%s", rp.Details, rp.SpecialBulletin)
	imgUrl := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, rp.ImageUrl)

	log.Printf("buildReportEmbed: Report: %v - %v\n", rp, time.Now())

	em = &discordgo.MessageEmbed{
		URL:         url,
		Title:       title,
		Description: description, 
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
	if len(avs) == 0 {
		log.Println("SendTodaysAvyList: no avys today, exiting")

		return
	}

	em := buildTodaysAvyEmbed(avs, date)

	err := a.DiscordBot.SendEmbedMessage(em, data.ChannelIDs.SkiPeeps)
	if err != nil {
		err = fmt.Errorf("SendAvyReport: Error sending embed message: %w", err)
		log.Println(err)

		return
	}

	log.Println("SendAvyReport: sent report message")
}

func (a *AvyCrawlerHandler) getAvyList() (avs []data.Avy, date string, err error) {
	avs, date, err = a.AvyCrawler.GetTodaysAvyList()
	if err != nil {
		err = fmt.Errorf("getAvyList: Error: %w", err)
		
		return
	}

	return
}

func buildTodaysAvyEmbed(avs []data.Avy, date string) (em *discordgo.MessageEmbed) {
	title := fmt.Sprintf("Avalanche Activity %s", date)

	slcList := filterByRegion(avs, slc)
	ogList := filterByRegion(avs, ogden)
	prList := filterByRegion(avs, provo)
	uList := filterByRegion(avs, uintas)
	mList := filterByRegion(avs, moab)

	slcAvy := formatByRegion(slcList, slc)
	ogAvy := formatByRegion(ogList, ogden)
	prAvy := formatByRegion(prList, provo)
	uAvy := formatByRegion(uList, uintas)
	mAvy := formatByRegion(mList, moab)

	content := fmt.Sprintf("%s\n%s%s%s%s", slcAvy, ogAvy, prAvy, uAvy, mAvy)

	em = &discordgo.MessageEmbed{
		Title:       title,
		Description: content,
	}

	return
}

func filterByRegion(avs []data.Avy, r string) (fAvs []data.Avy) {
	for _, a := range avs {
		if strings.EqualFold(a.Region, r) {
			fAvs = append(fAvs, a)
		}
	}
	
	return
}

func formatByRegion(avs []data.Avy, r string) (regionAvys string) {
	if len(avs) == 0 {
		log.Printf("formatByRegion: no avy activity for given region: %s\n", r)

		return
	}

	var avys string

	for _, a := range avs {
		avys = fmt.Sprintf("%s\n%s\n%s\n", avys, a.Title, fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, a.Url))
	}

	regionAvys = fmt.Sprintf("%s\n%s\n", r, avys)

	return
}

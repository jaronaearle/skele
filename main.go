package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"skele/config"
	"skele/internal/bot"
	"skele/internal/crawlers"
	"skele/internal/handlers"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
)

var (
	avyCenterDomains = []string{"https://utahavalanchecenter.org/", "https://utahavalanchecenter.org", "utahavalanchecenter.org/", "www.utahavalanchecenter.org/", "utahavalanchecenter.org"}
)

type Payload struct {
	Content string `json:"content"`
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("main: Error initializing config: %v", err)
	}

	c := colly.NewCollector(colly.AllowedDomains(avyCenterDomains...))
	ac := crawlers.NewAvyCrawler(c)

	session, err := discordgo.New(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	bot := bot.NewDiscordBot(session)
	avyReportHandler := handlers.AvyReportHandler{
		AvyCrawler: ac,
		DiscordBot: bot,
	}

	ctx, cancel := context.WithCancel((context.Background()))

	go startCron(ctx, avyReportHandler, cancel)
	go startBot(ctx, bot, cancel)

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, os.Interrupt)

	<-sig

	// TODO: make contexts/cancel work
	// select {
	// case <-sig:
	// case <-ctx.Done():
	// 	return
	// }
	// cancel()
	fmt.Println("Gracefully returning to the grave...")
}

func startCron(pCtx context.Context, a handlers.AvyReportHandler, exit context.CancelFunc) {
	fmt.Println("Starting cron...")
	defer exit()
	defer fmt.Println("Exiting cron...")

	mtnTZ, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Printf("Error loading america/denver tz: %s", err)
		return
	}
	s := gocron.NewScheduler(mtnTZ)

	s.Every(1).Days().At("07:30").Do(func() {
		a.SendAvyReport()
	})

	s.StartBlocking()

	// TODO: make contexts work
	// <- pCtx.Done()
}

func startBot(pCtx context.Context, bot *bot.DiscordBot, exit context.CancelFunc) {
	fmt.Println("Starting bot session...")
	defer exit()
	defer fmt.Println("Exiting bot session...")

	bot.RegisterHandlers()

	err := bot.Session.Open()
	if err != nil {
		fmt.Printf("startBot: Error opening websocket connection: %v", err)
		return
	}
}

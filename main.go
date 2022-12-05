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
type Payload struct {
	Content string `json:"content"`
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("main: Error initializing config: %v", err)
	}

	c := colly.NewCollector(colly.AllowedDomains(crawlers.AvyCenterDomains...))
	ac := crawlers.NewAvyCrawler(c)

	session, err := discordgo.New(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	bot := bot.NewDiscordBot(session)

	h := Handlers{
		AvyReportHandler: handlers.AvyReportHandler{
			AvyCrawler: ac,
			DiscordBot: bot,
		},
		ScheduledMessageHandler: handlers.ScheduledMessageHandler{
			DiscordBot: bot,
		},
	}

	ctx, cancel := context.WithCancel((context.Background()))

	go startCron(ctx, h, cancel)
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

func startCron(pCtx context.Context, h Handlers, exit context.CancelFunc) {
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
		h.AvyReportHandler.SendAvyReport()
	})

	// TODO: dev flag
	// skip while developing
	// s.Every(1).Days().At("11:00").Do(func() {
	// 	m, id := h.ScheduledMessageHandler.PrepareWordleMessage()
	// 	h.ScheduledMessageHandler.SendMessage(m, id)
	// })

	// s.Every(1).Days().At("09:30").Do(func() {
	// 	m, id := h.ScheduledMessageHandler.PrepareFHPMessage()
	// 	h.ScheduledMessageHandler.SendMessage(m, id)
	// })

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

type Handlers struct {
	AvyReportHandler        handlers.AvyReportHandler
	ScheduledMessageHandler handlers.ScheduledMessageHandler
}

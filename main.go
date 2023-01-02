package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/syslog"
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
	"github.com/honeybadger-io/honeybadger-go"
)

var (
	exp bool
)

func main() {
	defer honeybadger.Monitor()
	flag.BoolVar(&exp, "e", false, "Point channels to experiments channel")
	flag.Parse()

	cfg, err := config.GetConfig()
	if err != nil {
		honeybadger.Notify("main: Error initializing config: %w", err)
		panic(err)
	}

	w, err := syslog.Dial("udp", cfg.PaperTrailHost, syslog.LOG_EMERG | syslog.LOG_KERN , "skele-bot")
	if err != nil {
		log.Fatal("failed to dial syslog")
	}

	w.Info("Info log")
	w.Err("Err log")
	w.Notice("Notice log")

	c := colly.NewCollector(colly.AllowedDomains(crawlers.AvyCenterDomains...))
	ac := crawlers.NewAvyCrawler(c)

	session, err := discordgo.New(cfg.BotToken)
	if err != nil {
		honeybadger.Notify("main: Error creating bot session: %w", err)
		panic(err)
	}

	bot := bot.NewDiscordBot(session, exp)

	h := Handlers{
		AvyCrawlerHandler: handlers.AvyCrawlerHandler{
			AvyCrawler: ac,
			DiscordBot: bot,
		},
		ScheduledMessageHandler: handlers.ScheduledMessageHandler{
			DiscordBot: bot,
		},
	}

	ctx, cancel := context.WithCancel((context.Background()))

	go startCron(ctx, h, cancel)
	go startBot(ctx, bot, h, cancel)

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

	mtnTZ, _ := time.LoadLocation("America/Denver")

	s := gocron.NewScheduler(mtnTZ)

	s.Every(1).Days().At("07:30").Do(func() {
		h.AvyCrawlerHandler.SendAvyReport()
	})

	s.Every(1).Days().At("23:30").Do(func() {
		h.AvyCrawlerHandler.SendTodaysAvyList()
	})

	s.Every(1).Days().At("00:00").Do(func() {
		h.AvyCrawlerHandler.SendTodaysAvyList()
	})

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

func startBot(pCtx context.Context, bot *bot.DiscordBot, h Handlers, exit context.CancelFunc) {
	fmt.Println("Starting bot session...")
	bot.RegisterHandlers()

	err := bot.Session.Open()
	if err != nil {
		honeybadger.Notify("startBot: Error opening websocket connection: %w", err)
		return
	}

	// defer exit()
	// defer fmt.Println("Exiting bot session...")
}

type Handlers struct {
	AvyCrawlerHandler        handlers.AvyCrawlerHandler
	ScheduledMessageHandler handlers.ScheduledMessageHandler
}

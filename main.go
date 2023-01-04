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
	"skele/internal/loggers"
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

	w := loggers.NewWriter(cfg.PaperTrailHost)

	log.SetOutput(w.Writer)
	log.Println("hellooo from main")


	w.LogInfo("Info log")
	w.LogError("Err log")
	w.LogNotice("Notice log")

	c := colly.NewCollector(colly.AllowedDomains(crawlers.AvyCenterDomains...))
	ac := crawlers.NewAvyCrawler(c)

	session, err := discordgo.New(cfg.BotToken)
	if err != nil {
		err = fmt.Errorf("main: Error creating bot session: %w", err)
		log.Println("err.Error() ", err.Error())
		honeybadger.Notify(err)
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

	go startCron(ctx, h, w.Writer, cancel)
	go startBot(ctx, bot, h, w.Writer, cancel)

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
	log.Println("Gracefully returning to the grave...")
}

func startCron(pCtx context.Context, h Handlers, w *syslog.Writer, exit context.CancelFunc) {
	w.Info("Starting cron...")
	defer exit()
	defer log.Println("Exiting cron...")

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

	s.StartBlocking()

	// TODO: make contexts work
	// <- pCtx.Done()
}

func startBot(pCtx context.Context, bot *bot.DiscordBot, h Handlers, w *syslog.Writer, exit context.CancelFunc) {
	log.Println("Starting bot session...")
	bot.RegisterHandlers()

	err := bot.Session.Open()
	if err != nil {
		err= fmt.Errorf("startBot: Error opening websocket connection: %w", err)
		log.Println("err ",err)
		honeybadger.Notify(err)
		return
	}

	// defer exit()
	// defer fmt.Println("Exiting bot session...")
}

type Handlers struct {
	AvyCrawlerHandler        handlers.AvyCrawlerHandler
	ScheduledMessageHandler handlers.ScheduledMessageHandler
}

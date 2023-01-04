package crawlers

import (
	"fmt"
	"log"
	"log/syslog"
	"skele/internal/data"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/honeybadger-io/honeybadger-go"
)

var (
	AvyCenterDomains = []string{"https://utahavalanchecenter.org/", "https://utahavalanchecenter.org", "utahavalanchecenter.org/", "www.utahavalanchecenter.org/", "utahavalanchecenter.org"}
)

type AvyCrawler struct {
	Collector *colly.Collector
	Writer *syslog.Writer
}

func NewAvyCrawler(collector *colly.Collector) *AvyCrawler {
	return &AvyCrawler{
		Collector: collector,
	}
}

func (ac *AvyCrawler) GetReport() (rp data.AvyReport, err error) {
	configureCrawler(ac)

	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		rp.Date = e.ChildText(".text_01 .nowrap")
		rp.Details = e.ChildText(".text_03")
		rp.ImageUrl = e.ChildAttr(".compass-width", "src")
		rp.SpecialBulletin = e.ChildText(".page-content .mb3")

		log.Printf("\n\nDate: %v\nDetails: %v\nImgUrl: %v\nSpecialBulletin: %v\n\n", rp.Date, rp.Details, rp.ImageUrl, rp.SpecialBulletin)
	})
	
	fmt.Printf("GetReport - visiting at %v\n", time.Now())
	url := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, data.AvyUrlPaths.Forecast)
	err = ac.Collector.Visit(url)
	if err != nil {
		err = fmt.Errorf("GetAvyReport: Visit Error: %w ", err)
		log.Println(err)

		return
	}
	log.Printf("GetReport - visited at %v\n", time.Now())

	return
}

func (ac *AvyCrawler) GetTodaysAvyList() (avs []data.Avy, today string, err error) {
	configureCrawler(ac)

	mtnTZ, _ := time.LoadLocation("America/Denver")

	now := time.Now().In(mtnTZ)
	today = fmt.Sprintf("%v/%v/%v", now.Month(), now.Day(), now.Year())

	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		var avy data.Avy
		
		e.ForEach("tbody tr", func(_ int, e *colly.HTMLElement) {
			date := e.ChildText(".date-display-single")

			if strings.EqualFold(date, today) {
				log.Println("Avy today - adding to list")
				avy.Date = date
				avy.Title = e.ChildText(".views-field-title")
				avy.Url = e.ChildAttr(".views-field-title a", "href")
				avy.Region = e.ChildText(".views-field-field-region-forecaster")

				avs = append(avs, avy)
				} 

				log.Printf("%v avys %v\n", date ,avs)
			})
			
		log.Println("No avys today")
	})

	url := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, data.AvyUrlPaths.Avalanches)
	err = ac.Collector.Visit(url)
	if err != nil {
		err = fmt.Errorf("GetTodaysAvyList: Visit Error: %w", err)
		log.Println(err)

		return
	}

	return
}

func configureCrawler(ac *AvyCrawler) {
	ac.Collector.Limit(&colly.LimitRule{
		Parallelism: 1, RandomDelay: 7 * time.Second,
	})

	ac.Collector.SetRequestTimeout(60 * time.Second)

	ac.Collector.CheckHead = true

	ac.Collector.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %v\n", r.URL.String())
	})

	ac.Collector.OnResponse(func(r *colly.Response) {
		log.Printf("\nReceived response from: %v\n", r.Request.URL)
	})

	ac.Collector.OnError(func(r *colly.Response, err error) {
		err = fmt.Errorf("AvyCrawler: Error on response from: %v\n%w",r.Request.URL, err)
		log.Println(err)
		honeybadger.Notify(err)
	})
}

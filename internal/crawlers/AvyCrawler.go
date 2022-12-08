package crawlers

import (
	"fmt"
	"skele/internal/data"
	"time"

	"github.com/gocolly/colly"
)

var (
	AvyCenterDomains = []string{"https://utahavalanchecenter.org/", "https://utahavalanchecenter.org", "utahavalanchecenter.org/", "www.utahavalanchecenter.org/", "utahavalanchecenter.org"}
)

type AvyCrawler struct {
	Collector *colly.Collector
}

func NewAvyCrawler(collector *colly.Collector) *AvyCrawler {
	return &AvyCrawler{
		Collector: collector,
	}
}

func (ac *AvyCrawler) GetReport() (rp data.AvyReport, err error) {
	configureCrawler(ac)

	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		fmt.Println("attempting to crawl...")
		rp.Date = e.ChildText(".text_01 .nowrap")
		rp.Details = e.ChildText(".text_03")
		rp.ImageUrl = e.ChildAttr(".compass-width", "src")

		fmt.Printf("\n\nDate: %v\nDetails: %v\nImgUrl: %v\n\n", rp.Date, rp.Details, rp.ImageUrl)
	})

	url := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, data.AvyUrlPaths.Forecast)
	ac.Collector.Visit(url)

	return
}

func (ac *AvyCrawler) GetTodaysAvyList() (av data.AvyList, err error) {
	configureCrawler(ac)

	now := time.Now()
	// today := fmt.Sprintf("%d/%d/%d", now.Month(), now.Day(), now.Year())
	today := fmt.Sprintf("%d/%d/%d", now.Month(), 6, now.Year())
	fmt.Println(today)

	// var dates []string

	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		fmt.Println("attempting to crawl...")

		e.ForEach(".date-display-single", func(_ int, el *colly.HTMLElement) {
			if el.Text == today {
				fmt.Println("TODAY")
			}
		})
	})

	url := fmt.Sprintf("%s%s", data.AvyUrlPaths.BaseUrl, data.AvyUrlPaths.Avalanches)
	ac.Collector.Visit(url)

	return
}

func configureCrawler(ac *AvyCrawler) {
	ac.Collector.Limit(&colly.LimitRule{
		Parallelism: 1, RandomDelay: 7 * time.Second,
	})

	ac.Collector.SetRequestTimeout(60 * time.Second)

	ac.Collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %v...", r.URL.String())
	})

	ac.Collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("\nReceived response from: %v\n", r.Request.URL)
	})

	ac.Collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Response error: ", err)
	})
}

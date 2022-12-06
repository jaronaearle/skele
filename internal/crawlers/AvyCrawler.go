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
	URL       string
}

func NewAvyCrawler(collector *colly.Collector) *AvyCrawler {
	url := "https://utahavalanchecenter.org/forecast/salt-lake"
	return &AvyCrawler{
		Collector: collector,
		URL:       url,
	}
}

func (ac *AvyCrawler) GetReport() (rp data.AvyReport, err error) {
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

	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		fmt.Println("attempting to crawl...")
		rp.Date = e.ChildText(".text_01 .nowrap")
		rp.Details = e.ChildText(".text_03")
		rp.ImageUrl = e.ChildAttr(".compass-width", "src")

		fmt.Printf("\n\nDate: %v\nDetails: %v\nImgUrl: %v\n\n", rp.Date, rp.Details, rp.ImageUrl)
	})
	ac.Collector.Visit(ac.URL)
	return
}

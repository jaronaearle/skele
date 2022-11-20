package crawlers

import (
	"discord-hooks/internal/data"
	"fmt"
	"time"

	"github.com/gocolly/colly"
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
	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		rp.Date = e.ChildText(".text_01 .nowrap")
		rp.Details = e.ChildText(".text_3")
		rp.ImageUrl = e.ChildAttr(".compass-width", "src")

		fmt.Printf("\n\nDate: %v\nDetails: %v\nImgUrl: %v\n\n", rp.Date, rp.Details, rp.ImageUrl)
		ac.Collector.Visit(ac.URL)
	})

	ac.Collector.Limit(&colly.LimitRule{
		Parallelism: 1, RandomDelay: 7 * time.Second,
	})

	ac.Collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %v...", r.URL.String())
	})

	return
}

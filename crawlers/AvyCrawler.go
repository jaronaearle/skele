package crawlers

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type AvyCrawler struct {
	Collector colly.Collector
	URL string
}

type AvyReport struct {
	Date string
	Details string
	ImageUrl string
}

func NewAvyCrawler(collector colly.Collector) *AvyCrawler {
	url := "https://utahavalanchecenter.org/forecast/salt-lake"
	return &AvyCrawler{
		Collector: collector,
		URL: url,
	}
}

func (ac *AvyCrawler) GetReport() (rp AvyReport, err error) {
	ac.Collector.OnHTML(".view-content", func(e *colly.HTMLElement) {
		rp.Date = ".text_01 .nowrap"
		rp.Details =".text_3"
		rp.ImageUrl = e.ChildAttr(".compass-width", "src")
		
		fmt.Printf("Date: %v\nDetails: %v\nImgUrl: %v", rp.Date, rp.Details, rp.ImageUrl)
	})

	ac.Collector.Limit(&colly.LimitRule{
		Parallelism: 1, RandomDelay: 7 * time.Second,
	})

	ac.Collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %v...", r.URL.String())
	})

	ac.Collector.Visit(ac.URL)
	ac.Collector.Wait()

	return 
}

// func (ac *AvyCrawler) GetReportImageURL() (s string, err error) {
// 	return "", nil
// }

// func (ac *AvyCrawler) GetReportDetails() (d string, err error) {
// 	return "", nil
// }
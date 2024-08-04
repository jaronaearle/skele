package crawlers

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

var PermitDomains = []string{"https://www.recreation.gov/", "www.recreation.gov/", "recreation.gov/", "https://www.recreation.gov", "www.recreation.gov", "recreation.gov"}

// var ExpURL = "https://www.recreation.gov/permits/4251902"
var ExpURL = "https://www.recreation.gov/permits/4251902/registration/detailed-availability"

type PermitCrawler struct {
	Collector *colly.Collector
}

func NewPermitCrawler(collector *colly.Collector) *PermitCrawler {
	return &PermitCrawler{
		Collector: collector,
	}
}

func (pc *PermitCrawler) GetAvailablePermitDates() (err error) {
	// class="permit-availability-calendar"
	permitCollector := pc.Collector.Clone()
	configurePermitCrawler(permitCollector)

	// should be all needed selectors - need a way to get urls to chec

	permitCollector.OnHTML(".permit-availability-calendar", func(e *colly.HTMLElement) {
		test := e.ChildText("h1.h3")
		// calendar month header May 2024
		// month := e.ChildText("h2.rec-sr-only")

		// next month button
		// nextMonth := e.ChildText("button.next-prev-button") // probs need nth(1)

		// caldendar day grid cell
		// day := e.ChildText("div.calendar-cell-td > div.available")

		// fmt.Printf("\n\nDate: %v\nMonth: %v\nDay: %v\n\n", month, nextMonth, day)
		fmt.Printf("Title::::: >>>>>>>>> %v", test)
	})

	err = permitCollector.Visit(ExpURL)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	return
}

func configurePermitCrawler(collector *colly.Collector) {
	collector.Limit(&colly.LimitRule{
		Parallelism: 1, RandomDelay: 7 * time.Second,
	})

	collector.SetRequestTimeout(60 * time.Second)

	collector.CheckHead = true

	collector.AllowURLRevisit = true

	collector.OnRequest(func(r *colly.Request) {
		// log.Printf("Visiting %v\n", r.URL.String())
		fmt.Printf("Visiting %v\n", r.URL.String())
		fmt.Printf("Req %v\n", r.Method)
		fmt.Printf("Headers %v\n", r.Headers)
	})

	collector.OnResponse(func(r *colly.Response) {
		// log.Printf("\nReceived response from: %v\n", r.Request.URL)
		fmt.Printf("\nReceived response from: %v\n", r.Request.URL)
	})

	collector.OnError(func(r *colly.Response, err error) {
		err = fmt.Errorf("AvyCrawler: Error on response from: %v\n%w", r.Request.URL, err)
		// log.Println(err)
		fmt.Println(err)
		// honeybadger.Notify(err)
	})
}

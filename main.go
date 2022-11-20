package main

import (
	"discord-hooks/internal/crawlers"
	"discord-hooks/internal/data"
	"fmt"

	"github.com/gocolly/colly"
)

type Payload struct {
	Content string `json:"content"`
}

func main() {
	c := colly.NewCollector()

	ac := crawlers.NewAvyCrawler(c)
	
	test, _ := ac.GetReport()

	fmt.Println(test)
	
	testData := data.AvyReport{
		Date: "Saturday, November 19, 2022",
		Details: "lksdjfaksljdfaksljdflaksdjfsdkljf alskdfjaklsjdfaskldjf asdffsd",
		ImageUrl: "Saturday, November 19, 2022",
	}
	
	fmt.Println(testData)
	

	// data := Payload{
	// 	Content: "From cron, with <3",
	// }
	// s := gocron.NewScheduler(time.UTC)

	// s.Every(7).Seconds().Do(func() {
	// 	postTest(data)
	// })

	// s.StartBlocking()


	// err := postTest(data)
	// if err != nil {
	// 	panic(err)
	// }
}

// func postTest(p Payload) (err error) {
// 	pBytes, err := json.Marshal(p)
// 	if err != nil {
// 		fmt.Printf("postTest: Error marshalling payload: %v", err)
// 		return
// 	}

// 	body := bytes.NewReader(pBytes)

// 	req, err := http.NewRequest("POST", "https://discord.com/api/webhooks/1043741007398846604/wwGGUlBiV_OLKS-m3Xn0TVpM-dayDb5X48BXRuyBgKvKdAmrtZQCV9Jd_kO8hdqHJwT6", body)
// 	if err != nil {
// 		fmt.Printf("postTest: Error building new request: %v", err)
// 		return err
// 	}

// 	req.Header.Set("Accept", "application/json")
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Printf("postTest: Error sending request: %v", err)
// 		return err
// 	}
// 	fmt.Println("Request sent")
// 	defer resp.Body.Close()

// 	return nil
// }

// curl -i -H "Accept: application/json" -H "Content-Type:application/json" -X POST --data "{\"content\": \"Posted Via Command line\"}" https://discord.com/api/webhooks/1043741007398846604/wwGGUlBiV_OLKS-m3Xn0TVpM-dayDb5X48BXRuyBgKvKdAmrtZQCV9Jd_kO8hdqHJwT6

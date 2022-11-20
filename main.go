package main

import (
	"discord-hooks/internal/bot"
	"discord-hooks/internal/crawlers"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

const botToken = "<bot token>"
type Payload struct {
	Content string `json:"content"`
}

func main() {
	c := colly.NewCollector(colly.AllowedDomains("https://utahavalanchecenter.org/","https://utahavalanchecenter.org", "utahavalanchecenter.org/", "www.utahavalanchecenter.org/", "utahavalanchecenter.org"))

	ac := crawlers.NewAvyCrawler(c)

	test, err := ac.GetReport()
	if err != nil {
		fmt.Println("Error calling GetReport: ", err)
		panic(err)
	}

	fmt.Println(test)

	ds, err := discordgo.New(fmt.Sprintf("Bot %s", botToken))
	if err != nil {
		panic(err)
	}

	db := bot.NewDiscordBot(ds)
	db.StartBotSession()

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

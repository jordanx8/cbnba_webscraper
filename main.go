package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

// type Player struct {
// 	Rank     int    `json:"rank"`
// 	Name     string `json:"name"`
// 	School   string `json:"school"`
// 	Position string `json:"position"`
// 	Height   string `json:"height"`
// 	Weight   string `json:"weight"`
// }

type Player struct {
	Rank int    `json:"rank"`
	Name string `json:"name"`
	// School   string `json:"school"`
	// Position string `json:"position"`
	// Height   string `json:"height"`
	// Weight   string `json:"weight"`
}

func main() {
	allPlayers := make([]Player, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("espn.com", "www.espn.com"))

	collector.OnHTML(".draftTable__row li", func(element *colly.HTMLElement) {
		playerRank, err := strconv.Atoi(element.ChildText("a"))
		if err != nil {
			fmt.Println("Could not get data-rank")
		}

		playerText := element.Text

		player := Player{
			Rank: playerRank,
			Name: playerText,
		}

		allPlayers = append(allPlayers, player)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://www.espn.com/nba/draft/bestavailable")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(allPlayers)
}

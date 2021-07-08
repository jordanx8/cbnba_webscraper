package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

type Player struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`
	School   string `json:"school"`
	Position string `json:"position"`
}

func main() {
	allPlayers := make([]Player, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("espn.com", "www.espn.com"))

	collector.OnHTML(".draftTable__row li", func(element *colly.HTMLElement) {
		playerRank, err := strconv.Atoi(element.ChildText("span.draftTable__headline--pick"))
		if err != nil {
			fmt.Println("Could not get data-rank")
		}

		player := Player{
			Rank:     playerRank,
			Name:     element.ChildText("span.draftTable__headline--player"),
			School:   element.ChildText("span.draftTable__headline--school"),
			Position: element.ChildText("span.draftTable__headline--pos"),
		}

		allPlayers = append(allPlayers, player)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnHTML(".pagination a", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		collector.Visit(nextPage)
	})

	collector.Visit("https://www.espn.com/nba/draft/bestavailable/_/position/ovr/page/1")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(allPlayers)

	writeJSON(allPlayers)
}

func writeJSON(data []Player) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Unable to create JSON file.")
		return
	}

	_ = ioutil.WriteFile("../playerdata.json", file, 0644)
}

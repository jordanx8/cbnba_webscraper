package webscraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/thedevsaddam/gojsonq/v2"
)

type Player struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`
	School   string `json:"school"`
	Position string `json:"position"`
	NextGame string `json:"nextGame"`
}

func ScrapeESPNTop100() {
	// "children".Pluck()
	// teamsJSON := "https://site.web.api.espn.com/apis/v2/sports/basketball/mens-college-basketball/standings?region=us&lang=en&contentorigin=espn&group=50&sort=playoffseed%3Aasc%2Cvsconf_winpercent%3Adesc%2Cvsconf_wins%3Adesc%2Cvsconf_losses%3Aasc%2Cvsconf_gamesbehind%3Aasc&includestats=playoffseed%2Cvsconf%2Cvsconf_gamesbehind%2Cvsconf_winpercent%2Ctotal%2Cwinpercent%2Chome%2Croad%2Cstreak%2Cvsaprankedteams%2Cvsusarankedteams&season=2021"
	// client := resty.New()
	// resp, err := client.R().Get(teamsJSON)
	// if resp.StatusCode() == 200 {
	// 	if parsedResponse := resp.Body(); err == nil {
	// 		var teamInfo map[string]interface{}
	// 		err := json.Unmarshal(parsedResponse, &teamInfo)
	// 		if err != nil {
	// 			return
	// 		}
	// 		fmt.Println(teamInfo[""])
	// 	}
	// }

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
			NextGame: GetNextEvent(element.ChildText("span.draftTable__headline--school")),
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

	WriteJSON(allPlayers)
}

func WriteJSON(data []Player) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Unable to create JSON file.")
		return
	}

	_ = ioutil.WriteFile("../playerdata.json", file, 0644)
}

func GetNextEvent(school string) string {
	api := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/teams/%s", school)
	client := resty.New()
	resp, err := client.R().Get(api)
	if resp.StatusCode() == 200 {
		if parsedResponse := resp.String(); err == nil {
			// var playerInfo map[string]interface{}
			// err := json.Unmarshal(parsedResponse, &playerInfo)
			// if err != nil {
			// 	return school
			// }
			nextevent := gojsonq.New().FromString(parsedResponse).From("team.nextEvent").Pluck("date")
			// parsedResponse, err = json.MarshalIndent(playerInfo, "", " ")
			// if err != nil {
			// 	return school
			// }
			str := fmt.Sprintf("%v", nextevent)
			re := regexp.MustCompile(`\[([^\[\]]*)\]`)
			submatchall := re.FindAllString(str, -1)
			for _, element := range submatchall {
				element = strings.Trim(element, "[")
				element = strings.Trim(element, "]")
				return element
			}
			// date, err := time.Parse("2006-01-02 15:04", str)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println("My Date Reformatted:\t", date.Format(time.RFC822))
		}
	}
	return school
}

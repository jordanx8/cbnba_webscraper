package webscraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

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

type CollegeTeam struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func ScrapeESPNTop100() {
	fmt.Println("ScrapeESPNTop100()")

	allPlayers := make([]Player, 0)
	allCollegeTeams := make([]CollegeTeam, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("espn.com", "www.espn.com"))

	collector.OnHTML(".TeamLinks.flex.items-center div", func(element *colly.HTMLElement) {
		teamName := element.ChildText("h2.di.clr-gray-01.h5")
		if teamName == "" {
			return
		}
		teamURL := element.ChildAttr("a", "href")
		teamURLEnd := strings.Split(teamURL, "/mens-college-basketball/team/_/id/")
		re := regexp.MustCompile("[0-9]+")
		teamID := re.FindAllString(teamURLEnd[1], -1)
		collegeTeam := CollegeTeam{
			Name: teamName,
			ID:   teamID[0],
		}
		allCollegeTeams = append(allCollegeTeams, collegeTeam)
	})

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
			NextGame: GetNextEvent(element.ChildText("span.draftTable__headline--school"), allCollegeTeams),
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

	collector.Visit("https://www.espn.com/mens-college-basketball/teams")
	collector.Visit("https://www.espn.com/nba/draft/bestavailable/_/position/ovr/page/1")

	WritePlayerJSON(allPlayers)
}

func WritePlayerJSON(data []Player) {
	fmt.Println("Attempting to create playerdata.json.")
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Unable to create JSON file.")
		return
	}

	_ = ioutil.WriteFile("./playerdata.json", file, 0644)
	fmt.Println("Success.")
}

func WriteTeamJSON(data []CollegeTeam) {
	fmt.Println("Attempting to create collegeteams.json.")
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Unable to create JSON file.")
		return
	}

	_ = ioutil.WriteFile("./collegeteams.json", file, 0644)
	fmt.Println("Success.")
}

func GetNextEvent(school string, allCollegeTeams []CollegeTeam) string {
	api := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/mens-college-basketball/teams/%s", school)
	client := resty.New()
	resp, err := client.R().Get(api)
	if resp.StatusCode() == 200 {
		if parsedResponse := resp.String(); err == nil {
			nextevent := gojsonq.New().FromString(parsedResponse).From("team.nextEvent").Pluck("date")
			str := fmt.Sprintf("%v", nextevent)
			re := regexp.MustCompile(`\[([^\[\]]*)\]`)
			submatchall := re.FindAllString(str, -1)
			for _, element := range submatchall {
				element = strings.Trim(element, "[")
				element = strings.Trim(element, "]")
				layout := "2006-01-02T15:04Z"
				t, err := time.Parse(layout, element)
				if err != nil {
					fmt.Println(err)
				}
				location, err := time.LoadLocation("America/New_York")
				if err != nil {
					fmt.Println(err)
				}
				tlocal := t.In(location)
				return tlocal.Format(time.UnixDate)
			}
		}
	}
	for a := range allCollegeTeams {
		if strings.HasPrefix(allCollegeTeams[a].Name, school) {
			return GetNextEvent(allCollegeTeams[a].ID, allCollegeTeams)
		}
	}
	return school
}

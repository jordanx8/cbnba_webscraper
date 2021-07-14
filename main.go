package main

import (
	webscraper "github.com/jordanx8/webscraper/espn_webscraper"
	seeder "github.com/jordanx8/webscraper/mongodb_seeder"
)

func main() {
	fmt.Println("App starting")
	webscraper.ScrapeESPNTop100()
	seeder.SeedPlayerData()
}

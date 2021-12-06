package webscraper

import (
	"context"
	"fmt"

	seeder "github.com/jordanx8/webscraper/mongodb_seeder"
	pb "github.com/jordanx8/webscraper/proto"
)

func (s WebscraperService) ScrapeAndSeed(ctx context.Context, empty *pb.Empty) (*pb.ScrapeAndSeedResponse, error) {
	fmt.Println("Running ScrapeAndSeed()")
	err := ScrapeESPNTop100()
	if err != nil {
		return &pb.ScrapeAndSeedResponse{Success: 0}, err
	}
	err = seeder.SeedPlayerData()
	if err != nil {
		return &pb.ScrapeAndSeedResponse{Success: 0}, err
	}
	return &pb.ScrapeAndSeedResponse{Success: 1}, err
}

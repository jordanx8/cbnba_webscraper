package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	webscraper "github.com/jordanx8/webscraper/espn_webscraper"
	pb "github.com/jordanx8/webscraper/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("App starting")
	//webscraper.ScrapeESPNTop100()
	//seeder.SeedPlayerData()

	ctx := context.Background()
	log.Println("starting up")
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	s := grpc.NewServer()
	server := webscraper.WebscraperService{MongoClient: client}
	pb.RegisterWebscraperServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

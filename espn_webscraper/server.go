package webscraper

import (
	pb "github.com/jordanx8/webscraper/proto"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebscraperService struct {
	MongoClient *mongo.Client
	pb.UnimplementedWebscraperServiceServer
}

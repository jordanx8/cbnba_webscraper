package webscraper

import (
	"context"
	"fmt"
	"log"

	pb "github.com/jordanx8/webscraper/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p Player) ConvertToGRPC() *pb.Player {
	return &pb.Player{
		Rank:     int32(p.Rank),
		Name:     p.Name,
		School:   p.School,
		Position: p.Position,
		Nextgame: p.NextGame,
	}

}

func (s WebscraperService) GetPlayers(ctx context.Context, playerRequest *pb.PlayerRequest) (*pb.PlayerArray, error) {
	fmt.Println("Running GetPlayers()")
	client, err := GetClient()

	var players []*pb.Player
	col := client.Database("cbnba").Collection("PlayerData")

	var filters bson.D
	if playerRequest.GetName() != "" {
		filters = append(filters, primitive.E{Key: "name", Value: bson.M{"$regex": playerRequest.GetName(), "$options": "i"}})
	}
	if playerRequest.GetPosition() != "" {
		filters = append(filters, primitive.E{Key: "position", Value: bson.M{"$regex": playerRequest.GetPosition()}})
	}
	if playerRequest.GetSchool() != "" {
		filters = append(filters, primitive.E{Key: "school", Value: bson.M{"$regex": playerRequest.GetSchool(), "$options": "i"}})
	}
	if playerRequest.GetRank() != 0 {
		filters = append(filters, primitive.E{Key: "rank", Value: bson.M{"$eq": playerRequest.GetRank()}})
	}
	filters = append(filters, bson.E{})

	cur, _ := col.Find(ctx, filters)
	if err = cur.All(ctx, &players); err != nil {
		log.Fatal(err)
	}
	var x = 0
	for x < len(players) {
		players[x].Nextgame = string(players[x].Nextgame)
		x++
	}

	return &pb.PlayerArray{Players: players}, err
}

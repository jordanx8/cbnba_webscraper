package mongodb_seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Player struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`
	School   string `json:"school"`
	Position string `json:"position"`
	NextGame string `json:"nextGame"`
}

func SeedPlayerData() {

	fmt.Println("SeedPlayerData()")

	clientOptions := options.Client().ApplyURI("mongodb://localhost")
	//connect to MongoDb, if error then display the issue
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//check to ensure the connection went through
	err = client.Ping(context.TODO(), nil)
	ctx, _ := context.WithTimeout(context.Background(), 240*time.Second)

	col := client.Database("cbnba").Collection("PlayerData")
	if err = col.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	col = client.Database("cbnba").Collection("PlayerData")
	//check for error again
	if err != nil {
		log.Fatal(err)
	} //else mongo has been connected
	fmt.Println("Successfully connected to Mongo")

	byteValues, err := ioutil.ReadFile("./playerdata.json")
	if err != nil {
		fmt.Println(err)
	}

	var docs []Player

	err = json.Unmarshal(byteValues, &docs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.TypeOf(docs))

	for i := range docs {
		doc := docs[i]
		result, insertErr := col.InsertOne(ctx, doc)
		if insertErr != nil {
			fmt.Println(insertErr)
		} else {
			fmt.Println(result)
		}
	}

	fmt.Println("playerdata.json seeding finished.")

}

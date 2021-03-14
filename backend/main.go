package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Paper struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text       string             `json:"text,omitempty" bson:"firstname,text"`
	FirstTeam  string             `json:"firstteam,omitempty" bson:"firstteam,omitempty"`
	SecondTeam string             `json:"secondteam,omitempty" bson:"secondteam,omitempty"`
	Score      string             `json:"score,omitempty" bson:"score,omitempty"`
}

var databaseUrl string = ""

func main() {
	fmt.Println("Hello world")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseUrl))
	if err != nil {
		log.Fatal(err)
	}

	if client != nil {
		fmt.Println("Hello Database")
	}
	router := mux.NewRouter()

	http.ListenAndServe(":8000", router)

}

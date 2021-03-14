package main

import (
	"context"
	"encoding/json"
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

func CreatePaper(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var paper Paper
	_ = json.NewDecoder(request.Body).Decode(&paper)
	collection := client.Database("Blog").Collection("papers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, paper)
	json.NewEncoder(response).Encode(result)
}

func main() {
	fmt.Println("Hello world")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseUrl))
	if err != nil {
		log.Fatal(err)
	}

	if client != nil {
		fmt.Println("Hello Database")
	}
	router := mux.NewRouter()
	router.HandleFunc("/papers/add", CreatePaper).Methods("POST")
	http.ListenAndServe(":8000", router)

}

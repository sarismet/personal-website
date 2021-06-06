package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gomail "gopkg.in/gomail.v2"
)

type MessageInfo struct {
	FirstName string `json:"firstName,omitempty" bson:"firstname,text"`
	LastName  string `json:"lastName,omitempty" bson:"lastName,text"`
	Email     string `json:"email,omitempty" bson:"email,text"`
	Message   string `json:"message,omitempty" bson:"message,text"`
}

//REFERENCE: https://gist.github.com/ivanmrchk/e30eb45808536159bbec9aac20058b78
func (mi *MessageInfo) sendMail() {

	t := template.New("email-template.html")

	var err error
	t, err = t.ParseFiles("email-template.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mi); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", MyEmail)
	m.SetHeader("To", MailTo)
	m.SetHeader("Subject", "Contact")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, MyEmail, Password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func notification() {
	m := gomail.NewMessage()
	m.SetHeader("From", MyEmail)
	m.SetHeader("To", MailTo)
	m.SetHeader("Subject", "Notification")
	m.SetBody("text", "Someone has seen your resume")
	d := gomail.NewDialer("smtp.gmail.com", 587, MyEmail, Password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

var myclient *mongo.Client

type Paper struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text       string             `json:"text,omitempty" bson:"text,omitempty"`
	FirstTeam  string             `json:"firstteam,omitempty" bson:"firstteam,omitempty"`
	SecondTeam string             `json:"secondteam,omitempty" bson:"secondteam,omitempty"`
	Score      string             `json:"score,omitempty" bson:"score,omitempty"`
}

var databaseUrl string 
func CreatePaper(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var paper Paper
	_ = json.NewDecoder(request.Body).Decode(&paper)
	collection := myclient.Database("Blog").Collection("papers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, paper)
	json.NewEncoder(response).Encode(result)
}

func Notify(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	notification()
	json.NewEncoder(response).Encode("Notification is sent")
}

func sendEmail(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var messageInfo MessageInfo
	_ = json.NewDecoder(request.Body).Decode(&messageInfo)
	messageInfo.sendMail()
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

	myclient = client

	router := mux.NewRouter()
	router.HandleFunc("/papers/add", CreatePaper).Methods("POST")
	router.HandleFunc("/papers/sendEmail", sendEmail).Methods("POST")
	router.HandleFunc("/papers/notify", Notify).Methods("POST")

	http.ListenAndServe(":8000", router)

}

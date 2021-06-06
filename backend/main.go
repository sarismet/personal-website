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
	"net/smtp"

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
func (mi *MessageInfo) sendMail(_body string, _subject string) {

	from := MyEmail //ex: "John.Doe@gmail.com"
	password := Password   // ex: "ieiemcjdkejspqz"
	// receiver address
	toEmail := MailTo // ex: "Jane.Smith@yahoo.com"
	to := []string{toEmail}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := _subject
	body := _body
	message := []byte(subject + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// send mail
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err:", err)
		return
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
	var messageInfo MessageInfo
	messageInfo.sendMail("somebody has viewed your resume","notification")
	json.NewEncoder(response).Encode("Notification is sent")
}

func sendEmail(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var mi MessageInfo
	_ = json.NewDecoder(request.Body).Decode(&messageInfo)
	_body := "Hi my name is "+mi.FirstName+" "+mi.LastName+"and my message is \n"+mi.Message+"\n My connection email address is : "+mi.Email
	messageInfo.sendMail(_body,"Someone has sent you email")
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

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/smtp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Subscriber struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Age      uint32             `bson:"age"`
}

type HtmlTemplate struct {
	Name  string `bson:"name"`
	Title string `bson:"title"`
	Body  string `bson:"body"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:NewPassword123@cluster0.njmzc.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connected")

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	subscriberCollection := client.Database("MailGaner").Collection("Subscriber")
	htmlTemplateCollection := client.Database("MailGaner").Collection("HtmlTemplate")

	subscriberFilter := bson.D{}
	htmlTemplateFilter := bson.D{}

	subscriberList := getSubscriberList(subscriberCollection, subscriberFilter)
	htmlTemplateList := getHtmlTemplateList(htmlTemplateCollection, htmlTemplateFilter)

	for _, item := range subscriberList {
		sendMessage("horry.morry@mail.ru", item.Email, htmlTemplateList[0].Title, htmlTemplateList[0].Body, item.Username)
	}
}

func getSubscriberList(collection *mongo.Collection, filter primitive.D) []Subscriber {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var subscribers []Subscriber
	for cursor.Next(context.TODO()) {
		var subscriber Subscriber
		cursor.Decode(&subscriber)
		subscribers = append(subscribers, subscriber)
	}

	return subscribers
}

func getHtmlTemplateList(collection *mongo.Collection, filter primitive.D) []HtmlTemplate {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var htmlTemplates []HtmlTemplate
	for cursor.Next(context.TODO()) {
		var htmlTemplate HtmlTemplate
		cursor.Decode(&htmlTemplate)
		htmlTemplates = append(htmlTemplates, htmlTemplate)
	}

	return htmlTemplates
}

func sendMessage(from string, to string, title string, body string, subscriberName string) {
	smtpServer := "smtp.mail.ru"
	auth := smtp.PlainAuth(
		"",
		"horry.morry@mail.ru",
		"3NbdPBjE2PLph4KbM80M",
		smtpServer,
	)

	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = mime.QEncoding.Encode("UTF-8", title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	bodyMessage := fmt.Sprintf(body, subscriberName, to)
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(bodyMessage))

	err := smtp.SendMail(
		smtpServer+":25",
		auth,
		from,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}
}

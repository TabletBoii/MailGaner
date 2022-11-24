package main

import (
	"context"
	"fmt"
	"log"

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
	Name string `bson:"name"`
	Body string `bson:"body"`
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

	filter := bson.D{}

	cursor, err := subscriberCollection.Find(context.TODO(), filter)
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
}

package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionString := "mongodb+srv://sohammaha15:OCu9cwjwaI5yEHDF@cluster0.0upma.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	clientOptions := options.Client().ApplyURI(connectionString)
	Client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	Collection = Client.Database("Road-to-devsoc-25").Collection("items")
}

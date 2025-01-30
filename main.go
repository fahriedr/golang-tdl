package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fahriedr/golang-tdl/api"
	"github.com/fahriedr/golang-tdl/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient    *mongo.Client
	taskCollection *mongo.Collection
)

func init() {
	// context variable
	ctx := context.TODO()

	// database credentials
	dbUri := config.Envs.MongoUrl
	connectionOpts := options.Client().ApplyURI(dbUri)

	// database client
	mongoClient, err := mongo.Connect(ctx, connectionOpts)

	// error handling when connection fail
	if err != nil {
		fmt.Printf("an error ocurred when connect to mongoDB : %v", err)
		log.Fatal(err)
	}

	// Test connection to database
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB successfully connected")

	taskCollection = mongoClient.Database("go-tdl").Collection("tasks")
}

func main() {

	server := api.NewApiServer(taskCollection, config.Envs.Database, config.Envs.Port)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var collection = func() *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return client.Database("core").Collection("tink-linked")
}()

func getIsAccountLinked(username string) bool {
	user := UserItem{}
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != nil {
		sugar.Error(err)
	}
	return true
}

func setAccountLinked(username string) error {
	user := UserItem{}
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		sugar.Error(err)
	}
	user.AccountLinked = true
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		sugar.Error(err)
	}
	return nil
}

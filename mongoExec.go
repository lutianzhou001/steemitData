package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func insertInMongo(ctx context.Context, b []byte, cl *mongo.Client, id int, collection string) error {
	var bdoc interface{}
	err := bson.UnmarshalExtJSON(b, true, &bdoc)
	if err != nil {
		return err
	}
	coll := cl.Database("steemit").Collection(collection)
	_, err = coll.InsertOne(ctx, &bdoc)
	if err != nil {
		return err
	}
	log.Printf("Insert MongoDb %v "+collection+" Successfully", id)
	return nil
}

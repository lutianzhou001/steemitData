package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
)

func processingBlockMessage(ctx context.Context, block map[string]interface{}, cl *mongo.Client) error {
	if block["transactions"] != nil {
		txs := block["transactions"].([]interface{})
		for _, tx := range txs {
			operations := tx.(map[string]interface{})["operations"].([]interface{})
			if operations != nil {
				for _, operation := range operations {
					if operation.(map[string]interface{})["type"] == "vote_operation" {
						b, err := json.Marshal(operation)
						if err != nil {
							return err
						}
						err = insertInMongo(ctx, b, cl, 0, "vote")
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

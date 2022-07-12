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
						value := operation.(map[string]interface{})["value"].(map[string]interface{})
						if value != nil {
							b, err := json.Marshal(value)
							if err != nil {
								return err
							}
							err = insertInMongo(ctx, b, cl, 0, "vote")
							if err != nil {
								return err
							}
						}
					} else if operation.(map[string]interface{})["type"] == "comment_operation" {
						value := operation.(map[string]interface{})["value"].(map[string]interface{})
						// fmt.Println(value)
						if value != nil {
							if value["parent_author"] == "" {
								b, err := json.Marshal(value)
								if err != nil {
									return err
								}
								err = insertInMongo(ctx, b, cl, 0, "post")
								if err != nil {
									return err
								}
							} else {
								b, err := json.Marshal(value)
								if err != nil {
									return err
								}
								err = insertInMongo(ctx, b, cl, 0, "comment")
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
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
						value["timestamp"] = block["timestamp"]
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
								value["timestamp"] = block["timestamp"]
								value["body"] = ""
								b, err := json.Marshal(value)
								if err != nil {
									return err
								}
								err = insertInMongo(ctx, b, cl, 0, "post")
								if err != nil {
									return err
								}
							} else {
								continue
								//value["timestamp"] = block["timestamp"]
								//b, err := json.Marshal(value)
								//if err != nil {
								//	return err
								//}
								//err = insertInMongo(ctx, b, cl, 0, "comment")
								//if err != nil {
								//	return err
								//}
							}
						}
					} else {
						value := operation.(map[string]interface{})["value"].(map[string]interface{})
						if value != nil {
							collection := fmt.Sprint(operation.(map[string]interface{})["type"])
							value["timestamp"] = block["timestamp"]
							b, err := json.Marshal(value)
							if err != nil {
								return err
							}
							err = insertInMongo(ctx, b, cl, 0, collection)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}

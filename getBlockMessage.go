package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func convertBodyToBlockMessage(body []byte) ([]byte, map[string]interface{}, error) {
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)
	if result["result"] != nil {
		resMap := result["result"].(map[string]interface{})["block"].(map[string]interface{})
		res, err := json.Marshal(result["result"].(map[string]interface{})["block"])
		if err != nil {
			return nil, nil, err
		}
		return res, resMap, nil
	} else {
		return nil, nil, nil
	}
}

func getBlockMessage(ctx context.Context, id int, cl *mongo.Client) (map[string]interface{}, error) {
	url := "https://api.steemit.com"
	method := "POST"

	s := strconv.Itoa(id)
	payload := strings.NewReader(`{
    "jsonrpc":"2.0",
    "method":"block_api.get_block",
    "params":{"block_num":` + s + `},
    "id": ` + s + `
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	_, convertedMap, err := convertBodyToBlockMessage(body)
	if err != nil {
		return nil, err
	}

	return convertedMap, nil
}

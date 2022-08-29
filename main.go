package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoDB *mongo.Database
)

type Order struct {
	Product string
	Price   float64
	Client  string
}

func main() {
	ctx := context.Background()
	key := uuid.New().String()
	ctx = context.WithValue(ctx, "plunder", key)

	conf, err := NewConfig("config.yml")
	if err != nil {
		panic(err)
	}

	mongoDB = InitMongo(conf)
	requestURL := fmt.Sprintf("http://order:7373/order")
	client := http.Client{}

	order := Order{
		Product: "Product4",
		Price:   4.44,
		Client:  "tester4",
	}

	orderBytes, err := json.Marshal(&order)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(orderBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Add("plunder", key)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBytes))
}

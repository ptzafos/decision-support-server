package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var balance float64

type pandoraRequest struct {
	Key   string
	Close float64
	Macd  float64
	Rsi   float64
}

func decisionMaking(request pandoraRequest) {
	macd := request.Macd
	rsi := request.Rsi
	close := request.Close
	key := request.Key
	if rsi <= 30.0 && macd < 0.0 {
		buy(key, close)
	}
	if rsi >= 70.0 && macd > 0.0 {
		sell(key, close)
	}
}

func buy(key string, price float64) {
	balance -= price
	fmt.Println("Buy order:", key, "@", price, "- Current balance:", balance)
}

func sell(key string, price float64) {
	balance += price
	fmt.Println("Sell order:", key, "@", price, "- Current balance:", balance)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data pandoraRequest
		parseErr := json.NewDecoder(r.Body).Decode(&data)
		if parseErr != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			panic(parseErr)
		}
		decisionMaking(data)
		fmt.Fprint(w, "Pandora Analytics Server: Request Received!")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/pandora-analytics", postHandler)
	interrupt := make(chan os.Signal, 1)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()
	<-interrupt
}

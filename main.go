package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var balance float64
var sellTrades = 100
var buyTrades = 100
var mutex = &sync.Mutex{}

type pandoraRequest struct {
	Key   string  `json:"key"`
	Close float64 `json:"close"`
	Macd  float64 `json:"macd"`
	Rsi   float64 `json:"rsi"`
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
	if buyTrades > 0 {
		mutex.Lock()
		balance -= price
		buyTrades--
		checkAndresetVolTrades()
		mutex.Unlock()
		fmt.Println("Buy order:", key, "@", price, "- Current balance:", balance, "s:", sellTrades, "b:", buyTrades)
	}
}

func sell(key string, price float64) {
	if buyTrades < sellTrades && sellTrades > 0 {
		mutex.Lock()
		balance += price
		sellTrades--
		checkAndresetVolTrades()
		mutex.Unlock()
		fmt.Println("Sell order:", key, "@", price, "- Current balance:", balance, "s:", sellTrades, "b:", buyTrades)
	}
}

func checkAndresetVolTrades() {
	if sellTrades == 0 && buyTrades == 0 {
		sellTrades = 100
		buyTrades = 100
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data pandoraRequest
		parseErr := json.NewDecoder(r.Body).Decode(&data)
		if parseErr != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			panic(parseErr)
		}
		// strData, _ := json.Marshal(data)
		// fmt.Println(string(strData))
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

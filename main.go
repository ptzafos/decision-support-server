package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// balanceEur float32

type pandoraRequest struct {
	Key   string
	Close float64
	Macd  float64
	Rsi   float64
}

func decisionMaking() {}

func buy() {}

func sell() {}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data pandoraRequest
		parseErr := json.NewDecoder(r.Body).Decode(&data)
		if parseErr != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			panic(parseErr)
		}
		msg, _ := json.Marshal(data)
		fmt.Printf(string(msg))
		fmt.Fprint(w, "Pandora Analytics Server: Request Received!")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/pandora-analytics", postHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

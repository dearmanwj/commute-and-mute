package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Activity struct {
	Id           int64
	Type         string
	Start_latlng [2]float32
	End_latlng   [2]float32
}

func main() {
	http.HandleFunc("/", handlerHttp)
	http.ListenAndServe(":8080", nil)
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("In http handler")
	if r.Method == "GET" {
		challengeString := r.URL.Query().Get("hub.challenge")
		fmt.Fprintf(w, challengeString)
	} else {
		var a Activity
		json.NewDecoder(r.Body).Decode(&a)
		fmt.Fprintf(w, "Activity: %v", a)
	}
}

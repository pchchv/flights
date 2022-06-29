package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	r, err := json.Marshal("Flight Service. Version 0.1")
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(r)
	if err != nil {
		log.Panic(err)
	}
}

func toFlights(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(flightsJSON)
	if err != nil {
		log.Panic(err)
	}
}

func options(w http.ResponseWriter, req *http.Request) {
	// Need to add getting the best option
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(optionsJSON())
	if err != nil {
		log.Panic(err)
	}
}

func server() {
	log.Println("Server started!")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/flights", toFlights)
	http.HandleFunc("/options", options)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

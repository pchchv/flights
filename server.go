package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, _ *http.Request) {
	r, err := json.Marshal("Flight Service. Version 0.2")
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(r)
	if err != nil {
		log.Panic(err)
	}
}

func toFlights(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(flightsJSON)
	if err != nil {
		log.Panic(err)
	}
}

func variants(w http.ResponseWriter, req *http.Request) {
	opt := req.URL.Query().Get("options")
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(optionsJSON(opt))
	if err != nil {
		log.Panic(err)
	}
}

func server() {
	log.Println("Server started!")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/flights", toFlights)
	http.HandleFunc("/variants", variants)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

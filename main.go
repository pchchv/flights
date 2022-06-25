package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var flights [][]Flight
var flightsJSON []byte

type Flight struct {
	Carrier            string `xml:"Carrier"`
	FlightNumber       int    `xml:"FlightNumber"`
	Source             string `xml:"Source"`
	Destination        string `xml:"Destination"`
	DepartureTimeStamp string `xml:"DepartureTimeStamp"`
	ArrivalTimeStamp   string `xml:"ArrivalTimeStamp"`
	Class              string `xml:"Class"`
	NumberOfStops      string `xml:"NumberOfStops"`
	FareBasis          string `xml:"FareBasis"`
	TicketType         string `xml:"TicketType"`
}

func parseXML(file string) []Flight {
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panic(err)
		}
	}(f)
	decoder := xml.NewDecoder(f)
	var flights []Flight
	for {
		token, err := decoder.Token()
		if err != nil || token == nil {
			if err != io.EOF {
				log.Panic(err)
			}
			break
		}
		switch ty := token.(type) {
		case xml.StartElement:
			if ty.Name.Local == "Flight" {
				var f Flight
				err := decoder.DecodeElement(&f, &ty)
				if err != nil {
					log.Panic(err)
				}
				flights = append(flights, f)
			}
		}
	}
	return flights
}

func getFlights() [][]Flight {
	data := parseXML("RS_Via-3.xml")
	data = append(data, parseXML("RS_ViaOW.xml")...)
	var flights [][]Flight
	for _, v := range data {
		if len(flights) == 0 {
			flights = append(flights, []Flight{v})
		} else {
			for i, f := range flights {
				if strings.TrimSpace(f[0].FareBasis) == strings.TrimSpace(v.FareBasis) {
					flights[i] = append(flights[i], v)
					break
				}
				if i == len(flights)-1 {
					flights = append(flights, []Flight{v})
				}
			}
		}
	}
	var d [][]Flight
	for _, v := range flights {
		ft := 0
		for _, flight := range v {
			if flight.Source == "DXB" || flight.Destination == "BKK" {
				ft++
			}
		}
		if ft == 2 {
			d = append(d, v)
		}
	}
	return d
}

func getJSON(data [][]Flight) []byte {
	fl, err := json.MarshalIndent(data, " ", "\t")
	if err != nil {
		log.Panic(err)
	}
	return fl
}

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

func main() {
	flights = getFlights()
	flightsJSON = getJSON(flights)
	log.Println("Server started")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/flights", toFlights)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

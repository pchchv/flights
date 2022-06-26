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

var flights []Flights
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

type Price struct {
	ServiceCharges float64 `xml:"ServiceCharges"`
}

type Flights struct {
	flights []Flight
	price   Price
}

func parseXML(file string) ([]Flight, []Price) {
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
	var prices []Price
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
			if ty.Name.Local == "Pricing" {
				var p Price
				err := decoder.DecodeElement(&p, &ty)
				if err != nil {
					log.Panic(err)
				}
				prices = append(prices, p)
			}
		}
	}
	return flights, prices
}

func getFlights() []Flights {
	fData, pData := parseXML("RS_Via-3.xml")
	fd, pd := parseXML("RS_ViaOW.xml")
	fData = append(fData, fd...)
	pData = append(pData, pd...)
	var flights [][]Flight
	for _, v := range fData {
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
	var d []Flights
	if len(flights) == len(pData) {
		for i, fl := range flights {
			if fl[0].Source == "DXB" || fl[0].Destination == "BKK" {
				f := Flights{flights: flights[i], price: pData[i]}
				d = append(d, f)
			}
		}
	}
	return d
}

func getJSONArr(data []Flights) []byte {
	var fl []byte
	for _, v := range data {
		f, err := json.MarshalIndent(v.flights, " ", "\t")
		if err != nil {
			log.Panic(err)
		}
		p, err := json.MarshalIndent(v.price, " ", "\t")
		if err != nil {
			log.Panic(err)
		}
		f = append(f, p...)
		fl = append(fl, f...)
	}
	return fl
}

func getJSON(flight Flights, pre string) []byte {
	var fl []byte
	pr, err := json.MarshalIndent(pre, "\n", "\t")
	if err != nil {
		log.Panic(err)
	}
	f, err := json.MarshalIndent(flight.flights, "  ", "\t")
	if err != nil {
		log.Panic(err)
	}
	p, err := json.MarshalIndent(flight.price, "  ", "\t")
	if err != nil {
		log.Panic(err)
	}
	fl = append(fl, pr...)
	f = append(f, p...)
	fl = append(fl, f...)
	return fl
}

func prices() (Flights, Flights) {
	var expensive, cheap Flights
	for i, f := range flights {
		if i == 0 {
			cheap = f
		} else if i == 1 {
			if f.price.ServiceCharges < cheap.price.ServiceCharges {
				cheap, expensive = f, cheap
			} else if f.price.ServiceCharges > cheap.price.ServiceCharges {
				expensive = f
			}
		} else {
			if expensive.price.ServiceCharges < f.price.ServiceCharges {
				expensive = f
			}
			if cheap.price.ServiceCharges > f.price.ServiceCharges {
				cheap = f
			}
		}
	}
	return cheap, expensive
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

func options(w http.ResponseWriter, req *http.Request) {
	// Need to add getting the fastest/longest and best option
	ch, ex := prices()
	cheap := getJSON(ch, "The cheapest flight: ")
	expensive := getJSON(ex, "The most expensive flight: ")
	cap, err := json.MarshalIndent("Options:", "\n\n", "\n")
	if err != nil {
		log.Panic(err)
	}
	res := append(cap, cheap...)
	res = append(res, expensive...)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	flights = getFlights()
	flightsJSON = getJSONArr(flights)
	log.Println("Server started")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/flights", toFlights)
	http.HandleFunc("/options", options)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

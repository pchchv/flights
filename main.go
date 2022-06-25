package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"
)

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
	return flights
}

func allFlights(flights [][]Flight) {
	var data [][]Flight
	for _, v := range flights {
		ft := 0
		for _, flight := range v {
			if flight.Source == "DXB" || flight.Destination == "BKK" {
				ft++
			}
		}
		if ft == 2 {
			data = append(data, v)
		}
	}
}

func main() {
	flights := getFlights()
	allFlights(flights)
}

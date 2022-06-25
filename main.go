package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"
)

/*type Flight struct {
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
}*/

type Flight struct {
	FareBasis string `xml:"FareBasis"`
}

func parseXML(file string) []Flight {
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
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

func getFlights() []string {
	var flights []string
	data := parseXML("RS_Via-3.xml")
	data = append(data, parseXML("RS_ViaOW.xml")...)
	fl := make(map[Flight]int)
	for _, v := range data {
		fl[v]++
	}
	for f, _ := range fl {
		flights = append(flights, strings.TrimSpace(f.FareBasis))
	}
	return flights
}

func main() {
	log.Println(getFlights())
}

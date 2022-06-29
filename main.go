package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
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

func duration() (Flights, time.Duration, Flights, time.Duration) {
	var fastest, longest Flights
	var fdur, ldur time.Duration
	for i, f := range flights {
		dep, err := time.Parse("2006-01-02T1504", f.flights[0].DepartureTimeStamp)
		if err != nil {
			log.Panic(err)
		}
		ari, err := time.Parse("2006-01-02T1504", f.flights[len(f.flights)-1].ArrivalTimeStamp)
		if err != nil {
			log.Panic(err)
		}
		dur := ari.Sub(dep)
		if i == 0 {
			fdur = dur
			fastest = f
			continue
		} else if i == 1 {
			if fdur < dur {
				ldur = fdur
				longest = fastest
				fdur = dur
				fastest = f
			} else {
				ldur = dur
				longest = f
			}
			continue
		}
		if fdur < dur {
			fdur = dur
			fastest = f
		}
		if ldur > dur {
			ldur = dur
			longest = f
		}
	}
	return fastest, fdur, longest, ldur
}

func optionsJSON() []byte {
	fastest, fdur, longest, ldur := duration()
	ch, ex := prices()
	// Getting data on the cheapest and the most expensive flights
	cheap := getJSON(ch, "The cheapest flight: ")
	expensive := getJSON(ex, "The most expensive flight: ")

	// Getting data on the fastest and slowest flights
	fast := getJSON(fastest, "The fastest flight: ")
	f := fmt.Sprintf("Its duration: %v", fdur)
	fa, err := json.MarshalIndent(f, "\n\n", "\n")
	if err != nil {
		log.Panic(err)
	}
	fast = append(fast, fa...)
	long := getJSON(longest, "The longest flight: ")
	l := fmt.Sprintf("Its duration: %v", ldur)
	lo, err := json.MarshalIndent(l, "\n\n", "\n")
	if err != nil {
		log.Panic(err)
	}
	long = append(long, lo...)

	// Collecting data in one JSON
	cap, err := json.MarshalIndent("Options:", "\n\n", "\n")
	if err != nil {
		log.Panic(err)
	}
	res := append(cap, cheap...)
	res = append(res, expensive...)
	res = append(res, fast...)
	res = append(res, long...)
	return res
}

func main() {
	flights = getFlights()
	flightsJSON = getJSONArr(flights)
	server()
}

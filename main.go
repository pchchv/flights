package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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

type ByPrice []Flights

func (a ByPrice) Len() int           { return len(a) }
func (a ByPrice) Less(i, j int) bool { return a[i].price.ServiceCharges < a[j].price.ServiceCharges }
func (a ByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

func getDuration(flight Flights) time.Duration {
	dep, err := time.Parse("2006-01-02T1504", flight.flights[0].DepartureTimeStamp)
	if err != nil {
		log.Panic(err)
	}
	ari, err := time.Parse("2006-01-02T1504", flight.flights[len(flight.flights)-1].ArrivalTimeStamp)
	if err != nil {
		log.Panic(err)
	}
	return ari.Sub(dep)
}

func duration() (Flights, time.Duration, Flights, time.Duration) {
	var fastest, longest Flights
	var minDuration, maxDuration time.Duration
	for i, f := range flights {
		dur := getDuration(f)
		if i == 0 {
			minDuration = dur
			fastest = f
			continue
		} else if i == 1 {
			if minDuration < dur {
				maxDuration = minDuration
				longest = fastest
				minDuration = dur
				fastest = f
			} else {
				maxDuration = dur
				longest = f
			}
			continue
		}
		if minDuration < dur {
			minDuration = dur
			fastest = f
		}
		if maxDuration > dur {
			maxDuration = dur
			longest = f
		}
	}
	return fastest, minDuration, longest, maxDuration
}

func optimalFlight() Flights {
	var optimal Flights
	var optimalDuration time.Duration
	for i, f := range flights {
		if i == 0 {
			optimal = f
			optimalDuration = getDuration(f)
		} else if i <= 5 {
			dur := getDuration(f)
			if dur < optimalDuration {
				optimal = f
				optimalDuration = dur
			}
		}
		if i >= 5 {
			break
		}
	}
	return optimal
}

func optionsJSON(opt string) []byte {
	res, err := json.MarshalIndent("Options:", "\n\n", "\n")
	if err != nil {
		log.Panic(err)
	}
	if strings.Contains(opt, "price") || opt == "" {
		// Getting data on the cheapest and the most expensive flights
		ch := flights[0]
		ex := flights[len(flights)-1]
		cheap := getJSON(ch, "The cheapest flight: ")
		expensive := getJSON(ex, "The most expensive flight: ")
		res = append(res, cheap...)
		res = append(res, expensive...)
	}
	if strings.Contains(opt, "duration") || opt == "" {
		// Getting data on the fastest and slowest flights
		fastest, minDuration, longest, maxDuration := duration()
		fast := getJSON(fastest, "The fastest flight: ")
		f := fmt.Sprintf("Its duration: %v", minDuration)
		fa, err := json.MarshalIndent(f, "\n\n", "\n")
		if err != nil {
			log.Panic(err)
		}
		fast = append(fast, fa...)
		long := getJSON(longest, "The longest flight: ")
		l := fmt.Sprintf("Its duration: %v", maxDuration)
		lo, err := json.MarshalIndent(l, "\n\n", "\n")
		if err != nil {
			log.Panic(err)
		}
		long = append(long, lo...)
		res = append(res, fast...)
		res = append(res, long...)
	}
	if strings.Contains(opt, "optimal") || opt == "" {
		// Getting optimal flight
		optimal := getJSON(optimalFlight(), "Optimal flight:")
		res = append(res, optimal...)
	}
	return res
}

func main() {
	flights = getFlights()
	flightsJSON = getJSONArr(flights)
	sort.Sort(ByPrice(flights))
	server()
}

package main

import (
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServerPing(t *testing.T) {
	res, err := http.Get("http://127.0.0.1:8080/ping")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	b := string(body)
	if !strings.Contains(b, "Flight Service") {
		t.Fatal()
	}
}

func TestLoadPing(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/ping",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestAllFlights(t *testing.T) {
	res, err := http.Get("http://127.0.0.1:8080/flights")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
}

func TestLoadFlights(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/flights",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestVariants(t *testing.T) {
	params := []string{"",
		"?options=price",
		"?options=duration",
		"?options=optimal",
		"?options=duration,optimal",
		"?options=duration,price",
		"?options=optimal,price",
		"?options=duration,optimal,price"}
	for _, v := range params {
		res, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/variants" + v))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("status not OK")
		}
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
	}
}

func TestLoadVariants(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/variants",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

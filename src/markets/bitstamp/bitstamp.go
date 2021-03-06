package bitstamp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// {"high": "2559.98",
//  "last": "2559.25",
//  "timestamp": "1497016409",
//  "bid": "2555.00",
//  "vwap": "2492.13",
//  "volume": "1650.07586432",
//  "low": "2423.51",
//  "ask": "2559.25",
//  "open": "2519.00"}

type TBitstampTicker struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

type MyError struct {
	errcode int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Error code: %d", e.errcode)
}

// coinpair = btceur
func DoGet(coinpair string) (float64, error) {
	// Build the URL
	url := fmt.Sprintf("https://www.bitstamp.net/api/v2/ticker/%s/", coinpair)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("ERROR1")
		// log.Fatal("NewRequest: ", err)
		return -1, &MyError{13}
	}

	// Build HTTP client
	client := &http.Client{}

	// Do send an HTTP request and return an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR2")
		// log.Fatal("Do: ", err)
		return -1, &MyError{13}
	}

	var bst TBitstampTicker

	// Build JSON decoder with resp.Body
	decoder := json.NewDecoder(resp.Body)

	// Fill the record with JSON data
	err = decoder.Decode(&bst)
	if err != nil {
		fmt.Println("ERROR3")
		log.Println(err)
		return -1, &MyError{13}
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Check not empty
	if bst.Last == "" {
		fmt.Println("ERROR4")
		return -1, &MyError{13}
	}

	price, err := strconv.ParseFloat(bst.Last, 64)
	if err != nil {
		fmt.Println("ERROR5")
		// log.Fatal("ParseFloat: ", err)
		return -1, &MyError{13}
	}

	/*
		    volume, err := strconv.ParseFloat(bst.Volume, 64)
			if err != nil {
				log.Fatal("ParseFloat: ", err)
		        return -1, &MyError{13}
			}*/

	return price, nil
}

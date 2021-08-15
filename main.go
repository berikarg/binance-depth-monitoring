//!TODO figure out how to check if symbol exists

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var baseUrl = "https://api.binance.com" //consider putting inside makeGetRequest

type Depth struct {
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
	BidsOrderSum float64
	AsksOrderSum float64
}

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	for {
		depth1, err := getDepth("DASHUSDT", 15)
		if err != nil {
			log.Fatal(err)
			return
		}
		bids, err := json.MarshalIndent(depth1.Bids, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Bids: ", string(bids))

		asks, err := json.MarshalIndent(depth1.Asks, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Asks: ", string(asks))

		fmt.Println("Bids Order Sum: ", depth1.BidsOrderSum)
		fmt.Println("Asks Order Sum: ", depth1.AsksOrderSum)
		time.Sleep(time.Second)
	}
}

// returns depth of a symbol for limits between 1 and 100
func getDepth(symbol string, limit int) (Depth, error) {
	var endUrl string
	depth1 := Depth{}
	url := baseUrl + "/api/v3/depth"

	if symbol == "" {
		return depth1, errors.New("empty symbol")
	}
	if limit <= 0 || limit > 100 {
		return depth1, errors.New("invalid limit")
	}
	endUrl = "?symbol=" + symbol // if symbol is not 5, 10, 20, 50 it will return 100
	endUrl = endUrl + "&limit=" + strconv.Itoa(limit)
	url = url + endUrl

	body, err := makeGetRequest(url)
	if err != nil {
		log.Fatal(err)
		return depth1, err
	}
	err = json.Unmarshal(body, &depth1)
	if err != nil {
		log.Fatal(err)
		return depth1, err
	}

	if len(depth1.Bids) > limit {
		depth1.Bids = depth1.Bids[:limit]
		depth1.Asks = depth1.Asks[:limit]
	}

	//calculate the sum of bid orders
	for _, v := range depth1.Bids {
		bidOrder, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			log.Fatal(err)
			return depth1, err
		}
		depth1.BidsOrderSum += bidOrder
	}

	//calculate the sum of ask orders
	for _, v := range depth1.Asks {
		askOrder, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			log.Fatal(err)
			return depth1, err
		}
		depth1.AsksOrderSum += askOrder
	}

	return depth1, err
}

func getExchangeInfo(symbol string) ([]byte, error) {
	var endUrl string
	url := baseUrl + "/api/v3/exchangeInfo"
	if symbol == "" {
		endUrl = ""
	} else {
		endUrl = "?symbol=" + symbol
	}
	url = url + endUrl
	return makeGetRequest(url)
}

func makeGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body, err
}

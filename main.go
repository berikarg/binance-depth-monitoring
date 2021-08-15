package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var baseUrl = "https://api.binance.com" //consider putting inside makeGetRequest

type Depth struct {
	LastUpdateId int         `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}

func main() {
	body, err := getDepth("BTCUSDT", 5)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	depth1 := Depth{}
	err = json.Unmarshal(body, &depth1)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	fmt.Println(depth1.Asks[0][1])
}

func getDepth(symbol string, limit int) ([]byte, error) {
	var endUrl string
	url := baseUrl + "/api/v3/depth"
	if symbol == "" {
		return nil, errors.New("empty symbol")
	}
	if limit <= 0 && limit > 5000 {
		return nil, errors.New("invalid limit")
	}
	endUrl = "?symbol=" + symbol
	endUrl = endUrl + "&limit=" + strconv.Itoa(limit)
	url = url + endUrl
	return makeGetRequest(url)
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

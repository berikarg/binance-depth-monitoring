/*
This package includes functions to get the depth of a particular
symbol pair from Binance. It has functions to acquire depth via REST API and via Websocket
*/
package binancedepth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Depth struct {
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
	BidsOrderSum float64
	AsksOrderSum float64
}

// returns depth of a symbol for limits between 1 and 100
// uses REST API
func GetDepth(symbol string, limit int) (Depth, error) {
	var endUrl string
	depth1 := Depth{}
	url := "https://api.binance.com/api/v3/depth"

	if limit <= 0 || limit > 100 {
		err := errors.New("invalid limit range")
		log.Print(err)
		return depth1, err
	}
	endUrl = "?symbol=" + symbol // if symbol is not 5, 10, 20, 50 it will return 100
	endUrl = endUrl + "&limit=" + strconv.Itoa(limit)
	url = url + endUrl

	body, err := makeGetRequest(url)
	if err != nil {
		log.Print(err)
		return depth1, err
	}
	err = json.Unmarshal(body, &depth1)
	if err != nil {
		log.Print(err)
		return depth1, err
	}

	if len(depth1.Bids) > limit {
		depth1.Bids = depth1.Bids[:limit]
		depth1.Asks = depth1.Asks[:limit]
	}

	calcOrderSum(&depth1)

	return depth1, err
}

func makeGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// notify if response came with any error code
	if resp.StatusCode >= 400 {
		log.Println(string(body))
		return nil, errors.New(string(body))
	}
	return body, err
}

//calculate the sum of bid and asks orders
func calcOrderSum(depth1 *Depth) {
	for _, v := range depth1.Bids {
		bidOrder, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			log.Print(err)
			return
		}
		depth1.BidsOrderSum += bidOrder
	}
	for _, v := range depth1.Asks {
		askOrder, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			log.Print(err)
			return
		}
		depth1.AsksOrderSum += askOrder
	}
}

// initiates a websocket connection to get
// depth of a symbol for limits between 1 and 20
func DialDepth(symbol string, limit int) (*websocket.Conn, error) {
	var endUrl string
	var level int //to be used in request
	url := "wss://stream.binance.com:9443/ws/"

	if limit <= 0 || limit > 20 {
		err := errors.New("invalid limit range")
		log.Print(err)
		return nil, err
	}
	// adjust level
	if limit <= 5 {
		level = 5
	} else if limit <= 10 {
		level = 10
	} else {
		level = 20
	}
	endUrl = strings.ToLower(symbol)
	endUrl = endUrl + "@depth" + strconv.Itoa(level)
	url = url + endUrl

	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// notify if response came with any error code
	if resp.StatusCode >= 400 {
		log.Println(string(body))
		return nil, errors.New(string(body))
	}

	return conn, err
}

// reads from websocket stream returned by DialDepth
// limit is needed to cut extra orders
func ReadDepth(conn *websocket.Conn, limit int) (Depth, error) {
	depth := Depth{}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return depth, err
	}

	err = json.Unmarshal(msg, &depth)
	if err != nil {
		log.Println(err)
		return depth, err
	}

	if len(depth.Bids) > limit {
		depth.Bids = depth.Bids[:limit]
		depth.Asks = depth.Asks[:limit]
	}

	calcOrderSum(&depth)

	return depth, err
}

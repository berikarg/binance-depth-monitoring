package binance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Depth struct {
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
	BidsOrderSum float64
	AsksOrderSum float64
}

// returns depth of a symbol for limits between 1 and 100
func GetDepth(c *gin.Context) {
	var endUrl string
	depth1 := Depth{}
	url := "https://api.binance.com/api/v3/depth"
	symbol := c.Param("symbol")
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid parameter limit"})
		return
	}
	if limit <= 0 || limit > 100 {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid limit range"})
		return
	}
	endUrl = "?symbol=" + symbol // if symbol is not 5, 10, 20, 50 it will return 100
	endUrl = endUrl + "&limit=" + strconv.Itoa(limit)
	url = url + endUrl

	body, err := makeGetRequest(url)
	if err != nil {
		log.Print(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(body, &depth1)
	if err != nil {
		log.Print(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(depth1.Bids) > limit {
		depth1.Bids = depth1.Bids[:limit]
		depth1.Asks = depth1.Asks[:limit]
	}

	calcOrderSum(&depth1)

	c.IndentedJSON(http.StatusOK, depth1)
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

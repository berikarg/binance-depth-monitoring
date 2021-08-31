package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/binancedepth"
)

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)

	http.HandleFunc("/depth", showDepth)
	http.ListenAndServe(":8080", nil)
}

func showDepth(w http.ResponseWriter, r *http.Request) {
	depth, err := binancedepth.GetDepth("BTCUSDT", 5)
	if err != nil {
		log.Println(err)
		return
	}
	depthJson, err := json.Marshal(depth)
	if err != nil {
		log.Println(err)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(depthJson)
}

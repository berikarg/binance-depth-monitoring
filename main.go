package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/binance-api-test/websocket"
)

func depth(w http.ResponseWriter, r *http.Request) {
	// we call our new websocket package Upgrade
	// function in order to upgrade the connection
	// from a standard HTTP connection to a websocket one
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	// we then call our Writer function
	// which continually polls and writes the results
	// to this websocket connection
	go websocket.Writer(ws)
}

func setupRoutes() {
	http.HandleFunc("/depth", depth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)

	setupRoutes()
}

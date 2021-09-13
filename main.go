package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"example.com/binancedepth"
	"github.com/gorilla/websocket"
)

//Globals
var (
	prevSymbol    string
	prevLimit     int
	binanceWsConn *websocket.Conn
)

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)

	http.HandleFunc("/", index)
	http.HandleFunc("/depth", showDepth)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func showDepth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	symbol := r.FormValue("symbol")
	limitStr := r.FormValue("limit")
	symbol = strings.TrimPrefix(symbol, "symbol=") //for some reason ParseForm() unnecessarily adds keys
	limitStr = strings.TrimPrefix(limitStr, "limit=")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resp, err := json.Marshal(map[string]string{"message": "Symbol or limit are invalid"})
		if err != nil {
			log.Println("Marshal:", err)
		}
		w.Write(resp)
		return
	}

	//at start or if requested symbol,limit pair has been changed, establish a new WS connection
	if symbol != prevSymbol || limit != prevLimit {
		binanceWsConn, err = binancedepth.DialDepth(symbol, limit)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			resp, _ := json.Marshal(map[string]string{"message": "Binance: Symbol or limit are invalid"})
			w.Write(resp)
			return
		}
	}

	depth, err := binancedepth.ReadDepth(binanceWsConn, limit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	depthJson, err := json.Marshal(depth)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(depthJson)
}

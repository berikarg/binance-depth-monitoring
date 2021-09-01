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
	r.ParseForm()
	symbol := r.FormValue("symbol")
	limitStr := r.FormValue("limit")
	symbol = strings.TrimPrefix(symbol, "symbol=") //for some reason ParseForm() unnecessarily adds keys
	limitStr = strings.TrimPrefix(limitStr, "limit=")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println(err)
		return
	}
	depth, err := binancedepth.GetDepth(symbol, limit)
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
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-With")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(depthJson)
}

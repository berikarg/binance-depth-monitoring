package main

import (
	"fmt"
	"log"
	"os"

	"example.com/binance-api-test/binance"
	"github.com/gin-gonic/gin"
)

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(logFile)

	router := gin.Default()
	router.GET("/:symbol/:limit", binance.GetDepth)

	router.Run("localhost:8080")
}

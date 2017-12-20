package main

import (
	"github.com/gin-gonic/gin"
	"github.com/matejl/challenge/handlers"
	"flag"
	"strconv"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	port := flag.Int("port", 8080, "Which port should the server listen to.")
	flag.Parse()

	// example 1: /campaign?id=1234&dimensions=adId,adName&metrics=impressions,interactions,swipes
	// example 2: /campaign?dateRange=lastWeek&dimensions=adId,adName&metrics=uniqueUsers,impressions
	r.GET("/campaign", handlers.Campaign)

	r.Run(":" + strconv.Itoa(*port)) // listen and serve on 0.0.0.0:8080

}
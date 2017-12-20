package main

import (
	"github.com/gin-gonic/gin"
	"github.com/matejl/challenge/handlers"
)

func main() {

	r := gin.Default()

	// example 1: /campaign?id=1234&dimensions=adId,adName&metrics=impressions,interactions,swipes
	// example 2: /campaign?dateRange=lastWeek&dimensions=adId,adName&metrics=uniqueUsers,impressions
	r.GET("/campaign", handlers.Campaign)

	r.Run() // listen and serve on 0.0.0.0:8080

}
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"fmt"
	"github.com/matejl/challenge/model"
)

const (
	DateRangeLastWeek = "lastWeek"
)

type params struct {
	Id         int64    `form:"id"`
	Dimensions []string `form:"dimensions[]"`
	Metrics    []string `form:"metrics[]"`
	DateRange  string   `form:"dateRange"`
}

func Campaign(c *gin.Context) {

	p := params{}
	err := c.Bind(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if p.Id == 0 && p.DateRange == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either id or dateRange are mandatory"})
		return
	}

	// id based request
	if p.Id != 0 {
		if result, httpStatus, err := getCampaign(p.Id, p.Dimensions, p.Metrics); err == nil {
			c.JSON(http.StatusOK, result)
			return
		} else {
			c.JSON(httpStatus, gin.H{"error": err.Error()})
			return
		}
	}

	// date based request
	var dateFrom, dateTo time.Time
	switch p.DateRange {
	case DateRangeLastWeek:
		dateFrom = time.Now().AddDate(0, 0, -7)
		dateTo = time.Now()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unknown date range '%s'", p.DateRange)})
		return
	}

	if result, httpStatus, err := getAllCampaigns(dateFrom, dateTo, p.Dimensions, p.Metrics); err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(httpStatus, gin.H{"error": err.Error()})
	}
}

func getCampaign(campaignId int64, dimensions []string, metrics []string) ([]map[string]interface{}, int, error) {

	result, err := model.GetCampaignStats(campaignId, dimensions, metrics)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return result, http.StatusOK, nil

}

func getAllCampaigns(dateFrom time.Time, dateTo time.Time, dimensions []string, metrics []string) ([]map[string]interface{}, int, error) {

	result, err := model.GetCampaignsDateRangeStats(dateFrom, dateTo, dimensions, metrics)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return result, http.StatusOK, nil

}
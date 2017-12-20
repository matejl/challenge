package model

import (
	"cloud.google.com/go/bigquery"
	"github.com/matejl/challenge/database"
	"context"
	"github.com/matejl/challenge/utils/bigQueryUtils"
	"strings"
	"fmt"
	"errors"
	"time"
	"google.golang.org/api/iterator"
	"strconv"
	"log"
)

type Campaign struct {
	CampaignId   int64
	CampaignName string
}

func (c *Campaign) Save() (map[string]bigquery.Value, string, error) {

	return map[string]bigquery.Value{
		"campaign_id":   c.CampaignId,
		"campaign_name": c.CampaignName,
	}, "", nil

}

func MaxCampaignId() (int64, error) {

	client, err := database.GetStat()
	if err != nil {
		return 0, err
	}

	q := client.Query("SELECT MAX(campaign_id) AS m FROM dataset.campaign LIMIT 1")
	q.Dst = &bigquery.Table{}
	it, err := q.Read(context.Background())
	if err != nil {
		return 0, err
	}

	return bigQueryUtils.SingleRowInt64(it)

}

// GetCampaignStats returns campaign statistics depending on campaign id, dimensions and metrics.
func GetCampaignStats(campaignId int64, dimensions []string, metrics []string) ([]map[string]interface{}, error) {

	columns, groupBys, err := getColumnsGroupBys(dimensions, metrics)
	if err != nil {
		return nil, err
	}

	sqlQuery := `SELECT ` + strings.Join(columns, ",") + ` FROM dataset.campaign c
				LEFT JOIN dataset.ad a ON a.campaign_id = c.campaign_id
				LEFT JOIN dataset.user u ON u.ad_id = a.ad_id
				LEFT JOIN dataset.impression i ON i.user_id = u.user_id
				WHERE c.campaign_id = ` + strconv.FormatInt(campaignId, 10) + `
				GROUP BY ` + strings.Join(groupBys, ",")

	return getCampaignResults(sqlQuery, dimensions, metrics)

}

// GetCampaignsDateRangeStats returns statistics for multiple campaigns depending on date range, dimensions and metrics.
func GetCampaignsDateRangeStats(dateFrom time.Time, dateTo time.Time, dimensions []string, metrics []string) ([]map[string]interface{}, error) {

	columns, groupBys, err := getColumnsGroupBys(dimensions, metrics)
	if err != nil {
		return nil, err
	}

	sqlQuery := `SELECT ` + strings.Join(columns, ",") + ` FROM dataset.campaign c
				LEFT JOIN dataset.ad a ON a.campaign_id = c.campaign_id
				LEFT JOIN dataset.user u ON u.ad_id = a.ad_id
				LEFT JOIN dataset.impression i ON i.user_id = u.user_id
				WHERE DATE(i.datetime) >= '` + dateFrom.Format("2006-01-02") + `' AND
					DATE(i.datetime) <= '` + dateTo.Format("2006-01-02") + `'
				GROUP BY ` + strings.Join(groupBys, ",")

	return getCampaignResults(sqlQuery, dimensions, metrics)

}

func getColumnsGroupBys(dimensions []string, metrics []string) ([]string, []string, error) {
	columns := make([]string, 0)
	groupBys := make([]string, 0)

	for _, dimension := range dimensions {
		switch dimension {
		case "date":
			columns = append(columns, "i.datetime AS [date]")
			groupBys = append(groupBys, "date")
		case "campaignId":
			columns = append(columns, "c.campaign_id AS [campaignId]")
			groupBys = append(groupBys, "campaignId")
		case "campapignName":
			columns = append(columns, "c.campaign_name [campaignName]")
			groupBys = append(groupBys, "campaignName")
		case "adId":
			columns = append(columns, "a.ad_id AS [adId]")
			groupBys = append(groupBys, "adId")
		case "adName":
			columns = append(columns, "a.ad_name AS [adName]")
			groupBys = append(groupBys, "adName")
		default:
			return nil, nil, errors.New(fmt.Sprintf("unknown dimension '%s'", dimension))
		}

	}

	for _, metric := range metrics {
		switch metric {
		case "impressions":
			columns = append(columns, "COUNT(i.impression_id) AS [impressions]")
		case "clicks":
			columns = append(columns, "SUM(i.click) AS [clicks]")
		case "swipes":
			columns = append(columns, "SUM(i.swipe) AS [swipes]")
		case "pinches":
			columns = append(columns, "SUM(i.pinch) AS [pinches]")
		case "touches":
			columns = append(columns, "SUM(i.touch) AS [touches]")
		case "uniqueUsers":
			columns = append(columns, "COUNT(DISTINCT u.user_id) AS [uniqueUsers]")
		default:
			return nil, nil, errors.New(fmt.Sprintf("unknown metric '%s'", metric))
		}
	}

	return columns, groupBys, nil
}

func getCampaignResults(sqlQuery string, dimensions []string, metrics []string) ([]map[string]interface{}, error) {

	client, err := database.GetStat()
	if err != nil {
		return nil, err
	}

	q := client.Query(sqlQuery)
	q.UseLegacySQL = true
	iter, err := q.Read(context.Background())
	if err != nil {
		log.Println(sqlQuery)
		log.Println(err)
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	allKeys := append(dimensions, metrics...)
	for {
		result := make(map[string]interface{})
		var row []bigquery.Value
		err := iter.Next(&row)
		if err == iterator.Done {
			return results, nil
		}
		if err != nil {
			return nil, err
		}

		for i, val := range row {
			result[allKeys[i]] = val
		}
		results = append(results, result)
	}
	return results, nil

}

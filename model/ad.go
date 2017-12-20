package model

import (
	"cloud.google.com/go/bigquery"
	"github.com/matejl/challenge/utils/bigQueryUtils"
	"github.com/matejl/challenge/database"
	"golang.org/x/net/context"
)

type Ad struct {
	AdId       int64
	CampaignId int64
	AdName     string
}

func (a *Ad) Save() (map[string]bigquery.Value, string, error) {

	return map[string]bigquery.Value{
		"ad_id": a.AdId,
		"campaign_id": a.CampaignId,
		"ad_name": a.AdName,
	}, "", nil

}

func MaxAdId() (int64, error) {

	client, err := database.GetStat()
	if err != nil {
		return 0, err
	}

	q := client.Query("SELECT MAX(ad_id) AS m FROM [celtra-assignment:dataset.ad] LIMIT 1")
	q.Dst = &bigquery.Table{}
	it, err := q.Read(context.Background())
	if err != nil {
		return 0, err
	}

	return bigQueryUtils.SingleRowInt64(it)

}
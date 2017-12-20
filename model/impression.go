package model

import (
	"time"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/matejl/challenge/utils/bigQueryUtils"
	"github.com/matejl/challenge/database"
	"golang.org/x/net/context"
)

type Impression struct {
	ImpressionId int64
	UserId       int64
	Click        int
	Swipe        int
	Pinch        int
	Touch        int
	Datetime     time.Time
}

func (im *Impression) Save() (map[string]bigquery.Value, string, error) {

	return map[string]bigquery.Value{
		"impression_id": im.ImpressionId,
		"user_id":       im.UserId,
		"click":         im.Click,
		"swipe":         im.Swipe,
		"pinch":         im.Pinch,
		"touch":         im.Touch,
		"datetime":      civil.DateOf(im.Datetime),
	}, "", nil

}

func MaxImpressionId() (int64, error) {

	client, err := database.GetStat()
	if err != nil {
		return 0, err
	}

	q := client.Query("SELECT MAX(impression_id) AS m FROM [celtra-assignment:dataset.impression] LIMIT 1")
	q.Dst = &bigquery.Table{}
	it, err := q.Read(context.Background())
	if err != nil {
		return 0, err
	}

	return bigQueryUtils.SingleRowInt64(it)

}
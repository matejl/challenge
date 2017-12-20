package model

import (
	"cloud.google.com/go/bigquery"
	"github.com/matejl/challenge/utils/bigQueryUtils"
	"github.com/matejl/challenge/database"
	"golang.org/x/net/context"
)

type User struct {
	UserId   int64
	AdId     int64
	UserHash string
}

func (u *User) Save() (map[string]bigquery.Value, string, error) {

	return map[string]bigquery.Value{
		"user_id": u.UserId,
		"ad_id": u.AdId,
		"user_hash": u.UserHash,
	}, "", nil

}

func MaxUserId() (int64, error) {

	client, err := database.GetStat()
	if err != nil {
		return 0, err
	}

	q := client.Query("SELECT MAX(user_id) AS m FROM [celtra-assignment:dataset.user] LIMIT 1")
	q.Dst = &bigquery.Table{}
	it, err := q.Read(context.Background())
	if err != nil {
		return 0, err
	}

	return bigQueryUtils.SingleRowInt64(it)

}
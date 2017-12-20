package database

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
)

const (
	ProjectId = "celtra-assignment"
)

func GetStat() (*bigquery.Client, error){
	return bigquery.NewClient(context.Background(), ProjectId)
}

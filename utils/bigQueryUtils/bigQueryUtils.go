package bigQueryUtils

import (
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func SingleRowInt64(it *bigquery.RowIterator) (int64, error) {
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}

		if valInt64, valInt64Ok := row[0].(int64); valInt64Ok {
			return valInt64, nil
		}
	}
	return 0, nil
}
package es

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// build generic elastsicsearch index
func BuildESIndex(ctx context.Context, client *elastic.Client, index, id string, doc interface{}) error {
	_, err := client.Index().
		Index(index).
		Id(id).
		BodyJson(doc).
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		return fmt.Errorf("es index error: %w", err)
	}
	return nil
}

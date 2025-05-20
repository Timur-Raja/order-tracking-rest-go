package orderesrc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/es"

	"github.com/olivere/elastic/v7"
)

// -esrc suffix for elasticsearch, files contain elasticsearch related code

type orderIndexer struct {
	Client    *elastic.Client
	IndexName string
}

func NewOrderIndexer(client *elastic.Client, indexName string) *orderIndexer {
	return &orderIndexer{
		Client:    client,
		IndexName: indexName,
	}
}

// indexes the given OrderView into Elasticsearch.
// Uses the order ID as the ES document ID and forces a refresh so the document
func (i *orderIndexer) Run(ctx context.Context, o *order.OrderView) error {
	return es.BuildESIndex(ctx, i.Client, i.IndexName, fmt.Sprintf("%d", o.ID), o)
}

/////////////////////////////////////////////////////////////////////////////////////////////

type ordersViewSearchQuery struct {
	Client *elastic.Client
	Index  string

	Params struct {
		Search    string     `form:"q"`
		StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
		EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	}

	Result struct {
		OrderViewList []*order.OrderView
	}
}

func NewOrdersViewSearchQuery(client *elastic.Client, indexName string) *ordersViewSearchQuery {
	return &ordersViewSearchQuery{
		Client: client,
		Index:  indexName,
	}
}

// run text search and date filtering query
func (q *ordersViewSearchQuery) Run(ctx context.Context) error {
	// build bool query
	boolQ := elastic.NewBoolQuery()

	// text search on relevant fields
	if q.Params.Search != "" {
		boolQ.Must(elastic.NewMultiMatchQuery(q.Params.Search,
			"userName", "userEmail", "shippingAddress", "status", "orderItems.productName",
		))
	}

	// date range filter
	if q.Params.StartDate != nil || q.Params.EndDate != nil {
		r := elastic.NewRangeQuery("createdAt")
		if q.Params.StartDate != nil {
			r = r.Gte(q.Params.StartDate.Format(time.RFC3339))
		}
		if q.Params.EndDate != nil {
			next := q.Params.EndDate.AddDate(0, 0, 1)
			r = r.Lt(next.Format(time.RFC3339))
		}
		boolQ.Filter(r)
	}

	// setup query
	queryService := q.Client.Search().Index(q.Index).Sort("createdAt", false).Size(100)
	if q.Params.Search == "" && q.Params.StartDate == nil && q.Params.EndDate == nil {
		// fetch all
		queryService = queryService.Query(elastic.NewMatchAllQuery())
	} else {
		// fetch by params
		queryService = queryService.Query(boolQ)
	}

	// execute
	res, err := queryService.Do(ctx)
	if err != nil {
		return fmt.Errorf("search error: %w", err)
	}

	// load hits into OrderView
	for _, hit := range res.Hits.Hits {
		var o order.OrderView
		if err := json.Unmarshal(hit.Source, &o); err != nil {
			return fmt.Errorf("unmarshal ES hit: %w", err)
		}
		q.Result.OrderViewList = append(q.Result.OrderViewList, &o)
	}

	return nil
}

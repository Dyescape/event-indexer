package elastic

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/viper"
)

type ElasticSearchClient struct {
	client  *elasticsearch.Client
	index   string
	refresh string
}

func NewElasticSearchClient() *ElasticSearchClient {
	cfg := elasticsearch.Config{
		Addresses: viper.GetStringSlice("elasticsearch.address"),
		// ...
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return &ElasticSearchClient{
		client:  es,
		index:   viper.GetString("elasticsearch.index"),
		refresh: viper.GetString("elasticsearch.refresh"),
	}
}

func (e *ElasticSearchClient) Index(document string) error {
	request := esapi.IndexRequest{
		Index:   e.index,
		Body:    strings.NewReader(document),
		Refresh: e.refresh,
	}

	response, err := request.Do(context.Background(), e.client)
	if response != nil {
		defer response.Body.Close()
	}
	return err
}

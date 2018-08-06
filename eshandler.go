package main

import (
	"context"
	"log"

	"github.com/olivere/elastic"
)

var client *elastic.Client

// NewESClient create a new connection to es server
func NewESClient(config *Config) (*elastic.Client, error) {
	// https://github.com/olivere/elastic/issues/312
	// due to run a elasticsearch in docker container
	// need use elastic.SetSniff(false) or NewSimpleClient

	log.Printf("Try connect elastiSearch [%s]", config.EleasticSearch.URL)
	url := elastic.SetURL(config.EleasticSearch.URL)
	client, err := elastic.NewSimpleClient(url)
	if err != nil {
		return nil, err
	}
	config.ESClient = client
	info, code, err := client.Ping(config.EleasticSearch.URL).Do(context.Background())

	if err != nil {
		// Handle error
		return nil, err
	}
	log.Printf("Connect ElastiSearch [%s] Success [code=%d,info=%+v]", config.EleasticSearch.URL, code, info)
	return client, nil
}

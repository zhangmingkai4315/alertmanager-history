package main

import (
	"fmt"

	"github.com/olivere/elastic"
)

// NewESClient create a new connection to es server
func NewESClient(config *Config) (*elastic.Client, error) {
	connectType := "http"
	if config.ESHTTPS {
		connectType = "https"
	}
	url := elastic.SetURL(fmt.Sprintf("%s://%s:%d", connectType, config.ESHost, config.ESPort))
	client, err := elastic.NewClient(url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

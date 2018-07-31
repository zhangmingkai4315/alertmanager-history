package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/olivere/elastic"

	"github.com/prometheus/alertmanager/types"
)

// SearchRequest define user posted json message
type SearchRequest struct {
	Term string `json:"term"`
	Skip int    `json:"skip"`
	Take int    `json:"take"`
}

// SearchResponse define server reply json message
type SearchResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
	Code  int         `json:"status_code"`
}

func jsonQueryResponse(w http.ResponseWriter, data interface{}, errMsg string, stautsCode int) {
	search := SearchResponse{
		Data:  data,
		Error: errMsg,
		Code:  stautsCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stautsCode)
	err := json.NewEncoder(w).Encode(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	return
}

func searchHandlerWithConfig(config *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searchRequest := SearchRequest{
			Take: 20,
			Skip: 0,
		}
		if r.Body == nil {
			jsonQueryResponse(w, nil, "Post data not valid", http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&searchRequest)
		if err != nil {
			jsonQueryResponse(w, nil, "Decode post data fail", http.StatusBadRequest)
			return
		}
		results, err := searchFromElasticSearch(&searchRequest, config)
		if err != nil {
			jsonQueryResponse(w, nil, err.Error(), http.StatusServiceUnavailable)
			return
		}
		jsonQueryResponse(w, results, "", http.StatusOK)

	}
}

func searchFromElasticSearch(search *SearchRequest, config *Config) (results []types.Alert, err error) {
	if config.ESClient == nil {
		err = errors.New("ElasticSearch connection client not ready")
		return
	}
	multiQuery := elastic.NewMultiMatchQuery(search.Term, "labels.alertname", "labels.instance", "labels.job", "annotations.description", "annotations.summary")

	searchResults, err := config.ESClient.Search().Index(config.EleasticSearch.IndexName).Query(multiQuery).From(search.Skip).Size(search.Take).Do(context.Background())
	if err != nil {
		return
	}
	for _, hit := range searchResults.Hits.Hits {
		var alert types.Alert
		err := json.Unmarshal(*hit.Source, &alert)
		if err != nil {
			continue
		}
		results = append(results, alert)
	}

	return
}

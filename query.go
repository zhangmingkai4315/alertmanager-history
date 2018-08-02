package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/olivere/elastic"

	"github.com/prometheus/alertmanager/types"
)

// SearchRequest define user posted json message
type SearchRequest struct {
	Term        string    `json:"term"`        // Default value is ""
	Skip        int       `json:"skip"`        // Default value is 0
	Size        int       `json:"size"`        // Default value is 20
	StartsAt    time.Time `json:"startsAt"`    // Default value is one week ago
	EndsAt      time.Time `json:"endsAt"`      // Default value is now
	SortBy      string    `json:"sortBy"`      // Default value is startsAt
	SortReverse bool      `json:"sortReverse"` // Default value is true
	Current     bool      `json:"current"`     // Default value is false
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "GET" {
			jsonQueryResponse(w, nil, "Post method only", http.StatusMethodNotAllowed)
			return
		}
		searchRequest := SearchRequest{
			Size:        20,
			Skip:        0,
			Term:        "",
			StartsAt:    time.Now().AddDate(0, 0, -7),
			EndsAt:      time.Now(),
			SortBy:      "startsAt",
			SortReverse: true,
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

	queryBuilder := elastic.NewBoolQuery()

	if search.Current == false {
		endsTimeRange := elastic.NewRangeQuery("endsAt").Lte(search.StartsAt)
		startsTimeRange := elastic.NewRangeQuery("startsAt").Gte(search.EndsAt)
		queryBuilder.MustNot(endsTimeRange)
		queryBuilder.MustNot(startsTimeRange)
	} else {
		endsTimeRange := elastic.NewRangeQuery("endsAt").Gte(time.Now().Add(time.Minute * -2))
		queryBuilder.Must(endsTimeRange)
	}

	var matcher elastic.Query
	if search.Term != "" {
		matcher = elastic.NewMultiMatchQuery(search.Term, "labels.alertname", "labels.instance", "labels.job", "annotations.description", "annotations.summary")
		queryBuilder = queryBuilder.Must(matcher)
	}

	searchResults, searchError := config.ESClient.
		Search().
		Index(config.EleasticSearch.IndexName).
		Query(queryBuilder).
		Sort(search.SortBy, search.SortReverse).
		From(search.Skip).
		Size(search.Size).
		Pretty(true).
		Do(context.Background())
	if searchError != nil {
		switch {
		case elastic.IsNotFound(searchError):
			err = errors.New("Resouce Not Found")
		case elastic.IsTimeout(searchError):
			err = errors.New("Search Time out")
		case elastic.IsConnErr(searchError):
			err = errors.New("Connect elasticsearch error")
			return
		default:
			err = searchError
		}
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

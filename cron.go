package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/olivere/elastic"

	"github.com/prometheus/alertmanager/types"
	"github.com/robfig/cron"
)

// StartAlertsTransferCron will start the main cront func for transfer alerts
// from alertmanager to elasticsearch
func StartAlertsTransferCron(config *Config) {
	log.Printf("Start fetch alerts from alertmanager server [%s]", config.AlertManager.URL)
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %s", config.Interval), func() {
		alerts, err := fetchAlerts(config)
		if err != nil {
			config.fetchAlertsErrorsTotal.Inc()
			log.Printf("Error: %s", err.Error())
			return
		}
		config.fetchAlertsTimesTotal.Inc()
		// appendAlerts will save current alerts to
		err = appendAlerts(config, alerts)
		if err != nil {
			config.bulkInsertErrorsTotal.Inc()
			log.Printf("Error: %s", err.Error())
			return
		}
		config.bulkInsertsTimesTotal.Inc()
		// TODO save alerts to elasticsearch
	})
	c.Start()
}

type jsonResponse struct {
	Status string        `json:"status"`
	Alerts []types.Alert `json:"data"`
}

// startFetchAlerts requests get alerts from alertmanager
// and save it to elasticsearch
func fetchAlerts(config *Config) (alerts []types.Alert, err error) {
	// TODO fetch api to get alert
	res := jsonResponse{}
	api := fmt.Sprintf("%s%s", config.AlertManager.URL, config.AlertManager.AlertsAPI)
	resp, err := http.Get(api)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return
	}
	return res.Alerts, nil
}

func appendAlerts(config *Config, alerts []types.Alert) error {
	if config.ESClient == nil {
		return errors.New("Client to connect elasticsearch not ready")
	}
	if len(alerts) == 0 {
		return nil
	}
	bulk := config.ESClient.Bulk().Index(config.EleasticSearch.IndexName).Type(config.EleasticSearch.TypeName)

	for _, alert := range alerts {
		bulk.Add(elastic.NewBulkIndexRequest().Id(strconv.FormatUint(uint64(alert.Fingerprint()), 10)).Doc(alert))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := bulk.Do(ctx); err != nil {
		return fmt.Errorf("Save to elasticsearch failed : %s", err.Error())
	}
	return nil
}

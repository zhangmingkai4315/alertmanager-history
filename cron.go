package main

import (
	"fmt"
	"log"

	"github.com/olivere/elastic"
	"github.com/prometheus/alertmanager/types"
	"github.com/robfig/cron"
)

// StartAlertsTransferCron will start the main cront func for transfer alerts
// from alertmanager to elasticsearch
func StartAlertsTransferCron(config *Config, client *elastic.Client) {
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %s", config.Interval), func() {
		alerts, err := fetchAlerts(client)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		if len(alerts) == 0 {
			return
		}
		// TODO save alerts to elasticsearch
	})
	c.Start()
}

// startFetchAlerts requests get alerts from alertmanager
// and save it to elasticsearch
func fetchAlerts(client *elastic.Client) (alerts []types.Alert, err error) {
	// TODO fetch api to get alerts
	return
}

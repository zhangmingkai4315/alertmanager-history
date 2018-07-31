package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configfile             string
	showVersion            bool
	app                    = "alertmanager_history"
	version                = "v0.1"
	versionInfo            = fmt.Sprintf("%s %s (%s)", app, version, runtime.Version())
	fetchAlertsErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: app,
		Name:      "fetch_alerts_errors_total",
		Help:      "Total number of errors when fetching using cron job",
	})
	fetchAlertsTimesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: app,
		Name:      "fetch_alerts_times_total",
		Help:      "Total fetch times since start cron job",
	})
	bulkInsertErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: app,
		Name:      "bulk_insert_errors_total",
		Help:      "Total number of errors when bulk insert elasticsearch",
	})
	bulkInsertTimesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: app,
		Name:      "bulk_insert_times_total",
		Help:      "Total bulk insert times since start cron job",
	})
)

func init() {
	flag.StringVar(&configfile, "c", "config.yml", "Config file path")
	flag.BoolVar(&showVersion, "v", false, "Print version of application")
	prometheus.MustRegister(fetchAlertsErrorsTotal)
	prometheus.MustRegister(fetchAlertsTimesTotal)
	prometheus.MustRegister(bulkInsertErrorsTotal)
	prometheus.MustRegister(bulkInsertTimesTotal)
}

func main() {
	flag.Parse()
	if showVersion {
		fmt.Println(versionInfo)
		os.Exit(0)
	}

	config, err := NewConfig(configfile)
	if err != nil {
		panic(err)
	}

	_, err = NewESClient(config)
	if err != nil {
		panic(err)
	}

	config.fetchAlertsErrorsTotal = fetchAlertsErrorsTotal
	config.fetchAlertsTimesTotal = fetchAlertsTimesTotal
	config.bulkInsertErrorsTotal = bulkInsertErrorsTotal
	config.bulkInsertsTimesTotal = bulkInsertTimesTotal
	go StartAlertsTransferCron(config)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, versionInfo)
	})
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/search", searchHandlerWithConfig(config))
	serve := &http.Server{
		Addr:         config.ListenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(serve.ListenAndServe())
}

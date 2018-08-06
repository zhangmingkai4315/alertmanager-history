package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

var (
	configfile             string
	showVersion            bool
	app                    = "alertmanager_history"
	version                = "v0.1.0"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, versionInfo)
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/search", searchHandlerWithConfig(config))

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(config.ListenAddr, handler)
}

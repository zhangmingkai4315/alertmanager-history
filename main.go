package main

import (
	"flag"
	"net/http"
)

var configfile string

func init() {
	flag.StringVar(&configfile, "config", "config.yml", "config file path")
}
func main() {
	config, err := NewConfig(configfile)
	if err != nil {
		panic(err)
	}

	esclient, err := NewESClient(config)
	if err != nil {
		panic(err)
	}
	go StartAlertsTransferCron(config, esclient)
	http.ListenAndServe(":8080", newRouter())
}

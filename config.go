package main

import (
	"io/ioutil"
	"log"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/olivere/elastic"
	"gopkg.in/yaml.v2"
)

var globalConfig *Config

// Config hold all config items in one struct
type Config struct {
	//EleasticSearch define the connection url for elasticsearch
	EleasticSearch struct {
		URL       string `yaml:"url"`
		IndexName string `yaml:"index_name"`
		TypeName  string `yaml:"type_name"`
	} `yaml:"elasticsearch"`
	ESClient *elastic.Client
	//AlertManager define the connection url for alertmanager
	AlertManager struct {
		URL       string `yaml:"url"`
		AlertsAPI string `yaml:"alerts_api"`
	} `yaml:"alertmanager"`
	Interval   string `yaml:"interval"`
	ListenAddr string `yaml:"listenAddr"`

	fetchAlertsErrorsTotal prometheus.Counter
	fetchAlertsTimesTotal  prometheus.Counter
	bulkInsertErrorsTotal  prometheus.Counter
	bulkInsertsTimesTotal  prometheus.Counter
}

// NewConfig create a new config struct and
// read file to load infomation
func NewConfig(file string) (*Config, error) {
	log.Printf("Read config file from %s\n", file)
	config := &Config{}
	err := config.readConfigFile(file)
	if err != nil {
		return nil, err
	}
	log.Println("Load config file success")
	return config, nil
}

// GetGlobalConfig get global config object
func GetGlobalConfig() *Config {
	return globalConfig
}

func (c *Config) readConfigFile(file string) error {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal Error : %v", err)
	}
	return nil
}

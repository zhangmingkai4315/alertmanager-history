package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config hold all config items in one struct
type Config struct {
	ESHTTPS bool   `yaml:"es_https"`
	ESHost  string `yaml:"es_host"`
	ESPort  int    `yaml:"es_port"`

	AMHost string `yaml:"am_host"`
	AMPort int    `yaml:"am_port"`

	Interval string `yaml:"interval"`
}

// NewConfig create a new config struct and
// read file to load infomation
func NewConfig(file string) (*Config, error) {
	config := &Config{}
	return config.readConfigFile(file)
}

func (c *Config) readConfigFile(file string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal Error : %v", err)
	}
	return c, nil
}

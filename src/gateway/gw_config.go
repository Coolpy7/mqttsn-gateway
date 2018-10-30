package mqttsn

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type GatewayConfig struct {
	IsAggregate      bool   `yaml:"IsAggregate"`
	Host             string `yaml:"Host"`
	Port             int    `yaml:"Port"`
	BrokerHost       string `yaml:"BrokerHost"`
	BrokerPort       int    `yaml:"BrokerPort"`
	BrokerUser       string `yaml:"BrokerUser"`
	BrokerPassword   string `yaml:"BrokerPassword"`
	LogFilePath      string `yaml:"LogFilePath"`
	MessageQueueSize int    `yaml:"MessageQueueSize"`
	ReadBuffSize     int    `yaml:"ReadBuffSize"`
	WriteBuffSize    int    `yaml:"WriteBuffSize"`
}

func ParseConfig(filename string) (*GatewayConfig, error) {
	var config GatewayConfig
	ymlStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(ymlStr, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

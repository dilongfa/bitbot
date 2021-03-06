package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Hitbtc         Credential
	Poloniex       Credential
	Kraken         Credential
	TheRockTrading Credential `json:"The Rock Trading"`
}

type Credential struct {
	Key       string
	Secret    string
	Addresses map[string]string
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	dest := &Config{}
	err = json.Unmarshal(data, dest)
	return dest, err
}

package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Config struct {
	AdminUser   int    `json:"admin"`
	Token       string `json:"token"`
	Channel     string `json:"channel"`
	FixedFooter string `json:"fixed_footer"`
	Watermark   struct {
		Watermark string `json:"watermark"`
		Height    int    `json:"height"`
		Logo      string `json:"logo"`
	} `json:"watermark"`
}

var once sync.Once
var config Config

func Read() Config {
	once.Do(func() {
		bytes, err := ioutil.ReadFile("config.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
	})

	return config
}

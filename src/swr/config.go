package swr

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Name string
	Addr string
}

var _config *Configuration

func Config() *Configuration {
	if _config == nil {
		fp, err := ioutil.ReadFile("./data/sys/config.json")
		ErrorCheck(err)
		err = json.Unmarshal(fp, &_config)
		ErrorCheck(err)
		log.Printf("Configuration loaded.")
	}
	return _config
}

package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Port   string
	Secret []byte
)

type Configs struct {
	Port   string `json:"port"`
	Secret string `json:"secret"`
}

func init() {
	var config Configs
	jsonFile, errOpenFile := os.Open("config.json")
	if errOpenFile != nil {
		panic("error opening config file: " + errOpenFile.Error())
	}
	byteValueConfig, errReadJson := ioutil.ReadAll(jsonFile)
	if errReadJson != nil {
		log.Println("error reading config file: " + errReadJson.Error())
	}

	errParseJson := json.Unmarshal(byteValueConfig, &config)

	if errParseJson != nil {
		log.Println("error parsing config file: " + errParseJson.Error())
	}

	Secret = []byte(config.Secret)
	Port = config.Port
}

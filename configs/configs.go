package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	// Port No from config.json file
	Port string
	// Secret JWT Key from config.json file
	Secret []byte
)

// Configs hold config variables
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
	byteValueConfig, errReadJSON := ioutil.ReadAll(jsonFile)
	if errReadJSON != nil {
		log.Println("error reading config file: " + errReadJSON.Error())
	}

	errParseJSON := json.Unmarshal(byteValueConfig, &config)

	if errParseJSON != nil {
		log.Println("error parsing config file: " + errParseJSON.Error())
	}

	Secret = []byte(config.Secret)
	Port = config.Port
}

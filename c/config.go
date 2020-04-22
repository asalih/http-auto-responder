package config

import (
	"encoding/json"
	"io/ioutil"
)

//Config Application settings
type Config struct {
	DatabaseName    string `json:"databaseName"`
	JSONsFolderPath string `json:"jsonsFolderPath"`
	Port            int    `json:"port"`
}

//Configuration ...
var Configuration Config

//InitConfig ...
func InitConfig() {
	InitConfigFile("config.json")
}

//InitConfigFile initializes the config file
func InitConfigFile(cnfFile string) {

	jsonFile, err := ioutil.ReadFile(cnfFile)

	jerr := json.Unmarshal(jsonFile, &Configuration)

	if err != nil || jerr != nil {
		panic(err)
	}
}

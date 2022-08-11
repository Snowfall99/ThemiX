package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Id      uint64 `json:"id"`
	Batch   int    `json:"batchsize"`
	Port    int    `json:"port"`
	Address string `json:"adress"`
	Key     string `json:"key_path"`
	Cluster string `json:"cluster"`
	Pk      string `json:"pk"`
	Ck      string `json:"ck"`
}

func ReadConfig(path string) (Configuration, error) {
	var config Configuration
	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("readConfig: %v", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("readConfig: %v", err)
	}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		return config, fmt.Errorf("readConfig: %v", err)
	}
	return config, nil
}

package config

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"os"
)

var MainConf Config

type DbData struct {
	Name string
	Uri  string
}
type Config struct {
	Local   string
	Remotes []DbData
}

type data struct {
	Config Config
}

func LoadConfig() {
	file, err := os.Open("conf.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var yamlData string
	for scanner.Scan() {
		yamlData += scanner.Text() + "\n"
	}

	if scanner.Err() != nil {
		panic(err)
	}

	var cfg data
	err = yaml.Unmarshal([]byte(yamlData), &cfg)
	if err != nil {
		panic(err)
	}

	MainConf = cfg.Config
}

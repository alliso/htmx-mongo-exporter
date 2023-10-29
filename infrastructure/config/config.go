package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type DbData struct {
	Name string
	Uri  string
}
type Config struct {
	Local   DbData
	Remotes []DbData
}

type data struct {
	config Config
}

type Loader interface {
	loadConfig() error
}

func (c Config) loadConfig() (Config, error) {
	buf, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		return c, err
	}

	var data = &data{}

	err = yaml.Unmarshal(buf, data)

	return data.config, err
}

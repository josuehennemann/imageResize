package main

import (
	"github.com/josuehennemann/yaml"
	"io/ioutil"
)

type Config struct {
	LogPath     string `yaml:"logPath"`     // path do diretorio de log
	DataPath    string `yaml:"dataPath"`    // path do diretorio data da aplicação
	HttpAddress string `yaml:"httpAddress"` //endereço http
	LibPath     string `yaml:"libPath"`
	IsDev       bool
}

func initConfig() error {
	config = &Config{}
	data, err := ioutil.ReadFile(*fileConf)
	yaml.Unmarshal(data, config)
	if config.HttpAddress == "localhost" || config.HttpAddress == ":8080" {
		config.IsDev = true
	}
	return err
}

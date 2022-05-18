package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
}

var GlobalConfig *Config

func ConfigInit() {
	content, err := ioutil.ReadFile("./conf/config.yaml")
	if err != nil {
		log.Fatalf("read token from file failed, err: %v", err)
	}
	if err = yaml.Unmarshal(content, &GlobalConfig); err != nil {
		log.Fatalf("parse config failed, err: %v", err)
	}
}

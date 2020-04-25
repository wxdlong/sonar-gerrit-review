package main

import (
	"github.com/wxdlong/sonar-gerrit-review/gerrit"
	"github.com/wxdlong/sonar-gerrit-review/sonar"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Gerrit gerrit.Gerrit `yaml:"gerrit"`
	Sonar  sonar.Sonar
}

//ReadFromFile read config file "sonar-gerrit.yaml"
//convert to Config
func ReadFromFile(file string) {
	yam, err := ioutil.ReadFile("sonar-gerrit.yaml")
	if err != nil {
		log.Fatal("Read sonar-gerrit.yaml failed ", err)
	}
	config := &Config{}
	yaml.Unmarshal(yam, config)
	log.Println(config)
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

type (
	Config struct {
		Watches []Watch `json:"watches"`
	}
	Watch struct {
		ProjectRx string       `json:"tag"`
		SlackUrl  string       `json:"slackUrl"`
		SlackMsg  SlackMessage `json:"slackMsg"`
	}
)

func LoadConfig(path string) Config {
	cfg := Config{}
	if file, err := ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(file, &cfg)
		log.Printf("Watching %d Project filters, %s", len(cfg.Watches), cfg.Watches[0].ProjectRx)
	} else {
		log.Fatal("Unable to load config '%v'\n", err)
	}

	return cfg
}

func (c Config) Process(p Project) (url string, msg SlackMessage) {
	log.Printf("process::%s\n", p.Name)
	for _, watch := range c.Watches {
		if match, _ := regexp.MatchString(watch.ProjectRx, p.Name); match {
			// TODO: optimize? add "^" + name + "$" to map of projects with slack msg pointers
			return watch.SlackUrl, watch.SlackMsg
		}
	}
	return "", SlackMessage{}
}

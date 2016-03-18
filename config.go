package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

type (
	Config struct {
		Remotes []string `json:"remotes"`
		Watches []Watch  `json:"watches"`
	}
	Watch struct {
		ProjectRx    []string          `json:"tags"`
		SlackUrl     string            `json:"slackUrl"`
		Transitions  []string          `json:"transitions"`
		ColorMapping map[string]string `json:"colormapping"`
		SlackMsg     SlackMessage      `json:"slackMsg"`
	}
)

func ConfigChanged(path string) bool {
	return false
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	file, err := ioutil.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(file, &cfg)
		log.Printf("Watching %d Project filters, %s", len(cfg.Watches), cfg.Watches[0].ProjectRx)
	} else {
		log.Printf("Unable to load config '%v'\n", err)
	}
	return cfg, err
}

func InSlice(check string, slice []string) bool {
	for _, entry := range slice {
		if entry == check {
			return true
		}
	}
	return false
}

func (c Config) Process(p Project) (url string, msg SlackMessage) {
	log.Printf("process::%s\n", p.Name)
	for _, watch := range c.Watches {
		for _, projectRx := range watch.ProjectRx {
			if match, _ := regexp.MatchString(projectRx, p.Name); match {
				// TODO: optimize? add "^" + name + "$" to map of projects with slack msg pointers
				log.Printf("Lookig for %s in %q\n", p.Transition, watch.Transitions)

				if InSlice(p.Transition, watch.Transitions) {
					return watch.SlackUrl, watch.SlackMsg
				} else {
					return "", SlackMessage{}
				}
			}
		}
	}
	return "", SlackMessage{}
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

type (

	// Config is the entire configuration.
	// Remotes:  which remote ccTrays to mointor
	// Watches: see struct watches
	// SlackMessages: see struct SlackMessage
	Config struct {
		Remotes       []string                `json:"remotes"`
		Watches       []Watch                 `json:"watches"`
		SlackMessages map[string]SlackMessage `json:"slackmessages"`
	}

	// Watch is the mapping between ccTray Project name and slack.
	// SlckUrl is the web hook to slack to use adn channel the chanel to post messages to.
	// Transitions define whcih transitions to report
	Watch struct {
		ProjectRx   []string `json:"tags"`
		SlackURL    string   `json:"slackUrl"`
		Transitions []string `json:"transitions"`
		Channel     string   `json:"channel"`
	}
)

// ConfigChanged returns true
func ConfigChanged(path string) bool {
	return true
}

// LoadConfig reads the config from path given as argument
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

func inSlice(check string, slice []string) bool {
	for _, entry := range slice {
		if entry == check {
			return true
		}
	}
	return false
}

// Process returns a url and template slackmessage to be used for sending messages to slack given a certain Project
func (c Config) Process(p Project) (url string, msg SlackMessage) {
	log.Printf("process::%s\n", p.Name)

	for _, watch := range c.Watches {
		for _, projectRx := range watch.ProjectRx {
			if match, _ := regexp.MatchString(projectRx, p.Name); match {
				// TODO: optimize? add "^" + name + "$" to map of projects with slack msg pointers
				log.Printf("Lookig for %s in %q\n", p.Transition, watch.Transitions)

				if inSlice(p.Transition, watch.Transitions) {
					message := c.SlackMessages[p.Transition]
					message.Channel = watch.Channel
					return watch.SlackURL, message
				}
				return "", SlackMessage{}
			}
		}
	}
	return "", SlackMessage{}
}

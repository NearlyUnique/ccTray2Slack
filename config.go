package main

import (
	"encoding/json"
	"fmt"
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
		Identifier  string   `json:"identifier"`
		ProjectRx   []string `json:"tags"`
		SlackURL    string   `json:"slackUrl"`
		Transitions []string `json:"transitions"`
		Channel     string   `json:"channel"`
	}

	DefaultConfigArgs struct {
		RemoteURL string
		SlackHook string
	}
)

var (
	defaultRemotes = []string{"htttp://yourremotecctray:8153/go"}
	defaultWatches = []Watch{
		Watch{"Identifier",
			[]string{"Project1", "Project2"},
			"slackURL",
			[]string{"Fixed", "Broken"},
			"#slack_channel",
		},
		Watch{"Identifier2",
			[]string{"Project3", "Project4"},
			"slackURL",
			[]string{"Success", "Failure"},
			"#slack_channel",
		},
	}
)

func PrintDefaultConfig(args DefaultConfigArgs) {
	c := Config{}
	c.Watches = defaultWatches
	for i := range c.Watches {
		c.Watches[i].SlackURL = args.SlackHook
	}
	c.Remotes = defaultRemotes
	c.Remotes[0] = args.RemoteURL
	c.SlackMessages = make(map[string]SlackMessage)
	c.SlackMessages = defaultSlackMessages
	b, _ := json.MarshalIndent(c, " ", " ")
	fmt.Printf("%v", string(b))
}

// ConfigChanged returns true
func ConfigChanged(path string) bool {
	return true
}

// LoadConfig reads the config from path given as argument
func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	cfg.SlackMessages = make(map[string]SlackMessage)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return cfg, err
	}
	for _, file := range files {
		fileData, err := ioutil.ReadFile(path + file.Name())
		if err == nil {
			cfgTmp := Config{}
			err = json.Unmarshal(fileData, &cfgTmp)
			cfg.Add(cfgTmp)
		} else {
			log.Printf("Unable to load config '%v'\n", err)
		}
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

func (c *Config) Add(cfg Config) {
	for _, remote := range cfg.Remotes {
		c.Remotes = append(c.Remotes, remote)
	}
	for _, watch := range cfg.Watches {
		c.Watches = append(c.Watches, watch)
	}

	for key, slackMessage := range cfg.SlackMessages {
		c.SlackMessages[key] = slackMessage
	}

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

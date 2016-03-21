package main

import (
	"flag"
	"log"
	"time"
)

type CommandLineArgs struct {
	password   string
	username   string
	configPath string
}

var commandLineArgs CommandLineArgs

func main() {
	configPath, usr, pwd := parseCmdLine()
	commandLineArgs = CommandLineArgs{configPath: configPath, username: usr, password: pwd}
	var cc ccTray
	if config, err := LoadConfig(configPath); err == nil {
		cc = CreateCcTray(config.Remotes[0])
		cc.Username = usr
		cc.Password = pwd
		RunPollLoop(config, cc)
	} else {
		log.Fatal("Unable to load config stoping executions")
	}
}

func RunPollLoop(config Config, cc ccTray) {
	ticker := time.NewTicker(10 * time.Second)
	go cc.GetLatest()

	for {
		select {
		case p := <-cc.Ch:
			if url, msg := config.Process(p); url != "" {
				log.Printf("posting for %q\n", p.Name)
				msg.UpdateMessage(p)
				msg.PostSlackMessage(url)
			} else {
				log.Printf("skipped %q\n", p.Name)
			}
		case e := <-cc.ChErr:
			if e != nil {
				log.Fatalf("Failed to get ccTray Data \n%v\n", e)
			}
			log.Println("Cycle complete")
		case <-ticker.C:
			if ConfigChanged(commandLineArgs.configPath) {
				if temp, err := LoadConfig(commandLineArgs.configPath); err == nil {
					config = temp
				} else {
					log.Printf("Unable to read config file %s.", commandLineArgs.configPath)
				}
			}
			log.Println("checking ...")
			go cc.GetLatest()
		}
	}
}

func parseCmdLine() (configPath, user, password string) {
	path := flag.String("config", "config.json", "config for project filter and Slack integration")
	usr := flag.String("username", "", "cctray server account name")
	pwd := flag.String("password", "", "cctray server account password")

	flag.Parse()

	if *path == "" {
		log.Fatal("config path must be set")
	}

	return *path, *usr, *pwd
}

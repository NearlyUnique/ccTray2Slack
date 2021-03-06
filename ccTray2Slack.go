package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	ccTrayUrl, configPath, usr, pwd := parseCmdLine()
	config := LoadConfig(configPath)
	cc := CreateCcTray(ccTrayUrl)
	cc.Username = usr
	cc.Password = pwd
	RunPollLoop(config, cc)
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
			log.Println("checking ...")
			go cc.GetLatest()
		}
	}
}

func parseCmdLine() (ccTray, configPath, user, password string) {
	url := flag.String("url", "http://localhost/cctray.xml", "url for the CCTray xml data")
	path := flag.String("config", "watch.json", "config for project filter and Slack integration")
	usr := flag.String("username", "", "cctray server account name")
	pwd := flag.String("password", "", "cctray server account password")

	flag.Parse()

	if *url == "" {
		log.Fatal("url must be set")
	}
	if *path == "" {
		log.Fatal("config path must be set")
	}

	log.Printf("Downloading from '%s'", *url)
	return *url, *path, *usr, *pwd
}

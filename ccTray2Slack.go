package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	SlackMessage struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		IconUrl  string `json:"icon_url"`
		Channel  string `json:"channel"`
	}
)

func main() {
	ccTrayUrl := flag.String("url", "https://bigvisiblewall.googlecode.com/hg/web/cctray.xml", "url for the CCTray xml data")
	slackUrl := flag.String("slack", "https://hooks.slack.com/services/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", "url for Slack integration")
	flag.Parse()

	if *ccTrayUrl == "" {
		log.Fatal("url must be set")
	}
	log.Printf("Downloading from '%s'", *ccTrayUrl)

	cc := CreateCcTray(*ccTrayUrl)

	go cc.GetLatest()

	select {
	case p := <-cc.Ch:
		msg := fmt.Sprintf(
			"%q was %q at %v, <%s|Visit Page>",
			p.Name, p.LastBuildStatus, p.LastBuildTime, p.WebUrl)
		PostSlackMessage(*slackUrl, msg)
	case e := <-cc.ChErr:
		log.Fatalf("Failed to get ccTray Data \n%v\n", e)
	}
}

func PostSlackMessage(url, msg string) {
	slack := SlackMessage{msg, "GoBotGo", "http://fc01.deviantart.net/fs70/f/2009/348/2/0/Norton_Ghost_by_jakaneko.png", "#random"}
	jsonStr, _ := json.Marshal(slack)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

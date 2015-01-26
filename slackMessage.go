package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type (
	SlackMessage struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		IconUrl  string `json:"icon_url"`
		Channel  string `json:"channel"`
	}
)

func (s *SlackMessage) UpdateMessage(p Project) {
	rx := regexp.MustCompile("%.*?%")
	s.Text = rx.ReplaceAllStringFunc(s.Text, func(src string) string {
		switch src {
		case "%project%":
			return p.Name
		case "%status%":
			return p.LastBuildStatus
		case "%label%":
			return p.LastBuildLabel
		case "%url%":
			return p.WebUrl
		case "%time%":
			return p.LastBuildTime.Format("2006-01-02 15:04:05")
		}
		return src
	})
}

func (s *SlackMessage) PostSlackMessage(url string) error {
	if url == "debug" {
		log.Printf("%q\n%v\n", url, &s)
	} else {
		jsonStr, _ := json.Marshal(&s)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		if resp, err := client.Do(req); err == nil {
			defer resp.Body.Close()

			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				log.Printf("Err:%v\nBody: %s\n", err, body)
				return err
			}

		} else {
			return err
		}
	}
	return nil
}

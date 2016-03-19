package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Attachement struct {
	Title string `json:"title"`
	Color string `json:"color"`
	Text  string `json:"text"`
}

type (
	SlackMessage struct {
		Text         string        `json:"text"`
		Attachements []Attachement `json:"attachments"`
		Username     string        `json:"username"`
		IconUrl      string        `json:"icon_url"`
		Channel      string        `json:"channel"`
	}
)

var (
	rx = regexp.MustCompile("%.*?%")
)

func replaceString(s string, p Project) string {
	return rx.ReplaceAllStringFunc(s, func(src string) string {
		switch src {
		case "%project%":
			return p.Name
		case "%status%":
			return p.Transition
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
func (s *SlackMessage) UpdateMessage(p Project) {
	s.Text = replaceString(s.Text, p)
	for i, _ := range s.Attachements {
		s.Attachements[i].Text = replaceString(s.Attachements[i].Text, p)
	}
}

func (s *SlackMessage) PostSlackMessage(url string) error {
	if url == "debug" {
		log.Printf("HTTP POST -> Slack\n%v\n", *s)
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

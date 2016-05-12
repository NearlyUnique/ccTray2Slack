package slackmessage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/christer79/ccTray2Slack/cctray"
)

//Attachement defines an attchemnt included in a slack message
type Attachement struct {
	Title string `json:"title"`
	Color string `json:"color"`
	Text  string `json:"text"`
}

//SlackMessage is the payload sent to slack in the message request
type (
	SlackMessage struct {
		Text         string        `json:"text"`
		Attachements []Attachement `json:"attachments"`
		Username     string        `json:"username"`
		IonEmoji     string        `json:"icon_emoji"`
		Channel      string        `json:"channel"`
	}
)

var (
	rx                   = regexp.MustCompile("%.*?%")
	DefaultSlackMessages = map[string]SlackMessage{
		"Success": SlackMessage{
			"New update from your build-server",
			[]Attachement{Attachement{"<%url%|%Project%-%label%>", "#e23e3e", "was %status% at %time%"}},
			"GO CD Integration",
			":building_construction:",
			"",
		},
		"Failure": SlackMessage{
			"New update from your build-server",
			[]Attachement{Attachement{"<%url%|%Project%-%label%>", "#7fbb43", "was %status% at %time%"}},
			"GO CD Integration",
			":building_construction:",
			"",
		},
		"Broken": SlackMessage{
			"New update from your build-server",
			[]Attachement{Attachement{"<%url%|%Project%-%label%>", "#7fbb43", "was %status% at %time%"}},
			"GO CD Integration",
			":building_construction:",
			"",
		},
		"Fixed": SlackMessage{
			"New update from your build-server",
			[]Attachement{Attachement{"<%url%|%Project%-%label%>", "#e23e3e", "was %status% at %time%"}},
			"GO CD Integration",
			":building_construction:",
			"",
		},
	}
)

func replaceString(s string, p cctray.Project) string {
	return rx.ReplaceAllStringFunc(s, func(src string) string {
		switch src {
		case "%project%":
			return p.Name
		case "%status%":
			return p.Transition
		case "%label%":
			return p.LastBuildLabel
		case "%url%":
			return p.WebURL
		case "%time%":
			return p.LastBuildTime.Format("2006-01-02 15:04:05")
		}
		return src
	})
}

// UpdateMessage replaces keywords in a slack message with the matching values from a Project.
func (s *SlackMessage) UpdateMessage(p cctray.Project) {
	s.Text = replaceString(s.Text, p)
	for i := range s.Attachements {
		s.Attachements[i].Text = replaceString(s.Attachements[i].Text, p)
		s.Attachements[i].Title = replaceString(s.Attachements[i].Title, p)
	}
}

//PostSlackMessage posts a message to slack on the url passed as argument
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

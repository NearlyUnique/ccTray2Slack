package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	ccTray struct {
		URL         string
		Username    string
		Password    string
		Ch          chan Project
		ChErr       chan error
		interesting []string
		previous    map[string]Project
	}
)

// xml for ccTray schema
type (
	Project struct {
		Name            string   `xml:"name,attr"`
		Activity        string   `xml:"activity,attr"`
		LastBuildStatus string   `xml:"lastBuildStatus,attr"`
		LastBuildLabel  string   `xml:"lastBuildLabel,attr"`
		LastBuildTime   projTime `xml:"lastBuildTime,attr"`
		WebUrl          string   `xml:"webUrl,attr"`
		Transition      string
	}
	Projects struct {
		Projects []Project `xml:"Project"`
	}
)

func (p Project) String() string {
	return fmt.Sprintf(
		"name=%s, activity=%s, status=%s, label=%s, time=%v, url=%s, transition=%s",
		p.Name,
		p.Activity,
		p.LastBuildStatus,
		p.LastBuildLabel,
		p.LastBuildTime,
		p.WebUrl,
		p.Transition)
}

func CreateCcTray(url string) ccTray {
	return ccTray{
		URL:      url,
		Ch:       make(chan Project),
		ChErr:    make(chan error),
		previous: make(map[string]Project),
	}
}

func (cc ccTray) GetLatest() {
	var err error
	client := &http.Client{}
	req, err := http.NewRequest("GET", cc.URL, nil)
	req.SetBasicAuth(cc.Username, cc.Password)

	if resp, err := client.Do(req); err == nil {
		log.Printf("CC Tray http GET ok")
		defer resp.Body.Close()
		log.Printf("Code: %d, %q\n", resp.StatusCode, resp.Status)

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			p := Projects{}
			xml.Unmarshal(body, &p)
			cc.publishChanges(p.Projects)
			return
		}
	} else {
		log.Fatalf("CC Tray http GET failed %v\n", err)
	}
	cc.ChErr <- err
}

func (cc ccTray) publishChanges(projects []Project) {
	log.Printf("publishing %d\n", len(projects))
	for _, current := range projects {
		if prev, ok := cc.previous[current.Name]; ok {
			if prev != current {
				log.Printf("Replacing %q - \"%q\" \n", current.Name, current.Transition)
				cc.previous[current.Name] = current

				if current.Activity == "Sleeping" {
					current.Transition = current.LastBuildStatus
					cc.Ch <- current
				}
				if prev.LastBuildStatus != current.LastBuildStatus {
					if current.LastBuildStatus == "Success" {
						current.Transition = "Fixed"
					}
					if current.LastBuildStatus == "Failed" {
						current.Transition = "Broken"
					}
					if current.Activity == "Sleeping" {
						cc.Ch <- current
					}
				}

			} else {
				//log.Printf("No Change %q\n", current.Name)
			}
		} else {
			//log.Printf("Adding    %q\n", current.Name)
			cc.previous[current.Name] = current
		}
	}
	// everything is ok, finished looping - looks hacky to me but another channel, really?
	cc.ChErr <- nil
}

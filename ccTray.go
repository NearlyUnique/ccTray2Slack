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
		Url         string
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
		"name=%s, activity=%s, status=%s, label=%s, time=%v, url=%s",
		p.Name,
		p.Activity,
		p.LastBuildStatus,
		p.LastBuildLabel,
		p.LastBuildTime,
		p.WebUrl)
}

func CreateCcTray(url string) ccTray {
	return ccTray{
		Url:      url,
		Ch:       make(chan Project),
		ChErr:    make(chan error),
		previous: make(map[string]Project),
	}
}

func (cc ccTray) GetLatest() {
	var err error
	client := &http.Client{}
	req, err := http.NewRequest("GET", cc.Url, nil)
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
		if current.Activity == "Sleeping" {
			if prev, ok := cc.previous[current.Name]; ok {
				if prev != current {
					if prev.LastBuildStatus != current.LastBuildStatus {
						if current.LastBuildStatus == "Success" {
							current.Transition = "Fixed"
						}
						if current.LastBuildStatus == "Failure" {
							current.Transition = "Broken"
						}
					} else {
						current.Transition = current.LastBuildStatus
					}
					log.Printf("Replacing %q - \"%q\" \n", current.Name, current.Transition)
					cc.previous[current.Name] = current
					cc.Ch <- current
					current.Transition = ""
					cc.previous[current.Name] = current
				} else {
					//log.Printf("No Change %q\n", current.Name)
				}
			} else {
				//log.Printf("Adding    %q\n", current.Name)
				cc.previous[current.Name] = current
			}
		}
	}
	// everything is ok, finished looping - looks hacky to me but another channel, really?
	cc.ChErr <- nil
}

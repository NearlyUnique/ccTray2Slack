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
	if resp, err := http.Get(cc.Url); err == nil {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			p := Projects{}
			xml.Unmarshal(body, &p)
			cc.publishChanges(p.Projects)
			return
		}
	}
	cc.ChErr <- err
}

func (cc ccTray) publishChanges(projects []Project) {
	for _, current := range projects {
		if prev, ok := cc.previous[current.Name]; ok {
			if prev != current {
				log.Printf("Replacing %q\n", current.Name)
				cc.previous[current.Name] = current
				cc.Ch <- current
			} else {
				log.Printf("No Change %q\n", current.Name)
			}
		} else {
			log.Printf("Adding    %q\n", current.Name)
			cc.previous[current.Name] = current
		}
	}
	// everything is ok, finished looping - looks hacky to me but another channel, really?
	cc.ChErr <- nil
}

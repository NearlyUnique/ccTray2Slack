package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type (
	ccTray struct {
		Url      string
		Ch       chan Project
		ChErr    chan error
		previous map[string]Project
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

func CreateCcTray(url string) ccTray {
	return ccTray{
		Url:      url,
		Ch:       make(chan Project),
		ChErr:    make(chan error),
		previous: make(map[string]Project),
	}
}

func (c ccTray) GetLatest() {
	var e error
	if resp, e := http.Get(c.Url); e == nil {
		defer resp.Body.Close()

		if body, e := ioutil.ReadAll(resp.Body); e == nil {
			p := Projects{}
			xml.Unmarshal(body, &p)
			c.publishChanges(p.Projects)
			return
		}
	}
	c.ChErr <- e
}

func (c ccTray) publishChanges(projects []Project) {
	for _, current := range projects {
		if prev, ok := c.previous[current.Name]; ok {
			if prev != current {
				c.previous[current.Name] = current
				c.Ch <- current
			}
		} else {
			c.previous[current.Name] = current
		}
	}
}

package cctray

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var out io.Writer = os.Stdout

type (
	CcTray struct {
		URL         string
		Username    string
		Password    string
		Ch          chan Project
		ChErr       chan error
		interesting []string
		previous    map[string]Project
	}
)

// xml for CcTray schema
type (
	Project struct {
		Name            string   `xml:"name,attr"`
		Activity        string   `xml:"activity,attr"`
		LastBuildStatus string   `xml:"lastBuildStatus,attr"`
		LastBuildLabel  string   `xml:"lastBuildLabel,attr"`
		LastBuildTime   ProjTime `xml:"lastBuildTime,attr"`
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

func CreateCcTray(url string) CcTray {
	return CcTray{
		URL:      url,
		Ch:       make(chan Project),
		ChErr:    make(chan error),
		previous: make(map[string]Project),
	}
}

func (cc CcTray) GetProjects() (Projects, error) {
	var err error
	client := &http.Client{}
	req, err := http.NewRequest("GET", cc.URL, nil)
	req.SetBasicAuth(cc.Username, cc.Password)
	p := Projects{}

	if resp, err := client.Do(req); err == nil {
		log.Printf("CC Tray http GET ok")
		defer resp.Body.Close()
		log.Printf("Code: %d, %q\n", resp.StatusCode, resp.Status)

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			xml.Unmarshal(body, &p)
		}
	} else {
		log.Fatalf("CC Tray http GET failed %v\n", err)
	}
	return p, err
}

func (cc CcTray) ListProjects() {
	p, err := cc.GetProjects()
	if err == nil {
		for _, project := range p.Projects {
			fmt.Fprintf(out, "\"%v\",\n", project.Name)
		}
	} else {
		log.Fatalf("Error in ListProjects: %s\n", err)
	}
}

func (cc CcTray) GetLatest() {
	p, err := cc.GetProjects()
	if err == nil {
		cc.publishChanges(p.Projects)
	} else {
		cc.ChErr <- err
	}

}

func (cc CcTray) publishChanges(projects []Project) {
	log.Printf("publishing %d\n", len(projects))
	for _, current := range projects {
		if prev, ok := cc.previous[current.Name]; ok {
			if prev != current {
				log.Printf("Replacing %q - \"%q\" \n", current.Name, current.LastBuildStatus)
				cc.previous[current.Name] = current
				log.Printf("Status curr: %q prev: %q (%q)\n", current.LastBuildStatus, prev.LastBuildStatus, current.Activity)

				if current.Activity == "Sleeping" {
					current.Transition = current.LastBuildStatus
					cc.Ch <- current
				}

				if prev.LastBuildStatus != current.LastBuildStatus {
					if current.LastBuildStatus == "Success" {
						current.Transition = "Fixed"
					} else {
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

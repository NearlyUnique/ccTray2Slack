package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	xml1 = `<Projects>
<Project name="Project 1" activity="Sleeping" lastBuildStatus="Success" lastBuildLabel="1.2.3" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v1"/>
<Project name="Project 2" activity="Building" lastBuildStatus="Success" lastBuildLabel="1.2.6" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v2"/>
<Project name="Project 3" activity="Sleeping" lastBuildStatus="Failed" lastBuildLabel="1.2.7" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v3"/>
</Projects>`
	xml2 = `<Projects>
<Project name="Project 1" activity="Sleeping" lastBuildStatus="Passed" lastBuildLabel="1.2.3" lastBuildTime="2015-01-19T08:53:01" webUrl="http://localhost:8153/cruise/v1"/>
<Project name="Project 2" activity="Building" lastBuildStatus="Success" lastBuildLabel="1.2.6" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v2"/>
<Project name="Project 3" activity="Sleeping" lastBuildStatus="Failed" lastBuildLabel="1.2.7" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v3"/>
</Projects>`
	firstXml = true
)

func TestIt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if firstXml {
			fmt.Fprintln(w, xml1)
			firstXml = false
		} else {
			fmt.Fprintln(w, xml2)
		}
	}))
	defer ts.Close()

	sut := CreateCcTray(ts.URL)
	go func() {
		for {
			select {
			case p := <-sut.Ch:
				if p.Activity != "x" {
					t.Fail()
				}
			case e := <-sut.ChErr:
				if e != nil {
					t.Fail()
				}
			}
		}
	}()

	sut.GetLatest() // prime the system
	sut.GetLatest() // get changes
}

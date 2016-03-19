package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	ccXml = []string{
		`<Projects>
	<Project name="Project 1" activity="Sleeping" lastBuildStatus="Success" lastBuildLabel="1.2.3" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v1"/>
	<Project name="Project 2" activity="Building" lastBuildStatus="Success" lastBuildLabel="1.2.6" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v2"/>
	<Project name="Project 3" activity="Sleeping" lastBuildStatus="Failed" lastBuildLabel="1.2.7" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v3"/>
</Projects>`,
		`<Projects>
	<Project name="Project 1" activity="Sleeping" lastBuildStatus="Success" lastBuildLabel="1.2.3" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v1"/>
	<Project name="Project 2" activity="Building" lastBuildStatus="Success" lastBuildLabel="1.2.6" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v2"/>
	<Project name="Project 3" activity="Sleeping" lastBuildStatus="Failed" lastBuildLabel="1.2.7" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v3"/>
</Projects>`,
		`<Projects>
	<Project name="Project 1" activity="Sleeping" lastBuildStatus="Passed" lastBuildLabel="1.2.3" lastBuildTime="2015-01-19T08:53:01" webUrl="http://localhost:8153/cruise/v1"/>
	<Project name="Project 2" activity="Building" lastBuildStatus="Success" lastBuildLabel="1.2.6" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v2"/>
	<Project name="Project 3" activity="Sleeping" lastBuildStatus="Success" lastBuildLabel="1.2.7" lastBuildTime="2009-07-27T18:17:19" webUrl="http://localhost:8153/cruise/v3"/>
</Projects>`,
		`<Projects>
	<Project name="Project 1" activity="Sleeping" lastBuildStatus="Passed" lastBuildLabel="1.2.3" lastBuildTime="2015-01-19T08:53:01" webUrl="http://localhost:8153/cruise/v1"/>
	<Project name="Project 2" activity="Building" lastBuildStatus="Success" lastBuildLabel="1.2.6" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v2"/>
	<Project name="Project 3" activity="Sleeping" lastBuildStatus="Failed" lastBuildLabel="1.2.7" lastBuildTime="2009-07-27T14:17:19" webUrl="http://localhost:8153/cruise/v3"/>
</Projects>`,
	}

	expect = [6]string{"name=Project 3, activity=Sleeping, status=Success, label=1.2.7, time=2009-07-27 18:17:19 +0000 UTC, url=http://localhost:8153/cruise/v3, transition=Success",
		"name=Project 3, activity=Sleeping, status=Success, label=1.2.7, time=2009-07-27 18:17:19 +0000 UTC, url=http://localhost:8153/cruise/v3, transition=Fixed",
		"name=Project 1, activity=Sleeping, status=Success, label=1.2.3, time=2009-07-27 14:17:19 +0000 UTC, url=http://localhost:8153/cruise/v1, transition=Success",
		"name=Project 1, activity=Sleeping, status=Success, label=1.2.3, time=2009-07-27 14:17:19 +0000 UTC, url=http://localhost:8153/cruise/v1, transition=Fixed",
		"name=Project 3, activity=Sleeping, status=Failed, label=1.2.7, time=2009-07-27 14:17:19 +0000 UTC, url=http://localhost:8153/cruise/v3, transition=Failed",
		"name=Project 3, activity=Sleeping, status=Failed, label=1.2.7, time=2009-07-27 14:17:19 +0000 UTC, url=http://localhost:8153/cruise/v3, transition=Broken"}
)

func TestIt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var popped string
		w.Header().Set("Content-Type", "application/json")
		popped, ccXml = ccXml[len(ccXml)-1], ccXml[:len(ccXml)-1]
		fmt.Fprintln(w, popped)
	}))
	count := 0
	defer ts.Close()

	sut := CreateCcTray(ts.URL)

	go func() {
		for {
			select {
			case p := <-sut.Ch:
				equalString(t, fmt.Sprintf("%v", p), expect[count])
				count++
			case e := <-sut.ChErr:
				if e != nil {
					t.Error("Unexpected errors")
				}
			}
		}
	}()

	sut.GetLatest() // prime the system
	sut.GetLatest() // get changes
	sut.GetLatest() // get changes
}

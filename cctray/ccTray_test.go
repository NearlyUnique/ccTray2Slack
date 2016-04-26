package cctray

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"
)

var (
	ccXML = []string{
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

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CcTrayTestSuite struct{}

var _ = Suite(&CcTrayTestSuite{}) // Hook up gocheck into the "go test" runner.

func (s *CcTrayTestSuite) TestListProjects(c *C) {
	ccXMLCopy := ccXML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var popped string
		w.Header().Set("Content-Type", "application/json")
		popped, ccXMLCopy = ccXMLCopy[len(ccXMLCopy)-1], ccXMLCopy[:len(ccXMLCopy)-1]
		fmt.Fprintln(w, popped)
	}))

	defer ts.Close()
	sut := CreateCcTray(ts.URL)
	buf := &bytes.Buffer{}
	out = buf
	sut.ListProjects()
	expect := "\"Project 1\",\n\"Project 2\",\n\"Project 3\",\n"
	c.Assert(buf.String(), Equals, expect)
}

func (s *CcTrayTestSuite) TestIt(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var popped string
		w.Header().Set("Content-Type", "application/json")
		popped, ccXML = ccXML[len(ccXML)-1], ccXML[:len(ccXML)-1]
		fmt.Fprintln(w, popped)
	}))
	count := 0
	projectCount := 0
	defer ts.Close()

	sut := CreateCcTray(ts.URL)

	go func() {
		for {
			select {
			case p := <-sut.Ch:
				c.Assert(fmt.Sprintf("%v", p), Equals, expect[count])
				count++
			case <-sut.ChProjects:
				projectCount++
			case e := <-sut.ChErr:
				if e != nil {
					c.Error("Unexpected errors")
				}
			}
		}
	}()

	sut.GetLatest() // prime the system
	c.Assert(projectCount, Equals, 1)
	sut.GetLatest() // get changes
	c.Assert(count, Equals, 2)
	c.Assert(projectCount, Equals, 2)
	sut.GetLatest() // get changes
	c.Assert(count, Equals, 6)
	c.Assert(projectCount, Equals, 3)

}

package main

import (
	"testing"

	"github.com/christer79/ccTray2Slack/cctray"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ConfigTestSuite struct{}

var _ = Suite(&ConfigTestSuite{}) // Hook up gocheck into the "go test" runner.

var (
	groupProjectsConfig = `
	{
	"watches" : [
		{
			"tags": ["Project1", "Project2"],
			"transitions": ["Success","Failed"],
			"slackUrl": "project1",
			"channel": "#api_test",
			"identifier": "IdentifierA"
		},{
			"tags": ["Notinconfig"],
			"transitions": ["Success","Failed"],
			"slackUrl": "project1",
			"channel": "#api_test",
			"identifier": "IdentifierB"
		}

	]
}`
	testProjects = []cctray.Project{
		cctray.Project{Name: "Project1", Transition: "Fixed"},
		cctray.Project{Name: "Project2", Transition: "Success"},
		cctray.Project{Name: "Notinconfig", Transition: "Failed"}}
	expectedWatches = []Watch{Watch{"Identifier 1", []string{"^Openstack.*"}, "_", []string{"Success", "Failed"}, "#api_test"},
		Watch{"Identifier 2", []string{"^Provision.*"}, "_", []string{"Success", "Failed"}, "#api_test"}}

	groupProjectsProjects = cctray.Projects{testProjects}

	groupProjectsResult = map[string][]cctray.Project{
		"IdentifierB": []cctray.Project{testProjects[2]},
		"IdentifierA": []cctray.Project{testProjects[0], testProjects[1]},
	}
)

func (s *ConfigTestSuite) TestProcess(c *C) {
	config, _ := LoadConfig("testdata/config2.d")
	// When Processing a correct project but with wrong Transiton
	url, msg := config.Process(testProjects[0])
	// ... return empty url
	c.Assert(url, Equals, "")
	// When processing a correct project with corret transition
	url, msg = config.Process(testProjects[1])
	//... return message for the correct transition
	c.Assert(msg.Text, Equals, "Success text")
	// ... a non empty url
	c.Assert(url, Equals, "project2_3")
	// ... correct channel should be set
	c.Assert(msg.Channel, Equals, "#api_test")
	// When Processing a project which does not match a watched project
	url, msg = config.Process(testProjects[2])
	// ... return an empty url
	c.Assert(url, Equals, "")
}

func (s *ConfigTestSuite) TestLoadConfig(c *C) {

	config, _ := LoadConfig("testdata/config1.d/")
	c.Assert(len(config.Watches), Equals, 3)
	c.Assert(len(config.Remotes), Equals, 2)

	//TODO: Compare structs
	c.Assert(config.Watches[0].Channel, Equals, expectedWatches[0].Channel)
	c.Assert(config.Watches[1].ProjectRx[0], Equals, expectedWatches[1].ProjectRx[0])

	configPath := "testdata/config_which_does_not_exist.json"
	_, err := LoadConfig(configPath)
	c.Assert(err, NotNil)

}

func (s *ConfigTestSuite) TestGroupProjects(c *C) {
	config := Config{}
	generateConfig([]byte(groupProjectsConfig), &config)
	statuses := config.GroupProjects(groupProjectsProjects)
	c.Assert(statuses["IdentifierB"][0], Equals, groupProjectsResult["IdentifierB"][0])
	c.Assert(statuses["IdentifierA"][0], Equals, groupProjectsResult["IdentifierA"][0])
	c.Assert(statuses["IdentifierA"][1], Equals, groupProjectsResult["IdentifierA"][1])
}

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
	testProjects = []cctray.Project{
		cctray.Project{Name: "Project1", Transition: "Fixed"},
		cctray.Project{Name: "Project1", Transition: "Success"},
		cctray.Project{Name: "Notinconfig", Transition: "Failed"}}
	expectedWatches = []Watch{Watch{"Identifier 1", []string{"^Openstack.*"}, "_", []string{"Success", "Failed"}, "#api_test"},
		Watch{"Identifier 2", []string{"^Provision.*"}, "_", []string{"Success", "Failed"}, "#api_test"}}
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
	c.Assert(url, Equals, "project1")
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

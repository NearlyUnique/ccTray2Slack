package main

import "testing"

func equalString(t *testing.T, got string, expect string) {
	if got != expect {
		t.Errorf("Expected '%s' got '%s'", expect, got)
	}
}

var (
	testProjects = []Project{
		Project{Name: "Project1", Transition: "Fixed"},
		Project{Name: "Project1", Transition: "Success"},
		Project{Name: "Notinconfig", Transition: "Failed"}}
	expectedWatches = []Watch{Watch{"Identifier 1", []string{"^Openstack.*"}, "_", []string{"Success", "Failed"}, "#api_test"},
		Watch{"Identifier 2", []string{"^Provision.*"}, "_", []string{"Success", "Failed"}, "#api_test"}}
)

func equalInt(t *testing.T, got int, expect int) {
	if got != expect {
		t.Errorf("Expected %d got %d", expect, got)
	}
}

func TestProcess(t *testing.T) {
	config, _ := LoadConfig("testdata/config2.d")
	// When Processing a correct project but with wrong Transiton
	url, msg := config.Process(testProjects[0])
	// ... return empty url
	equalString(t, url, "")
	// When processing a correct project with corret transition
	url, msg = config.Process(testProjects[1])
	//... return message for the correct transition
	equalString(t, msg.Text, "Success text")
	// ... a non empty url
	equalString(t, url, "project1")
	// ... correct channel should be set
	equalString(t, msg.Channel, "#api_test")
	// When Processing a project which does not match a watched project
	url, msg = config.Process(testProjects[2])
	// ... return an empty url
	equalString(t, url, "")
}

func TestLoadConfig(t *testing.T) {

	config, _ := LoadConfig("testdata/config1.d/")
	if len(config.Watches) != 3 {
		t.Errorf("Expected 3 watches got %q", len(config.Watches))
	}
	if len(config.Remotes) != 2 {
		t.Errorf("Expected 3 watches got %q", len(config.Watches))
	}

	//TODO: Compare structs
	if config.Watches[0].Channel != expectedWatches[0].Channel {
		t.Errorf("Expected %v got %v", expectedWatches[0].Channel, config.Watches[0].Channel)
	}

	//TODO: Compare structs
	if config.Watches[1].ProjectRx[0] != expectedWatches[1].ProjectRx[0] {
		t.Errorf("Expected %v got %v", expectedWatches[1].ProjectRx[0], config.Watches[1].ProjectRx[0])
	}

	configPath := "testdata/config_which_does_not_exist.json"
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Errorf("Expecte to fail when loading config: %q", configPath)
	}

}

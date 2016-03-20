package main

import "testing"

func equalString(t *testing.T, got string, expect string) {
	if got != expect {
		t.Errorf("Expected %s got %s", expect, got)
	}
}

func equalInt(t *testing.T, got int, expect int) {
	if got != expect {
		t.Errorf("Expected %d got %d", expect, got)
	}
}

func TestInSlice(t *testing.T) {

	check := []string{"apa", "bepa", "cepa"}
	if InSlice("feg", check) {
		t.Error("feg is not in slice")
	}
	if !InSlice("bepa", check) {
		t.Error("bepa is in slice")
	}
}

var (
	expectedWatches = []Watch{Watch{[]string{"^Openstack.*"}, "_", []string{"Success", "Failed"}, "#api_test"},
		Watch{[]string{"^Provision.*"}, "_", []string{"Success", "Failed"}, "#api_test"}}
)

func TestLoadConfig(t *testing.T) {

	config, _ := LoadConfig("testdata/config1.json")
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

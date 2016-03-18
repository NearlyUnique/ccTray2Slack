package main

import "testing"

func TestInSlice(t *testing.T) {

	check := []string{"apa", "bepa", "cepa"}
	if InSlice("feg", check) {
		t.Error("feg is not in slice")
	}
	if !InSlice("bepa", check) {
		t.Error("bepa is in slice")
	}
}

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

func TestLoadConfig(t *testing.T) {

	config, _ := LoadConfig("testdata/config1.json")
	if len(config.Watches) != 3 {
		t.Errorf("Expected 3 watches got %q", len(config.Watches))
	}
	if len(config.Remotes) != 2 {
		t.Errorf("Expected 3 watches got %q", len(config.Watches))
	}

	equalString(t, config.Remotes[1], "otherhost")
	equalString(t, config.Watches[2].Transitions[0], "Broken")
	equalString(t, config.Watches[0].Transitions[1], "Failed")
	equalInt(t, len(config.Watches[2].ProjectRx), 2)

	equalString(t, config.Watches[0].ColorMapping["Success"], "#00ff00")

	equalString(t, config.Watches[1].SlackMsg.Attachements[0].Text, "This is attachments text")

	configPath := "testdata/config_which_does_not_exist.json"
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Errorf("Expecte to fail when loading config: %q", configPath)
	}
}

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

func TestLoadConfig(t *testing.T) {

	config := LoadConfig("testdata/config1.json")
	if len(config.Watches) != 3 {
		t.Errorf("Expected 3 watches got %q", len(config.Watches))
	}

	got := config.Watches[2].Transitions[0]
	expect := "Broken"
	if got != expect {
		t.Errorf("Expected %s got %s", expect, got)
	}

	got = config.Watches[0].Transitions[1]
	expect = "Failed"
	if got != expect {
		t.Errorf("Expected %s got %s", expect, got)
	}

	got_i := len(config.Watches[2].ProjectRx)
	expect_i := 2
	if got_i != expect_i {
		t.Errorf("Expected %d got %d", expect_i, got_i)
	}

}

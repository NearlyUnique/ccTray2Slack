package main

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/christer79/ccTray2Slack/cctray"
)

// func TestCompareDates(t *testing.T) {
//     p := Project {LastBuildTime:"2009-07-27T14:17:19"}
//     if p.LastBuildTime != "x" {
//         t.Error("Expected 1.5, got ", p)
//     }
// }

func Test_lastBuildTime_can_be_parsed_from_xml(t *testing.T) {
	p := cctray.Project{}
	raw := []byte(`<Project lastBuildTime="2009-07-27T14:17:19" name="N" activity="Sleeping" lastBuildStatus="Success" lastBuildLabel="3.0.754" webUrl="http://localhost:8153/cruise/tab/stage/detail/enterprisecorp-3/3.0.754/build/1"/>`)
	expected := time.Date(2009, time.July, 27, 14, 17, 19, 0, time.UTC)

	xml.Unmarshal(raw, &p)

	if !p.LastBuildTime.Equal(expected) {
		t.Errorf("got %v, expected %v", p.LastBuildTime, expected)
	}
}

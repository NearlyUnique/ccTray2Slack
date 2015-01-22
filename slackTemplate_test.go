package main

import (
	"testing"
	"time"
)

func Test_the_project_information_is_put_into_the_slack_msg_template(t *testing.T) {
	p := Project{
		"#name#",
		"#activity#",
		"#status#",
		"#label#",
		projTime{time.Date(2001, time.January, 15, 14, 16, 1, 0, time.UTC)},
		"#url#",
	}
	s := SlackMessage{`As of %time%, %project%, is %status% <%url%|%label%>`, "", "", ""}
	s.UpdateMessage(p)
	expected := `As of 2001-01-15 14:16:01, #name#, is #status# <#url#|#label#>`
	if s.Text != expected {
		t.Errorf("\nunexpected: %q\n       got: %q", s.Text, expected)
	}
}

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
		"#transition#",
	}

	var attachments []Attachement

	//If no attachement exists please do not Failed
	s := SlackMessage{`Project %project% <%url%|%label%>`, attachments, "", "", ""}

	s.UpdateMessage(p)

	attachments = append(attachments, Attachement{"", "", `has status %status%`})
	s = SlackMessage{`Project %project% <%url%|%label%>`, attachments, "", "", ""}

	s.UpdateMessage(p)
	expected := `Project #name# <#url#|#label#>`
	if s.Text != expected {
		t.Errorf("\nunexpected: %q\n       got: %q", s.Text, expected)
	}
	expected = `has status #transition#`
	if s.Attachements[0].Text != expected {
		t.Errorf("\nunexpected: %q\n       got: %q", s.Attachements[0].Text, expected)
	}

}

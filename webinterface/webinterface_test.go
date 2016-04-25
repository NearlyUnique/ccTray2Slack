package webinterface

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type WebInterfaceTestSuite struct{}

var _ = Suite(&WebInterfaceTestSuite{}) // Hook up gocheck into the "go test" runner.

type testStruct struct {
	Name   string
	Status string
}

var (
	testTemplate = `
  {{ range . }}
    {{ .Name }}<br>{{ .Status }}
  {{ end }}
  `

	testIndata = []testStruct{testStruct{"Name1", "Status1"}, testStruct{"Name2", "Status2"}}
)

func (s *WebInterfaceTestSuite) TestTemplate(c *C) {
	t := template.Must(template.New("test").Parse(testTemplate))
	b := new(bytes.Buffer)
	fmt.Println(testIndata)
	t.Execute(b, testIndata)
	fmt.Println(b.String())
}

package senders

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/chyle/chyle/types"
)

func TestBuildStdoutSender(t *testing.T) {
	config := stdoutConfig{FORMAT: "json"}
	assert.IsType(t, jSONStdoutSender{}, buildStdoutSender(config))

	config = stdoutConfig{FORMAT: "template", TEMPLATE: "{{.}}"}
	assert.IsType(t, templateStdoutSender{}, buildStdoutSender(config))
}

func TestJSONStdoutSender(t *testing.T) {
	buf := &bytes.Buffer{}

	s := jSONStdoutSender{buf}

	c := types.Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = []map[string]interface{}{
		{
			"id":   1,
			"test": "test",
		},
		{
			"id":   2,
			"test": "test",
		},
	}

	err := s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, `{"datas":[{"id":1,"test":"test"},{"id":2,"test":"test"}],"metadatas":{}}`, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations as json")
}

func TestTemplateStdoutSender(t *testing.T) {
	buf := &bytes.Buffer{}

	s := templateStdoutSender{buf, "{{ range $key, $value := .Datas }}{{$value.id}} : {{$value.test}} | {{ end }}"}

	c := types.Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = []map[string]interface{}{
		{
			"id":   1,
			"test": "test",
		},
		{
			"id":   2,
			"test": "test",
		},
	}

	err := s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, `1 : test | 2 : test | `, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations shaped on given template")
}
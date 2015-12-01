package main

import (
	"bytes"
	"testing"
)

func TestRenderSetsValues(t *testing.T) {
	in := bytes.NewBuffer([]byte("{{ .FOO }}, plain"))
	out := new(bytes.Buffer)

	err := render(in, out, map[string]string{"FOO": "from var"}, false)
	if err != nil {
		t.Fatalf("error from render: %s", err)
	}

	result := out.String()
	if result != "from var, plain" {
		t.Fatalf(
			"wrong rendered result. expected 'from var, plain', got '%s'",
			result,
		)
	}
}

func TestRenderEatsWhitespace(t *testing.T) {
	in := bytes.NewBuffer([]byte("\t{{- .FOO -}} \n   , plain"))
	out := new(bytes.Buffer)

	err := render(in, out, map[string]string{"FOO": "from var"}, true)
	if err != nil {
		t.Fatalf("error from render: %s", err)
	}

	result := out.String()
	if result != "from var, plain" {
		t.Fatalf(
			"wrong rendered result. expected 'from var, plain', got '%s'",
			result,
		)
	}
}

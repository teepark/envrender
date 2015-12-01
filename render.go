package main

import (
	"io"
	"io/ioutil"
	"regexp"
	"text/template"
)

var (
	startTag = regexp.MustCompile(`\s*{{-`)
	endTag   = regexp.MustCompile(`-}}\s*`)
)

func render(in io.Reader, out io.Writer, env map[string]string, wsEaters bool) error {
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	if wsEaters {
		input = startTag.ReplaceAll(input, []byte("{{"))
		input = endTag.ReplaceAll(input, []byte("}}"))
	}

	t, err := template.New("").Parse(string(input))
	if err != nil {
		return err
	}

	return t.Execute(out, env)
}

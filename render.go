package main

import (
	"io"
	"io/ioutil"
	"text/template"
)

func render(in io.Reader, out io.Writer, env map[string]string) error {
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	t, err := template.New("").Parse(string(input))
	if err != nil {
		return err
	}

	return t.Execute(out, env)
}

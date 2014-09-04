package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/template"

	"github.com/kballard/go-shellquote"
)

func main() {
	var (
		t *template.Template
		err error
		wr io.WriteCloser
	)

	env := whole_environ()

	/* parse flags */
	execp := flag.String("e", "", "command to exec after processing all files")
	flag.Parse()

	for _, path := range flag.Args() {

		if path == "-" {
			/* read stdin and parse as a template */
			text, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			t = template.Must(template.New("stdin").Parse(string(text)))
		} else {
			/* parse the file as a template */
			t = template.Must(template.ParseFiles(path))
		}

		if path == "-" {
			/* send output to stdout */
			wr = os.Stdout
		} else {
			/* empty and then open that same file for writing */
			wr, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
			if err != nil {
				panic(err)
			}
		}

		/* execute the template with environment vars,
		 * write the output back to the original file */
		err := t.Execute(wr, env)
		if err != nil {
			panic(err)
		}

		if path != "-" {
			/* close the file writer */
			err = wr.Close()
			if err != nil {
				panic(err)
			}
		}
	}

	if *execp != "" {
		/* parse shell arguments in the -e command */
		cmd, err := shellquote.Split(*execp)
		if err != nil {
			panic(err)
		}

		/* find the executable */
		cmdpath, err := exec.LookPath(cmd[0])
		if err != nil {
			panic(err)
		}

		/* and exec it */
		panic(syscall.Exec(cmdpath, cmd, os.Environ()))
	}
}

func whole_environ() map[string]string {
	env := os.Environ()
	result := map[string]string{}
	var item []string

	for _, pair := range env {
		item = strings.SplitN(pair, "=", 2)
		result[item[0]] = item[1]
	}

	return result
}

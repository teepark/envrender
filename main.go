package main

import (
	"flag"
	"fmt"
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
		t   *template.Template
		err error
		wr  io.WriteCloser
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
				failwith("file read <%s>", err)
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
				failwith("open file for write <%s>", err)
			}
		}

		/* execute the template with environment vars,
		 * write the output back to the original file */
		err := t.Execute(wr, env)
		if err != nil {
			failwith("render <%s>", err)
		}

		if path != "-" {
			/* close the file writer */
			err = wr.Close()
			if err != nil {
				failwith("close writer <%s>", err)
			}
		}
	}

	if *execp != "" {
		/* parse shell arguments in the -e command */
		cmd, err := shellquote.Split(*execp)
		if err != nil {
			failwith("shell split <%s>", err)
		}

		/* find the executable */
		cmdpath, err := exec.LookPath(cmd[0])
		if err != nil {
			failwith("executable path lookup <%s>", err)
		}

		/* and exec it */
		failwith("exec <%s>", syscall.Exec(cmdpath, cmd, os.Environ()))
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

func failwith(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

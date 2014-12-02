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

func whole_environ() map[string]string {
	env := os.Environ()
	result := map[string]string{}
	var item []string

	for _, pair := range env {
		item = strings.SplitN(pair, "=", 2)
		if len(item) < 2 {
			failwith("invalid item from os.Environ")
		}
		result[item[0]] = item[1]
	}

	return result
}

func failwith(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func processStdio(env map[string]string) {
	var (
		text []byte
		err  error
		t    *template.Template
	)

	/* read stdin to exhaustion */
	text, err = ioutil.ReadAll(os.Stdin)
	if err != nil {
		failwith("reading stdin <%s>", err)
	}

	/* parse it as a template */
	t, err = template.New("stdin").Parse(string(text))
	if err != nil {
		failwith("parsing stdin template <%s>", err)
	}

	/* render with the environ to stdout */
	if err = t.Execute(os.Stdout, env); err != nil {
		failwith("rendering stdin template <%s>", err)
	}
}

func processFile(sourcePath, destPath string, env map[string]string) {
	var (
		t   *template.Template
		err error
		wr  io.WriteCloser
	)

	/* parse the file as a template */
	t, err = template.ParseFiles(sourcePath)
	if err != nil {
		failwith("parsing file %s <%s>", sourcePath, err)
	}

	/* empty and open the destination file for writing */
	wr, err = os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		failwith("open file for write <%s>", err)
	}

	/* render with the environ to the open file */
	if err = t.Execute(wr, env); err != nil {
		failwith("render to file <%s>", err)
	}

	/* close the file */
	if err = wr.Close(); err != nil {
		failwith("closing writer <%s>", err)
	}
}

func execCmd(cmd string) {
	var (
		err     error
		args    []string
		cmdpath string
	)

	/* shell-aware split the arguments */
	args, err = shellquote.Split(cmd)
	if err != nil {
		failwith("shell split <%s>", err)
	}

	/* find the executable */
	cmdpath, err = exec.LookPath(args[0])
	if err != nil {
		failwith("executable path lookup <%s>", err)
	}

	/* exec it */
	if err = syscall.Exec(cmdpath, args, os.Environ()); err != nil {
		failwith("exec <%s>", err)
	}
}

func main() {
	env := whole_environ()

	/* parse flags */
	execp := flag.String("e", "", "command to exec after processing all files")
	flag.Parse()

	/* process files */
	for _, path := range flag.Args() {
		if path == "-" {
			processStdio(env)
		} else if i := strings.Index(path, ":"); i >= 0 {
			processFile(path[:i], path[i+1:], env)
		} else {
			processFile(path, path, env)
		}
	}

	/* exec cmd if given */
	if *execp != "" {
		execCmd(*execp)
	}
}

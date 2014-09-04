package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/template"

	"github.com/kballard/go-shellquote"
)

func main() {
	env := whole_environ()

	/* parse flags */
	execp := flag.String("e", "", "command to exec after processing all files")
	flag.Parse()

	for _, path := range flag.Args() {

		/* parse the file as a template */
		t, err := template.ParseFiles(path)
		if err != nil {
			fmt.Println(path, "parsing:", err)
			os.Exit(1)
		}

		/* empty and then open that same file for writing */
		wr, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
		if err != nil {
			fmt.Println(path, "writing:", err)
			os.Exit(1)
		}

		/* execute the template with environment vars,
		 * write the output back to the original file */
		err = t.Execute(wr, env)
		if err != nil {
			fmt.Println(path, "executing:", err)
			os.Exit(1)
		}

		/* close the file writer */
		err = wr.Close()
		if err != nil {
			fmt.Println(path, "closing:", err)
			os.Exit(1)
		}
	}

	if *execp != "" {
		cmd, err := shellquote.Split(*execp)
		if err != nil {
			fmt.Println("cmd parse:", err)
			os.Exit(1)
		}
		cmdpath, err := exec.LookPath(cmd[0])
		if err != nil {
			fmt.Println("cmd lookup:", err)
			os.Exit(1)
		}
		syscall.Exec(cmdpath, cmd, os.Environ())
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

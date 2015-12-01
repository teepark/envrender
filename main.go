package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	var wsEaters bool

	flag.BoolVar(
		&wsEaters,
		"w",
		true,
		"support whitespace {{- eating -}} tag-delimiters",
	)

	flag.Usage = usage
	flag.Parse()

	env := envMap()
	jobs, rest := splitArgs(flag.Args())

	var stdinUsed, stdoutUsed bool

	for _, job := range jobs {
		input, err := getInput(job.src)
		if err != nil {
			log.Fatal(err)
		}
		if job.src == "-" {
			if stdinUsed {
				log.Fatal("attempted to use stdin twice")
			}
			stdinUsed = true
		}

		output, err := getOutput(job.dest)
		if err != nil {
			log.Fatal(err)
		}
		if job.dest == "-" {
			if stdoutUsed {
				log.Fatal("attempted to use stdout twice")
			}
			stdoutUsed = true
		}

		err = render(input, output, env, wsEaters)
		if err != nil {
			log.Fatal(err)
		}

		if input != os.Stdin {
			err = input.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
		if output != os.Stdout {
			err = output.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if len(rest) > 0 {
		cmd, err := exec.LookPath(rest[0])
		if err != nil {
			log.Fatal(err)
		}

		err = syscall.Exec(cmd, rest, os.Environ())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"Usage of %s: %s [flags] src:dest... prog...\n",
		os.Args[0], os.Args[0],
	)
	flag.PrintDefaults()
}

func getInput(src string) (io.ReadCloser, error) {
	if src == "-" {
		return os.Stdin, nil
	}
	return os.Open(src)
}

func getOutput(dest string) (io.WriteCloser, error) {
	if dest == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func main() {
	var createPerms string
	flag.StringVar(
		&createPerms,
		"create_perms",
		"644",
		"octal permissions to set on created output files",
	)

	flag.Usage = usage
	flag.Parse()

	perms, err := parsePerms(createPerms)
	if err != nil {
		log.Fatalf("error parsing -create_perms: %s", err)
	}

	newMask := 0777 - perms
	oldMask := syscall.Umask(newMask)
	newMask |= oldMask
	syscall.Umask(newMask)

	env := envMap()
	jobs, rest := splitArgs(flag.Args())

	var stdinUsed, stdoutUsed bool

	for _, job := range jobs {
		if job.src == "-" {
			if stdinUsed {
				log.Fatal("attempted to use stdin twice")
			}
			stdinUsed = true
		}
		if job.dest == "-" {
			if stdoutUsed {
				log.Fatal("attempted to use stdout twice")
			}
			stdoutUsed = true
		}

		input, err := getInput(job.src)
		if err != nil {
			log.Fatal(err)
		}
		output, err := getOutput(job.dest)
		if err != nil {
			log.Fatal(err)
		}

		err = render(input, output, env)
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

	syscall.Umask(oldMask)

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
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] src:dest... prog...\n", os.Args[0])
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
	return os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
}

func parsePerms(perms string) (int, error) {
	n, err := strconv.ParseUint(perms, 8, 64)
	return int(n), err
}

package main

import "strings"

type job struct {
	src, dest string
}

func splitArgs(args []string) ([]job, []string) {
	var jobs []job

	for i, arg := range args {
		if arg == "--" {
			return jobs, args[i+1:]
		}

		split := strings.SplitN(arg, ":", 2)
		if len(split) < 2 {
			return jobs, args[i:]
		}

		if split[0] == "" {
			split[0] = split[1]
		} else if split[1] == "" {
			split[1] = split[0]
		}

		jobs = append(jobs, job{
			src:  split[0],
			dest: split[1],
		})
	}

	return jobs, nil
}

package main

import "testing"

func TestSplitArgsSplits(t *testing.T) {
	jobs, rest := splitArgs([]string{
		"foo:bar",
		"spam:eggs",
		"/bin/sh",
		"-c",
		"echo hi",
	})

	Eq(t, "job count", 2, len(jobs))
	Eq(t, "jobs[0]", job{"foo", "bar"}, jobs[0])
	Eq(t, "jobs[1]", job{"spam", "eggs"}, jobs[1])

	Eq(t, "rest count", 3, len(rest))
	Eq(t, "rest[0]", "/bin/sh", rest[0])
	Eq(t, "rest[1]", "-c", rest[1])
	Eq(t, "rest[2]", "echo hi", rest[2])
}

func TestSplitArgsStopsAtDoubleDash(t *testing.T) {
	jobs, rest := splitArgs([]string{
		"foo:bar",
		"--",
		"spam:eggs",
	})

	Eq(t, "job count", 1, len(jobs))
	Eq(t, "job[0]", job{"foo", "bar"}, jobs[0])

	Eq(t, "rest count", 1, len(rest))
	Eq(t, "rest[0]", "spam:eggs", rest[0])
}

func TestSplitArgsCopiesSrcToDest(t *testing.T) {
	jobs, _ := splitArgs([]string{"foo:"})

	Eq(t, "job count", 1, len(jobs))
	Eq(t, "jobs[0].src", "foo", jobs[0].src)
	Eq(t, "jobs[0].dest", "foo", jobs[0].dest)
}

func TestSplitArgsCopiesDestToSrc(t *testing.T) {
	jobs, _ := splitArgs([]string{":foo"})

	Eq(t, "job count", 1, len(jobs))
	Eq(t, "jobs[0].dest", "foo", jobs[0].dest)
	Eq(t, "jobs[0].src", "foo", jobs[0].src)
}

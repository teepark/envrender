package main

import "testing"

func Eq(t *testing.T, msg string, expected, got interface{}) {
	if expected != got {
		t.Fatalf("%s: expected %+v, got %+v", msg, expected, got)
	}
}

func NotEq(t *testing.T, msg string, expected, got interface{}) {
	if expected == got {
		t.Fatalf("%s: expected NOT %+v, got %+v", msg, expected, got)
	}
}

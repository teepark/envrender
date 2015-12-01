package main

import (
	"os"
	"testing"
)

func TestEnvMap(t *testing.T) {
	if err := os.Setenv("foo", "bar"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv("foo"); err != nil {
			t.Fatal(err)
		}
	}()

	m := envMap()
	if m["foo"] != "bar" {
		t.Fatalf("expected 'foo' to be 'bar', got '%s'", m["foo"])
	}
}

func TestEnvMapMissing(t *testing.T) {
	m := envMap()
	if v, ok := m["DONT_SET_ME"]; ok {
		t.Fatalf("'DONT_SET_ME' was found? '%s'", v)
	}
}

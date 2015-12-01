package main

import (
	"log"
	"os"
	"strings"
)

func envMap() map[string]string {
	env := os.Environ()
	result := map[string]string{}
	var item []string

	for _, pair := range env {
		item = strings.SplitN(pair, "=", 2)
		if len(item) < 2 {
			log.Fatalf("invalid item from os.Environ")
		}
		result[item[0]] = item[1]
	}

	return result
}

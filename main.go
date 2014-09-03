package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func main() {
	for _, path := range os.Args[1:] {

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
		err = t.Execute(wr, environ())
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
}

func environ() map[string]string {
	env := os.Environ()
	result := map[string]string{}
	var item []string

	for _, pair := range env {
		item = strings.SplitN(pair, "=", 2)
		result[item[0]] = item[1]
	}

	return result
}

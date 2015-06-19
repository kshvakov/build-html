package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/flosch/pongo2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	input       = ""
	output      = ""
	contextFile = ""
	context     pongo2.Context
	trimSpaces  bool
)

func init() {

	flag.StringVar(&input, "i", ".", "input dir")
	flag.StringVar(&output, "o", "", "output dir")
	flag.StringVar(&contextFile, "c", "", "context json (pongo2.Context)")
	flag.BoolVar(&trimSpaces, "s", true, "trim spaces")

	flag.Usage = func() {

		fmt.Fprint(os.Stderr, "./build-html -i=template/src -o=template/html -c=context.json -s=false\n")

		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()

	if output == "" {

		flag.Usage()

		os.Exit(2)
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {

		log.Fatal(err)
	}

	if contextFile != "" {

		if _, err := os.Stat(contextFile); os.IsNotExist(err) {

			log.Fatal(err)
		}

		if data, err := ioutil.ReadFile(contextFile); err == nil {

			if err := json.Unmarshal(data, &context); err != nil {

				log.Fatal(err)
			}

		} else {

			log.Fatal(err)
		}
	}

	filepath.Walk(input, func(path string, fi os.FileInfo, _ error) error {

		name, err := filepath.Rel(input, path)

		if name == "internal" || err != nil {

			return nil
		}

		if fi.IsDir() {

			if err := os.MkdirAll(output+"/"+name, 0755); err != nil {

				log.Fatal(err)
			}

			return nil
		}

		if !strings.HasSuffix(path, ".tpl") {

			return nil
		}

		tpl := pongo2.Must(pongo2.FromFile(path))

		if content, err := tpl.Execute(context); err == nil {

			ioutil.WriteFile(output+"/"+name, trim(content), 0644)

		} else {

			log.Fatal(err)
		}

		return nil
	})
}

func trim(content string) []byte {

	if !trimSpaces {

		return []byte(content)
	}

	i := 0

	html := make([]byte, len(content))

	previousSpace := false

	for _, x := range []byte(content) {

		if x == ' ' || x == '\t' || x == '\n' || x == '\r' || x == '\f' {

			if !previousSpace {

				previousSpace = true

				html[i] = ' '

				i++
			}

		} else {

			previousSpace = false

			html[i] = x

			i++
		}
	}

	return html[:i]
}

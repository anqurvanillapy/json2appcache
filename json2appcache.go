package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Manifest struct {
	Res    []string `json:"initial"`
	Config string   `json:"configPath"`
}

const tpl = `
CACHE MANIFEST
# {{ .timestamp }}

CACHE:
{{ .timestamp }}

NETWORK:
*

FALLBACLK:
{{}}`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: json2appcache manifest.json")
		os.Exit(1)
	}

	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	raw, err := ioutil.ReadFile(os.Args[1])
	check(err)

	var m Manifest
	json.Unmarshal(raw, &m)

	timestamp := time.Now()

	// Avoid using relative paths for caching.
	getAbsPathAndEntry := func(fpath string, isConfig bool) (string, string) {
		var root string

		if isConfig {
			root = "/resource/"
		} else {
			root = "/"
		}

		fdir := filepath.Dir(fpath)
		fbase := path.Base(fpath)
		fentry := fbase[:strings.LastIndex(fbase, "_")]

		fmt.Println(fentry)

		var buf bytes.Buffer
		buf.WriteString(root)
		buf.WriteString(fpath)

		return fpath, fentry
	}

	t, err := template.New("appcache").Parse(tpl)
	check(err)

	// err = t.Execute(os.Stdout)
	check(err)
}

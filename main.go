package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/artonge/go-gtfs"
)

func main() {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		exit("unable to create temporary directory", err)
	}
	defer os.RemoveAll(dir)

	resp, err := http.Get("https://transitfeeds.com/p/mta/79/latest/download")
	if err != nil {
		exit("unable to download GTFS feed", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exit("unable to read response body", err)
	}

	buf := bytes.NewReader(data)
	r, err := zip.NewReader(buf, int64(buf.Len()))
	if err != nil {
		exit("unable to read zip archive", err)
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			exit("unable to open file in zip archive", err)
		}

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			exit("unable to read file in zip archive", err)
		}

		filename := filepath.Join(dir, f.Name)
		if err := ioutil.WriteFile(filename, data, 0644); err != nil {
			exit("unable to write temporary file", err)
		}

		if err := rc.Close(); err != nil {
			exit("unable to close file in zip archive", err)
		}
	}

	g, err := gtfs.Load(dir, nil)
	if err != nil {
		exit("unable to load GTFS feed from temporary directory", err)
	}

	fmt.Printf("Successfully GTFS feed for %s!", g.Agency.Name)
}

func exit(msg string, err error) {
	fmt.Printf("%s: %v", msg, err)
	os.Exit(1)
}

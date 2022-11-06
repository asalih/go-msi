package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/asalih/go-msi"
)

func main() {
	files, err := filepath.Glob("testdata/*.msi")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file == "testdata/.DS_Store" {
			continue
		}

		rdr, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			panic(err)
		}

		msi, err := msi.Open(rdr)
		if err != nil {
			panic(err)
		}

		streams := msi.Streams()

		for {
			streamName := streams.Next()
			if streamName == "" {
				break
			}

			if path.Ext(streamName) != ".cab" {
				continue
			}

			dStream, err := msi.ReadStream(streamName)
			if err != nil {
				panic(err)
			}

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(dStream)
			if err != nil {
				panic(err)
			}

			fmt.Println(len(buf.Bytes()))
		}
	}
}

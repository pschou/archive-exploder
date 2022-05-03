// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_helloworld_get]

// Package helloworld provides a set of Cloud Functions samples.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/pschou/tease"
)

type Archive interface {
	Next() (name, path string, r io.Reader, err error)
	IsEOF() bool
	Close()
	Type() string
}
type formatTest struct {
	Test func(*tease.Reader) bool
	Read func(*tease.Reader, int64) (Archive, error)
	Type string
}

var formatTests = []formatTest{}

var version = "test"
var debug *bool
var maxRec *int

// Main is a function to fetch the HTTP repodata from a URL to get the latest
// package list for a repo
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Archive Exploder,  Version: %s\n\nUsage: %s [options...]\n\n", version, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nFormats supported:")
		for _, f := range formatTests {
			fmt.Fprintln(os.Stderr, "  -", f.Type)
		}
	}

	maxRec = flag.Int("r", 1, "Levels of recusion (archives-inside-archives) to expand")
	var output = flag.String("output", ".", "Path to put the extracted file(s)")
	var input = flag.String("input", "", "Path to input archive file")
	debug = flag.Bool("debug", false, "Turn on debug, more verbose")
	flag.Parse()

	if *input == "" {
		flag.Usage()
		os.Exit(0)
	}

	if *maxRec == 0 || *maxRec < -1 {
		fmt.Println("Invalid recursion value")
		os.Exit(1)
	}
	//repoPath := strings.TrimSuffix(strings.TrimPrefix(*inRepoPath, "/"), "/")

	info, err := os.Stat(*input)
	if err != nil {
		log.Fatalf("failed to stat file: %s", err)
	}

	f, err := os.Open(*input)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	//fmt.Printf("%+v %s\n", rdr, err)
	explode(*output, f, info.Size(), 0)
}

func explode(filePath string, in io.Reader, size int64, rec int) (err error) {
	if rec == *maxRec {
		// If we have reached the max depth, print out any file / archive without testing
		//fmt.Println("reached max")
		var n int64
		n, err = writeFile(filePath, in)
		if size >= 0 && n != size {
			log.Fatal("Reader.MaxDepth: copied file size does not match")
		}
		return
	}

	tr := tease.NewReader(in)
	//defer tr.Close()
	//var arch Archive

	//if arch, err = readISO9660(tr); err == nil {
	//} else if arch, err = readGzip(tr); err == nil {
	//} else {
	//	fmt.Println("not an archive", file)
	//}

	//fmt.Println("dir", dir, "file", file, "r", r, "err", err)

	matches := []formatTest{}

	for _, ft := range formatTests {
		if ft.Test(tr) {
			matches = append(matches, ft)
		}
	}

	tr.Seek(0, io.SeekStart)

	switch len(matches) {
	case 0:
		if *debug {
			fmt.Println("no archive match for", filePath)
		}
		var n int64
		tr.Seek(0, io.SeekStart)
		tr.Pipe()
		n, err = writeFile(filePath, tr)
		if size >= 0 && n != size {
			log.Fatal("copied file size,", n, ", ane expected,", size, ", does not match")
		}
	case 1:
		// We found only one potential archive match, go ahead and explode it.
		tr.Seek(0, io.SeekStart)
		ft := matches[0]
		if *debug {
			fmt.Println("archive match for", filePath, "type", ft.Type)
		}
		if arch, err := ft.Read(tr, size); err == nil {
			//defer arch.Close()
			for !arch.IsEOF() {
				//if ft.Type == "gzip" {
				//	fmt.Println("calling next")
				//}
				//var a_dir, a_file string
				//	var r io.Reader
				//var err error
				a_dir, a_file, r, err := arch.Next()
				to_read := int64(-1)
				if lr, ok := (r).(*io.LimitedReader); ok {
					to_read = lr.N
					//fmt.Printf("limited readertype %+v\n", lr)
				}
				//if ft.Type == "gzip" {
				//fmt.Println("a_dir", a_dir, "a_file", a_file, r)
				//}
				if err != nil {
					if err != io.EOF {
						fmt.Println("  error advancing to next file", err)
					}
					break
				}
				explode(path.Join(filePath, a_dir, a_file), r, to_read, rec+1)
				//pos, _ := tr.Seek(0, io.SeekCurrent)
				//if pos == size {
				//	break
				//}
			}
		} else {
			if *debug {
				fmt.Println("Archive", filePath, "failed to expand with type", ft.Type, ", ", err)
			}
			tr.Seek(0, io.SeekStart)
			_, err = writeFile(filePath, tr)
		}
	default:
		if *debug {
			fmt.Println("Archive", filePath, "matches multiple formats, what to do?")
			for _, ft := range matches {
				fmt.Println("  ", ft.Type)
			}
		}
		_, err = writeFile(filePath, tr)
	}

	return
}

func writeFile(filePath string, in io.Reader) (int64, error) {
	if *debug {
		fmt.Println("Writing out file", filePath)
	}
	dir, _ := path.Split(filePath)
	ensureDir(dir)
	out, err := os.Create(filePath)
	if err != nil {
		log.Println("= Error creating file", filePath, "err:", err)
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, in)
}

/*
func testFormats(tr *tease.Reader) (ret []func(*tease.Reader) (Archive, error)) {
	if testISO9660(tr) {
		ret = append(ret, (func(*tease.Reader) (Archive, error))(readISO9660))
	}
}
*/

package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	//"github.com/kdomanski/iso9660"

	"github.com/pschou/tease"
)

type GzipFile struct {
	buf_reader *bufio.Reader
	gz_reader  *gzip.Reader
	tr_reader  *tease.Reader
	eof        bool
	count      int
}

func init() {
	formatTests = append(formatTests, formatTest{
		Test: testGzip,
		Read: readGzip,
		Type: "gzip",
	})
}

func testGzip(tr *tease.Reader) bool {
	tr.Seek(0, io.SeekStart)
	buf := make([]byte, 2)
	tr.Read(buf)
	tr.Seek(0, io.SeekStart)
	return bytes.Compare(buf, []byte{0x1f, 0x8b}) == 0
}

func readGzip(tr *tease.Reader) (Archive, error) {
	//tr.Seek(0, io.SeekStart)
	//buf := make([]byte, 20)
	//tr.Read(buf)
	//fmt.Printf("   %x\n", buf)
	tr.Seek(0, io.SeekStart)
	br := bufio.NewReader(tr)
	gzr, err := gzip.NewReader(br)
	if err != nil {
		//tr.Seek(0, io.SeekStart)
		//tr.Stats()
		if *debug {
			fmt.Println("Error reading gzip", err)
		}
		//log.Fatal("Error reading gzip", err)
		return nil, err
	}
	gzr.Multistream(false)
	//if err != nil {
	//	return nil, err
	//}

	//fmt.Println("gzip returning reader", gzr, tr)

	ret := GzipFile{
		buf_reader: br,
		gz_reader:  gzr,
		tr_reader:  tr,
		eof:        false,
	}

	tr.Pipe()
	return &ret, nil
}

func (i *GzipFile) Type() string {
	return "gzip"
}

func (i *GzipFile) IsEOF() bool {
	return i.eof
}

func (c *GzipFile) Close() {
	//if c.buf_reader != nil {
	//	c.buf_reader.Reset(nil)
	//}
	if c.gz_reader != nil {
		c.gz_reader.Close()
	}
}

func (i *GzipFile) Next() (path, name string, r io.Reader, err error) {
	//fmt.Println("gzip next called")
	//i.tr_reader.Pipe()
	if i.count == 0 {
		i.count = 1
		return ".", "pt_1", i.gz_reader, nil
	}
	//i.eof = true
	//return "", "", nil, io.EOF

	/*
		buf := make([]byte, 100)
		n := 100
		for err != io.EOF && n == 100 {
			fmt.Println("dumping out rest of file")
			n, err = i.gz_reader.Read(buf)
		}
	*/
	fmt.Println("reset")
	err = i.gz_reader.Reset(i.buf_reader)
	if err != nil {
		i.eof = true
		//if err == io.EOF {
		//}
		return "", "", nil, err
	}
	i.gz_reader.Multistream(false)
	i.count++
	return ".", fmt.Sprintf("pt_%d", i.count), i.gz_reader, nil
}

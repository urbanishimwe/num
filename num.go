// Copyright 2020 Urban I. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Data storage units
const (
	_ = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

var (
	x, o, d, c bool
	s, stdin   bool
	verbose, b bool
	u, format  string
	_fmt       string
	fCount     int
)

func main() {
	parseFormater()

	// command-line inputs
	leftArgs := flag.Args()
	if len(leftArgs) > 0 {
		for _, str := range leftArgs {
			if c {
				for _, v := range str {
					fmt.Println(outputChar(v))
				}
			} else {
				fmt.Println(outputInt(str))
			}
		}
	}
	var scanner *bufio.Scanner
	if stdin {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		scanner = bufio.NewScanner(nonBlockStdin(0))
	}
	readFromScanner(scanner)
}

func outputInt(v string) string {
	t := 1
	uToB := []byte(u)
	switch {
	case endsWithFold(uToB, []byte("kb")):
		t = KB
	case endsWithFold(uToB, []byte("mb")):
		t = MB
	case endsWithFold(uToB, []byte("gb")):
		t = GB
	case endsWithFold(uToB, []byte("tb")):
		t = TB
	case endsWithFold(uToB, []byte("pb")):
		t = PB
	default:
		t = 1
	}
	t2, cut := 1, 2
	vToB := []byte(v)
	switch {
	case endsWithFold(vToB, []byte("kb")):
		t2 = KB
	case endsWithFold(vToB, []byte("mb")):
		t2 = MB
	case endsWithFold(vToB, []byte("gb")):
		t2 = GB
	case endsWithFold(vToB, []byte("tb")):
		t2 = TB
	case endsWithFold(vToB, []byte("pb")):
		t2 = PB
		cut = 1
	default:
		cut = 0
	}
	v = v[:len(v)-cut]
	dt, err := strconv.ParseInt(v, 0, 64)
	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "output int error: %q", err)
		}
		return ""
	}
	dt = (dt * int64(t2)) / int64(t)
	vInts := make([]interface{}, fCount)
	for i := range vInts {
		vInts[i] = dt
	}
	return fmt.Sprintf(_fmt, vInts...)
}

func outputChar(v rune) string {
	vChar := make([]interface{}, fCount)
	for i := 0; i < fCount; i++ {
		vChar[i] = v
	}
	return fmt.Sprintf(_fmt, vChar...)
}

func endsWithFold(a, b []byte) bool {
	if len(a) < len(b) {
		return false
	}
	return bytes.EqualFold(a[len(a)-len(b):], b)
}

func parseFormater() {
	if format != "" {
		_fmt = format
		countFmt()
	} else {
		if x {
			_fmt += "0x%x "
			fCount++
		}
		if d {
			_fmt += "%v "
			fCount++
		}
		if o {
			_fmt += "0%o "
			fCount++
		}
		if b {
			_fmt += "0b%b "
			fCount++
		}
		if s {
			_fmt += "%q "
			fCount++
		}
		if _fmt == "" {
			_fmt = "%v"
			fCount++
		}
	}
}

func countFmt() {
	rd := strings.NewReader(format)
	var r0, r1 rune
	var err error
	for err == nil {
		r0, _, err = rd.ReadRune()
		if r0 == '%' {
			r1, _, err = rd.ReadRune()
			if r1 != '%' {
				fCount++
			}
		}
	}
}

func readFromScanner(reader *bufio.Scanner) {
	if c {
		reader.Split(bufio.ScanRunes)
	} else {
		reader.Split(bufio.ScanWords)
	}
	for reader.Scan() {
		if c {
			token := reader.Text()
			fmt.Println(outputChar([]rune(token)[0]))
			continue
		}
		fmt.Println(outputInt(string(reader.Bytes())))
	}
	if reader.Err() != nil && verbose {
		fmt.Fprintf(os.Stderr, "scanner error: %q", reader.Err())
	}
}

type nonBlockStdin int

func (b nonBlockStdin) Read(p []byte) (int, error) {
	fs, err := os.Stdin.Stat()
	if err != nil {
		return 0, err
	}
	if fs.Size() > 0 {
		return os.Stdin.Read(p)
	}
	return 0, io.EOF
}

func init() {
	flag.Usage = func() {
		var title = `Num is the CLI to transform integers and characters.
Results of multiple inputs are separated by newline(\n).

USAGE: flags must be entered before inputs in the command-line.
Parsing order start with command line inputs and then standard input or terminal pipe.

OUTPUT Order: [hex-format] [base 10] [octal] [binary] [utf-8].
`

		fmt.Fprintln(os.Stderr, title)
		flag.PrintDefaults()
	}
	flag.BoolVar(&x, "x", false, "append output in hexadecimal")
	flag.BoolVar(&d, "d", false, "apend output in decimal(default)")
	flag.BoolVar(&o, "o", false, "append output in octal")
	flag.BoolVar(&b, "b", false, "append output in binary")
	flag.BoolVar(&s, "s", false, "append output of an integer converted to a character")
	flag.BoolVar(&c, "c", false, "treat input as utf-8 characters and convert them to integers.")
	flag.BoolVar(&verbose, "v", false, "verbose: prints parser errors on standard error stream")
	flag.StringVar(&u, "u", "b", "data units for the output i.e KB, MB, GB, TB, or PB. can be used with -f")
	flag.StringVar(&format, "f", "", "output format with valid printf flags. It will be used instead of other output flags")
	flag.BoolVar(&stdin, "stdin", false, "allow blocking for inputs from standard input stream")
	flag.Parse()
}

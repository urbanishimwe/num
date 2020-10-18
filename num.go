// Copyright 2020 Urban I. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// Data storage units
const (
	_ = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
	EB
)

var (
	x, o, d, c bool
	s, stdin   bool
	verbose, b bool
	u, format  string
	_fmt, f    string
	fCount     int
)

func main() {
	flag.BoolVar(&x, "x", false, "append output in hexadecimal")
	flag.BoolVar(&d, "d", false, "apend output in decimal(default)")
	flag.BoolVar(&o, "o", false, "append output in octal")
	flag.BoolVar(&b, "b", false, "append output in binary")
	flag.BoolVar(&s, "s", false, "append output of an integer converted to a character")
	flag.BoolVar(&c, "c", false, "treat input as characters and convert them to integers.")
	flag.BoolVar(&verbose, "v", false, "verbose: prints parser errors on standard error stream")
	flag.StringVar(&u, "u", "b", "data units for the output i.e B, KB, MB, GB, TB, PB or EB")
	flag.StringVar(&format, "format", "", "custom output format with valid printf flags, this override bases flags")
	flag.StringVar(&f, "f", "", "name of the file to read inputs from")
	flag.BoolVar(&stdin, "stdin", false, "allow blocking for inputs from standard input stream")
	flag.Parse()
	parseFormater()

	// command-line inputs
	leftArgs := flag.Args()
	if len(leftArgs) > 0 {
		for i := range leftArgs {
			if c {
				fmt.Println(outputChar(leftArgs[i]))
			} else {
				fmt.Println(outputInt(leftArgs[i]))
			}
		}
	}

	var reader *bufio.Scanner
	switch {

	// check for file inputs
	case f != "":
		f, err := os.Open(f)
		mayBeExit(err)
		reader = bufio.NewScanner(f)
		readFromScanner(reader)

	// check for inputs from pipe and stdin
	default:
		stat, err := os.Stdin.Stat()
		// no other inputs we got
		if len(leftArgs) < 1 && reader == nil {
			mayBeExit(err)
		}
		if stat == nil {
			return
		}
		charDevice := stat.Mode()&os.ModeDevice != 0 && stat.Mode()&os.ModeCharDevice != 0
		// read from a pipe?
		if !charDevice {
			reader = bufio.NewScanner(os.Stdin)
			readFromScanner(reader)
		}
		// block for input?
		if charDevice && stdin {
			reader = bufio.NewScanner(os.Stdin)
			readFromScanner(reader)
		}
	}
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
	case endsWithFold(uToB, []byte("eb")):
		t = EB
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
	case endsWithFold(vToB, []byte("eb")):
		t2 = EB
	case endsWithFold(vToB, []byte("b")):
		cut = 1
	default:
		cut = 0
	}
	v = v[:len(v)-cut]
	vInt := make([]interface{}, fCount)
	dt, err := strconv.ParseInt(v, 0, 64)
	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "output int error: %q", err)
		}
		return ""
	}
	for i := range vInt {
		if t > t2 {
			vInt[i] = float64(dt) * float64(t2) / float64(t)
			continue
		}
		vInt[i] = dt * int64(t2/t)
	}
	return fmt.Sprintf(_fmt, vInt...)
}

func outputChar(v string) string {
	var res strings.Builder
	vChar := make([]interface{}, fCount)
	for _, val := range v {
		for i := 0; i < fCount; i++ {
			vChar[i] = val
		}
		res.WriteString(fmt.Sprintf(_fmt, vChar...))
	}
	buf := res.String()
	return buf
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
	if fCount > 1 && c {
		_fmt += "\n"
	}
}

func readFromScanner(reader *bufio.Scanner) {
	if !c {
		reader.Split(bufio.ScanWords)
	}
	for reader.Scan() {
		if c {
			fmt.Println(outputChar(stringy(reader.Bytes())))
			continue
		}
		fmt.Println(outputInt(stringy(reader.Bytes())))
	}
	if reader.Err() != nil && verbose {
		fmt.Fprintf(os.Stderr, "scanner error: %q", reader.Err())
	}
}

func mayBeExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}

func stringy(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
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

func init() {
	flag.Usage = func() {
		const title = `
Num is the CLI to transform integers and characters.

USAGE: flags must be entered before inputs in the command-line
`
		fmt.Fprintln(os.Stderr, title)
		flag.PrintDefaults()
	}
}

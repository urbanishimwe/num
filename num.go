// Copyright 2020 Urban I. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
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
	EB
)

// V array of values
type V []string

var (
	x, o, d, b, stdin bool
	u, f              string
	fC                int
	v                 = new(V)
)

func init() {
	flag.Usage = func() {
		const title = `
Num is the command line tool to parse integers in/from different bases and data units.

USAGE: (flag must be entered before input in the command line)
`
		fmt.Println(title)
		flag.PrintDefaults()
		const footer = `
you can add multiple bases like -x -d -o -b. the output of every input will be on a single line.
Inputs can be separated with white space or newline. Data unit(flag -u) receive single units,
in converting from Da to Db where Da > Db it is better to use custom format
that receive floating e.g: num -u GB -f="%0.4fGB" 10TB, note that everything is case insensitive.
`
		fmt.Println(footer)
	}
}

func main() {
	flag.BoolVar(&x, "x", false, "output in hexadecimal")
	flag.BoolVar(&d, "d", false, "output in decimal(default)")
	flag.BoolVar(&o, "o", false, "output in octal")
	flag.BoolVar(&b, "b", false, "output in binary")
	flag.StringVar(&u, "u", "b", "data units for the output i.e KB, MB, GB, TB, PB or EB")
	flag.StringVar(&f, "f", "", "custom output format with valid printf flags, it does not affect data unit but it will override other formats")
	flag.IntVar(&fC, "f-count", 1, "number of flags parsed in custom format(--f). e.g --f '%q %x' --f-count must be 2")
	flag.BoolVar(&stdin, "stdin", false, "read input from stdin pipe line by line until EOF")
	flag.Parse()
	if stdin {
		var in string
		buf := bufio.NewScanner(os.Stdin)
		for buf.Scan() {
			in += " " + buf.Text()
		}
		v.Set(in)
	} else {
		v.Set(strings.Join(flag.Args(), " "))
	}
	for _, val := range *v {
		fmt.Println(output(val))
	}
}

func output(v string) string {
	var format string
	var fCount int
	if f != "" {
		format = f
		fCount = fC
	} else {
		if x {
			format += "0x%x "
			fCount++
		}
		if d {
			format += "%v "
			fCount++
		}
		if o {
			format += "0o%o "
			fCount++
		}
		if b {
			format += "0b%b "
			fCount++
		}
		if !(x || b || o || d) {
			format += "%v "
			fCount++
		}
	}
	Lf := len(format)
	if Lf > 1 && format[Lf-1] == ' ' {
		format = format[:Lf-1]
	}

	t := 1
	uToB := []byte(u)
	switch {
	case ASCIIEndsWithFold(uToB, []byte("kb")):
		t = KB
	case ASCIIEndsWithFold(uToB, []byte("mb")):
		t = MB
	case ASCIIEndsWithFold(uToB, []byte("gb")):
		t = GB
	case ASCIIEndsWithFold(uToB, []byte("tb")):
		t = TB
	case ASCIIEndsWithFold(uToB, []byte("pb")):
		t = PB
	case ASCIIEndsWithFold(uToB, []byte("eb")):
		t = EB
	default:
		t = 1
	}
	t2, cut := 1, 2
	vToB := []byte(v)
	switch {
	case ASCIIEndsWithFold(vToB, []byte("kb")):
		t2 = KB
	case ASCIIEndsWithFold(vToB, []byte("mb")):
		t2 = MB
	case ASCIIEndsWithFold(vToB, []byte("gb")):
		t2 = GB
	case ASCIIEndsWithFold(vToB, []byte("tb")):
		t2 = TB
	case ASCIIEndsWithFold(vToB, []byte("pb")):
		t2 = PB
	case ASCIIEndsWithFold(vToB, []byte("eb")):
		t2 = EB
	case ASCIIEndsWithFold(vToB, []byte("b")):
		cut = 1
	default:
		cut = 0
	}
	v = v[:len(v)-cut]
	vInt := make([]interface{}, fCount)
	for i := range vInt {
		dt, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			vInt[i] = 0
			continue
		}
		if t > t2 {
			vInt[i] = float64(dt) * float64(t2) / float64(t)
			continue
		}
		vInt[i] = dt * uint64(t2/t)
	}
	return fmt.Sprintf(format, vInt...)
}

// Set ...
func (t *V) Set(a string) {
	a = strings.TrimLeft(a, "-")
	v := strings.Split(a, " ")

	// removing empty entries
	La := len(v)
	i := 0
	for i < La {
		v[i] = strings.ReplaceAll(v[i], " ", "")
		if v[i] == "" {
			if (i + 1) == La {
				v = v[:i]
			} else {
				v = append(v[:i], v[i+1:]...)
			}
			La--
			continue
		}
		i++
	}
	*t = v
}

// ASCIIEndsWithFold ...
func ASCIIEndsWithFold(a, b []byte) bool {
	La, Lb := len(a), len(b)
	if La < Lb {
		return false
	}
	const (
		A byte = 0x41 // 'A'
		// Z  byte = 0x5a // 'Z'
		_A byte = 0x61 // 'a'
		_Z byte = 0x7a // 'z'
	)
	for i, v := range a[La-Lb:] {
		n := b[i]
		if v == n {
			continue
		}
		if n >= _A && n <= _Z {
			n -= (_A - A)
		}
		if v >= _A && v <= _Z {
			v -= (_A - A)
		}
		if v != n {
			return false
		}
	}
	return true
}

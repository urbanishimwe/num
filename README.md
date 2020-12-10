# num
[![GoDoc](https://godoc.org/github.com/google/gopacket?status.svg)](https://godoc.org/github.com/urbanishimwe/num)

Num is the CLI to transform integers and characters.

```
go get -u github.com/urbanishimwe/num
```

**FLAGS:**
```
USAGE: flags must be entered before inputs in the command-line

  -b    append output in binary
  -c    treat input as characters and convert them to integers.
  -d    apend output in decimal(default)
  -f string
        name of the file to read inputs from
  -format string
        custom output format with valid printf flags, this override bases flags
  -o    append output in octal
  -s    append output of an integer converted to a character
  -stdin
        allow blocking for inputs from standard input stream
  -u string
        data units for the output i.e B, KB, MB, GB, TB, or PB (default "b")
  -v    verbose: prints parser errors on standard error stream
  -x    append output in hexadecimal
```

**examples**

- converting to GB
```
num -u GB 10TB 8GB
```

- converting to GB and in binary
```
num -u -b 10TB
```

- converting from hexadecimal to decimal
```
num 0x_dad_face
```

- converting from octal with Data unit and custom format
```
num -u TB -f="%gTB" 0X_dad_face_dead_faceGB
```

- UTF8 strings
```
num -c 😍 // output: 128525
```

- UTF8 Code points
```
num -s 128525 // output: 😍
```

- reading input from file with in multiple bases
```
cat input.in | num -x -d -o -u=KB
```

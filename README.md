# num
[![GoDoc](https://godoc.org/github.com/google/gopacket?status.svg)](https://godoc.org/github.com/urbanishimwe/num)

Num is the command-line tool to parse integers from/to different formats.
Integers can be parsed to/from decimal, binary, octal, hexadecimal and byte(b) with data storage units kilobyte(kb),
gigabyte(gb), telabyte(tb), petabyte(pt) and exabyte(eb). custom format also supported

```
go get -u github.com/urbanishimwe/num
```

**FLAGS:**
```
USAGE: (flag must be entered before input in the command line)

  -b    output in binary
  -c    input will be treated as UTF8 characters. it doesn't allow data units and multiple bases, custom format will be applied on every character
  -d    output in decimal(default)
  -f string
        custom output format with valid printf flags, it does not affect data unit but it will override other formats
  -f-count int
        number of flags parsed in custom format(--f). e.g --f '%q %x' --f-count must be 2 (default 1)
  -o    output in octal
  -s    input numbers will be converted to string UTF8
  -stdin
        read input from stdin pipe line by line until EOF
  -u string
        data units for the output i.e KB, MB, GB, TB, PB or EB (default "b")
  -x    output in hexadecimal

you can add multiple bases like -x -d -o -b. the output of every input will be on a single line.
Inputs can be separated with white space or newline. Data unit(flag -u) receive single units,
in converting from Da to Db where Da > Db it is better to use custom format
that receive floating e.g: num -u GB -f="%0.4fGB" 10TB, note that everything is case insensitive.
```

**examples**

- converting to GB
```
num -u GB 10TB 8EB
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

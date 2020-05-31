// Copyright 2020 Urban Ishimwe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestASCIIEndsWithFold(t *testing.T) {
	var c = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	var C = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if !ASCIIEndsWithFold(c, C) {
		t.Error("ASCIIEndsWithFold error")
	}
	c[9] = '+'
	C[9] = '`'
	if ASCIIEndsWithFold(c, C) {
		t.Error("ASCIIEndsWithFold error")
	}
}

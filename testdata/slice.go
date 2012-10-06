// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package main

import "fmt"

var PASS = true

func initialValue() {
	pass := true

	var s1 []byte
	s2 := []byte{}
	s3 := make([]byte, 0)

	tests := []struct {
		msg string
		in  bool
		out bool
	}{
		{"nil s1", s1 == nil, true},
		{"nil s2", s2 == nil, false},
		{"nil s3", s3 == nil, false},
		{"len s1", len(s1) == 0, true},
		{"len s2", len(s2) == 0, true},
		{"len s3", len(s3) == 0, true},
		{"cap s1", cap(s1) == 0, true},
		{"cap s2", cap(s2) == 0, true},
		{"cap s3", cap(s3) == 0, true},
	}

	for _, t := range tests {
		if t.in != t.out {
			fmt.Printf("\tFAIL: %s => got %v, want %v\n", t.msg, t.in, t.out)
			pass, PASS = false, false
		}
	}
	if pass {
		fmt.Println("\tpass")
	}
}

func shortHand() {
	pass := true

	var array = [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	var a_slice, b_slice []byte

	// == 1. Slice of an array

	a_slice = array[4:8]
	if string(a_slice) == "efgh" && len(a_slice) == 4 && cap(a_slice) == 6 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 1. [4:8] => got %v, len=%v, cap=%v\n",
			a_slice, len(a_slice), cap(a_slice))
		pass, PASS = false, false
	}

	a_slice = array[6:7]
	if string(a_slice) != "g" {
		fmt.Printf("\tFAIL: 1. [6:7] => got %v\n", a_slice)
		pass, PASS = false, false
	}

	a_slice = array[:3]
	if string(a_slice) == "abc" && len(a_slice) == 3 && cap(a_slice) == 10 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 1. [:3] => got %v, len=%v, cap=%v\n",
			a_slice, len(a_slice), cap(a_slice))
		pass, PASS = false, false
	}

	a_slice = array[5:]
	if string(a_slice) != "fghij" {
		fmt.Printf("\tFAIL: 1. [5:] => got %v\n", a_slice)
		pass, PASS = false, false
	}

	a_slice = array[:]
	if string(a_slice) != "abcdefghij" {
		fmt.Printf("\tFAIL: 1. [:] => got %v\n", a_slice)
		pass, PASS = false, false
	}

	a_slice = array[3:7]
	if string(a_slice) == "defg" && len(a_slice) == 4 && cap(a_slice) == 7 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 1. [3:7] => got %v, len=%v, cap=%v\n",
			a_slice, len(a_slice), cap(a_slice))
		pass, PASS = false, false
	}

	// == 2. Slice of a slice

	b_slice = a_slice[1:3]
	if string(b_slice) == "ef" && len(b_slice) == 2 && cap(b_slice) == 6 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 2. [1:3] => got %v, len=%v, cap=%v\n",
			b_slice, len(b_slice), cap(b_slice))
		pass, PASS = false, false
	}

	b_slice = a_slice[:3]
	if string(b_slice) != "def" {
		fmt.Printf("\tFAIL: 2. [:3] => got %v\n", b_slice)
		pass, PASS = false, false
	}

	b_slice = a_slice[:]
	if string(b_slice) != "defg" {
		fmt.Printf("\tFAIL: 2. [:] => got %v\n", b_slice)
		pass, PASS = false, false
	}

	// * * *

	if pass {
		fmt.Println("\tpass")
	}
}

func useFunc() {
	pass := true

	// Returns the biggest value in a slice of ints.
	Max := func(slice []int) int {
		max := slice[0] // The first element is the max for now.
		for index := 1; index < len(slice); index++ {
			if slice[index] > max {
				max = slice[index]
			}
		}
		return max
	}

	A1 := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	A2 := [4]int{1, 2, 3, 4}
	A3 := [1]int{1}

	var slice []int

	slice = A1[:] // Take all A1 elements.
	if Max(slice) != 9 {
		fmt.Printf("\tFAIL: A1 => got %v, want 9\n", Max(slice))
		pass, PASS = false, false
	}

	slice = A2[:]
	if Max(slice) != 4 {
		fmt.Printf("\tFAIL: A2 => got %v, want 4\n", Max(slice))
		pass, PASS = false, false
	}

	slice = A3[:]
	if Max(slice) != 1 {
		fmt.Printf("\tFAIL: A3 => got %v, want 1\n", Max(slice))
		pass, PASS = false, false
	}

	if pass {
		fmt.Println("\tpass")
	}
}

func reference() {
	pass := true

	fmtSlice := func(slice []byte) string {
		s := ("[")
		for index := 0; index < len(slice)-1; index++ {
			s += fmt.Sprintf("%q,", slice[index])
		}
		s += fmt.Sprintf("%q]", slice[len(slice)-1])

		return s
	}

	A := [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	slice1 := A[3:7]
	slice2 := A[5:]
	slice3 := slice1[:2]

	// == 1. Current content of A and the slices.

	tests := []struct {
		msg string
		in  string
		out string
	}{
		{"A", fmtSlice(A[:]), "['a','b','c','d','e','f','g','h','i','j']"},
		{"slice1", fmtSlice(slice1), "['d','e','f','g']"},
		{"slice2", fmtSlice(slice2), "['f','g','h','i','j']"},
		{"slice3", fmtSlice(slice3), "['d','e']"},
	}

	for _, t := range tests {
		if t.in != t.out {
			fmt.Printf("\tFAIL: 1. %s => got %v, want %v\n", t.msg, t.in, t.out)
			pass, PASS = false, false
		}
	}

	// == 2. Let's change the 'e' in A to 'E'.
	A[4] = 'E'

	tests = []struct {
		msg string
		in  string
		out string
	}{
		{"A", fmtSlice(A[:]), "['a','b','c','d','E','f','g','h','i','j']"},
		{"slice1", fmtSlice(slice1), "['d','E','f','g']"},
		{"slice2", fmtSlice(slice2), "['f','g','h','i','j']"},
		{"slice3", fmtSlice(slice3), "['d','E']"},
	}

	for _, t := range tests {
		if t.in != t.out {
			fmt.Printf("\tFAIL: 2. %s => got %v, want %v\n", t.msg, t.in, t.out)
			pass, PASS = false, false
		}
	}

	// == 3. Let's change the 'g' in slice2 to 'G'.
	slice2[1] = 'G'

	tests = []struct {
		msg string
		in  string
		out string
	}{
		{"A", fmtSlice(A[:]), "['a','b','c','d','E','f','G','h','i','j']"},
		{"slice1", fmtSlice(slice1), "['d','E','f','G']"},
		{"slice2", fmtSlice(slice2), "['f','G','h','i','j']"},
		{"slice3", fmtSlice(slice3), "['d','E']"},
	}

	for _, t := range tests {
		if t.in != t.out {
			fmt.Printf("\tFAIL: 3. %s => got %v, want %v\n", t.msg, t.in, t.out)
			pass, PASS = false, false
		}
	}

	// * * *

	if pass {
		fmt.Println("\tpass")
	}
}

func resize() {
	pass := true

	var slice []byte

	// == 1.
	slice = make([]byte, 4, 5) // [0 0 0 0]

	if len(slice) == 4 && cap(slice) == 5 &&
		slice[0] == 0 && slice[1] == 0 && slice[2] == 0 && slice[3] == 0 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 1. got %v, want [0 0 0 0])\n", slice)
		pass, PASS = false, false
	}

	// == 2.
	slice[1], slice[3] = 2, 3 // [0 2 0 3]

	if slice[0] == 0 && slice[1] == 2 && slice[2] == 0 && slice[3] == 3 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 2. got %v, want [0 2 0 3])\n", slice)
		pass, PASS = false, false
	}

	// == 3.
	slice = make([]byte, 2) // Resize: [0 0]

	if len(slice) == 2 && cap(slice) == 2 && slice[0] == 0 && slice[1] == 0 {
		// ok
	} else {
		fmt.Printf("\tFAIL: 3. got %v, want [0 0])\n", slice)
		pass, PASS = false, false
	}
	//==

	if pass {
		fmt.Println("\tpass")
	}
}

func main() {
	fmt.Print("\n\n== Slices\n")

	fmt.Print("\n=== RUN initialValue\n")
	initialValue()
	fmt.Print("\n=== RUN shortHand\n")
	shortHand()
	fmt.Print("\n=== RUN useFunc\n")
	useFunc()
	fmt.Print("\n=== RUN reference\n")
	reference()
	fmt.Print("\n=== RUN resize\n")
	resize()

	if PASS {
		fmt.Print("\nPASS\n")
	} else {
		fmt.Print("\nFAIL\n")
		print("Fail: Slices")
	}
}

/*
== reference()

=== First content of A and the slices
A is : ['a','b','c','d','e','f','g','h','i','j']
slice1 is : ['d','e','f','g']
slice2 is : ['f','g','h','i','j']
slice3 is : ['d','e']

=== Content of A and the slices, after changing 'e' to 'E' in array A
A is : ['a','b','c','d','E','f','g','h','i','j']
slice1 is : ['d','E','f','g']
slice2 is : ['f','g','h','i','j']
slice3 is : ['d','E']

=== Content of A and the slices, after changing 'g' to 'G' in slice2
A is : ['a','b','c','d','E','f','G','h','i','j']
slice1 is : ['d','E','f','G']
slice2 is : ['f','G','h','i','j']
slice3 is : ['d','E']

*/
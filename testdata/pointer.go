// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package main

import "fmt"

var PASS = true

// Global declaration of a pointer
var i int
var hello string
var p *int

func init() {
	p = &i             // p points to i (p stores the address of i)
	helloPtr := &hello // pointer variable of type *string which points to hello

	fmt.Println("== init()")
	fmt.Print("\t\"helloPtr\": ", helloPtr)
}

func declaration() {
	var i int
	var hello string
	var p *int

	p = &i
	helloPtr := &hello
	fmt.Println("\t\"p\":", p, "\n\t\"helloPtr\":", helloPtr)
}

func showAddress() {
	var (
		i     int     = 9
		hello string  = "Hello world"
		pi    float32 = 3.14
		b     bool    = true
	)

	fmt.Println("\t\"i\":", &i)
	fmt.Println("\t\"hello\":", &hello)
	fmt.Println("\t\"pi\":", &pi)
	fmt.Println("\t\"b\":", &b)
}

func nilValue() {
	pass := true

	var num = 10
	var p *int

	if p == nil {
		// ok
	} else {
		fmt.Printf("\tFAIL: declaration => got %v\n", p == nil)
		pass, PASS = false, false
	}

	p = &num
	if p != nil {
		// ok
	} else {
		fmt.Printf("\tFAIL: assignment => got %v\n", p == nil)
		pass, PASS = false, false
	}

	if pass {
		fmt.Println("\tpass")
	}
}

func access() {
	pass := true

	hello := "Hello, mina-san!"
	var helloPtr *string
	helloPtr = &hello

	i := 6
	iPtr := &i

	if *helloPtr != "Hello, mina-san!" {
		fmt.Printf("\tFAIL: *helloPtr => got %v, want %v\n", *helloPtr, hello)
		pass, PASS = false, false
	}
	if *iPtr != 6 {
		fmt.Printf("\tFAIL: *iPtr => got %v, want %v\n", *iPtr, i)
		pass, PASS = false, false
	}

	// * * *

	x := 3
	y := &x

	*y++
	if x != 4 {
		fmt.Printf("\tFAIL: x => got %v, want 4\n", x)
		pass, PASS = false, false
	}

	*y++
	if x != 5 {
		fmt.Printf("\tFAIL: x => got %v, want 5\n", x)
		pass, PASS = false, false
	}

	if pass {
		fmt.Println("\tpass")
	}
}

func allocation() {
	sum := 0
	var doubleSum *int // a pointer to int
	for i := 0; i < 10; i++ {
		sum += i
	}

	doubleSum = new(int) // allocate memory for an int and make doubleSum point to it
	*doubleSum = sum * 2 // use the allocated memory, by dereferencing doubleSum

	if sum == 45 && *doubleSum == 90 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: sum=%v, *doubleSum=%v\n", sum, *doubleSum)
		PASS = false
	}
}

func parameterByValue() {
	// Returns 1 plus its input parameter.
	var add = func(v int) int {
		v = v + 1
		return v
	}

	x := 3
	x1 := add(x)

	if x == 3 && x1 == 4 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: x=%v, x1=%v\n", x, x1)
		PASS = false
	}
}

func byReference_1() {
	add := func(v *int) int { // pointer to int
		*v = *v + 1 // we dereference and change the value pointed by a
		return *v
	}

	x := 3
	x1 := add(&x) // by passing the adress of x to it

	if x1 == 4 && x == 4 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: x=%v, x1=%v\n", x, x1)
		PASS = false
	}

	x1 = add(&x)
	if x == 5 && x1 == 5 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: x=%v, x1=%v\n", x, x1)
		PASS = false
	}
}

func byReference_2() {
	add := func(v *int, i int) { *v += i }
	value := 6
	incr := 1

	add(&value, incr)
	if value == 7 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: value=%v\n", value)
		PASS = false
	}

	add(&value, incr)
	if value == 8 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: value=%v\n", value)
		PASS = false
	}
}

func byReference_3() {
	x := 3
	f := func() {
		x = 4
	}
	y := &x

	f()
	if *y == 4 {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: 3. *y=%v\n", *y)
		PASS = false
	}
}

func main() {
	fmt.Print("\n\n== Pointers\n")

	fmt.Print("\n=== RUN declaration\n")
	declaration()
	fmt.Print("\n=== RUN showAddress\n")
	showAddress()

	fmt.Print("\n=== RUN nilValue\n")
	nilValue()
	fmt.Print("\n=== RUN access\n")
	access()
	fmt.Print("\n=== RUN allocation\n")
	allocation()

	fmt.Print("\n=== RUN parameterByValue\n")
	parameterByValue()
	fmt.Print("\n=== RUN byReference_1\n")
	byReference_1()
	fmt.Print("\n=== RUN byReference_2\n")
	byReference_2()
	fmt.Print("\n=== RUN byReference_3\n")
	byReference_3()

	if PASS {
		fmt.Print("\nPASS\n")
	} else {
		fmt.Print("\nFAIL\n")
		print("Fail: Pointers")
	}
}

/*
== init()
	"helloPtr": 0x4e02b8

== Pointers

=== RUN declaration
	"p": 0xf840038018 
	"helloPtr": 0xf840028070

=== RUN showAddress
	"i": 0xf840038020
	"hello": 0xf840028030
	"pi": 0xf840038028
	"b": 0xf840038030
*/
// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package main

import "fmt"

var PASS = true

type person struct {
	name string
	age  int
}

// Return the older person of p1 and p2, and the difference in their ages.
func older(p1, p2 person) (person, int) {
	if p1.age > p2.age {
		return p1, p1.age - p2.age
	}
	return p2, p2.age - p1.age
}

// Return the older person in a group of 10 persons.
func older10(people [10]person) person {
	older := people[0] // The first one is the older for now.

	// Loop through the array and check if we could find an older person.
	for index := 1; index < 10; index++ { // We skipped the first element here.
		if people[index].age > older.age {
			older = people[index]
		}
	}
	return older
}

// == Array
//

func zeroArray() {
	pass := true

	var a1 [4]byte
	a2 := [4]byte{}

	tests := []struct {
		msg string
		in  bool
		out bool
	}{
		//{"nil a1", a1 == nil, true},
		//{"nil a2", a2 == nil, false},
		{"len a1", len(a1) == 40, true},
		{"len a2", len(a2) == 4, true},
		{"cap a1", cap(a1) == 4, true},
		{"cap a2", cap(a2) == 4, true},
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

func initArray() {
	pass := true

	// Declare and initialize an array A of 10 person.
	array1 := [10]person{
		person{"", 0},
		person{"Paul", 23},
		person{"Jim", 24},
		person{"Sam", 84},
		person{"Rob", 54},
		person{"", 0},
		person{"", 0},
		person{"", 0},
		person{"Karl", 10},
		person{"", 0},
	}

	// Declare and initialize an array of 10 persons, but let the compiler guess the size.
	array2 := [...]person{ // Substitute '...' instead of an integer size.
		person{"", 0},
		person{"Paul", 23},
		person{"Jim", 24},
		person{"Sam", 84},
		person{"Rob", 54},
		person{"", 0},
		person{"", 0},
		person{"", 0},
		person{"Karl", 10},
		person{"", 0}}

	tests := []struct {
		msg string
		in  bool
		out bool
	}{
		{"len", len(array1) == len(array2), true},
		{"cap", cap(array1) == cap(array2), true},
		{"equality", array1 == array2, true},
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

func testArray() {
	// Declare an example array variable of 10 person called 'array'.
	var array [10]person

	// Initialize some of the elements of the array, the others are by default
	// set to person{"", 0}
	array[1] = person{"Paul", 23}
	array[2] = person{"Jim", 24}
	array[3] = person{"Sam", 84}
	array[4] = person{"Rob", 54}
	array[8] = person{"Karl", 19}

	older := older10(array) // Call the function by passing it our array.

	if older.name == "Sam" {
		fmt.Println("\tpass")
	} else {
		fmt.Printf("\tFAIL: got %v, want Sam\n", older.name)
		PASS = false
	}
}

func multiArray() {
	// Declare and initialize an array of 2 arrays of 4 ints
	doubleArray_1 := [2][4]int{[4]int{1, 2, 3, 4}, [4]int{5, 6, 7, 8}}

	// Simplify the previous declaration, with the '...' syntax
	doubleArray_2 := [2][4]int{
		[...]int{1, 2, 3, 4}, [...]int{5, 6, 7, 8}}

	// Super simpification!
	doubleArray_3 := [2][4]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}

	if doubleArray_1 == doubleArray_2 && doubleArray_2 == doubleArray_3 {
		fmt.Println("\tpass")
	} else {
		fmt.Print("\tFAIL: got different arraies\n")
		PASS = false
	}
}

// == Struct
//

func testStruct() {
	pass := true

	var tom person
	tom.name, tom.age = "Tom", 18

	bob := person{age: 25, name: "Bob"} // specify the fields and their values
	paul := person{"Paul", 43}          // specify values of fields in their order

	TB_older, TB_diff := older(tom, bob)
	TP_older, TP_diff := older(tom, paul)
	BP_older, BP_diff := older(bob, paul)

	tests := []struct {
		msg       string
		inPerson  person
		outPerson person
		inDiff    int
		outDiff   int
	}{
		{"Tom,Bob", TB_older, bob, TB_diff, 7},
		{"Tom,Paul", TP_older, paul, TP_diff, 25},
		{"Bob,Paul", BP_older, paul, BP_diff, 18},
	}

	for _, t := range tests {
		if t.inPerson != t.outPerson {
			fmt.Printf("\tFAIL: %s => person got %v, want %v\n",
				t.msg, t.inPerson, t.outPerson)
			pass, PASS = false, false
		}
		if t.inDiff != t.outDiff {
			fmt.Printf("\tFAIL: %s => difference got %v, want %v\n",
				t.msg, t.inDiff, t.outDiff)
			pass, PASS = false, false
		}
	}
	if pass {
		fmt.Println("\tpass")
	}
}

func main() {
	fmt.Print("\n\n== Composite types\n")

	fmt.Print("\n=== RUN zeroArray\n")
	zeroArray()
	fmt.Print("\n=== RUN initArray\n")
	initArray()
	fmt.Print("\n=== RUN testArray\n")
	testArray()
	fmt.Print("\n=== RUN multiArray\n")
	multiArray()

	fmt.Print("\n=== RUN testStruct\n")
	testStruct()

	if PASS {
		fmt.Print("\nPASS\n")
	} else {
		fmt.Print("\nFAIL\n")
		print("Fail: Composite types")
	}
}
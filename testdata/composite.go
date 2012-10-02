// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package main

import "fmt"

// We declare our new type
type person struct {
	name string
	age  int
}

// Return the older person of p1 and p2, and the difference in their ages.
func Older(p1, p2 person) (person, int) {
	if p1.age > p2.age { // Compare p1 and p2's ages
		return p1, p1.age - p2.age
	}
	return p2, p2.age - p1.age
}

func testStruct() {
	var tom person

	tom.name, tom.age = "Tom", 18

	// Look how to declare and initialize easily.
	bob := person{age: 25, name: "Bob"} //specify the fields and their values
	paul := person{"Paul", 43}          //specify values of fields in their order

	tb_Older, tb_diff := Older(tom, bob)
	// Checking
	if tb_Older == bob && tb_diff == 7 {
		println("[OK] Tom, Bob")
	} else {
		fmt.Printf("[Error] Of %s and %s, %s is older by %d years\n",
			tom.name, bob.name, tb_Older.name, tb_diff)
	}
	//==

	tp_Older, tp_diff := Older(tom, paul)
	// Checking
	if tp_Older == paul && tp_diff == 25 {
		println("[OK] Tom, Paul")
	} else {
		fmt.Printf("[Error] Of %s and %s, %s is older by %d years\n",
			tom.name, paul.name, tp_Older.name, tp_diff)
	}
	//==

	bp_Older, bp_diff := Older(bob, paul)
	// Checking
	if bp_Older == paul && bp_diff == 18 {
		println("[OK] Bob, Paul")
	} else {
		fmt.Printf("[Error] Of %s and %s, %s is older by %d years\n",
			bob.name, paul.name, bp_Older.name, bp_diff)
	}
}

// * * *

// Return the older person in a group of 10 persons.
func Older10(people [10]person) person {
	older := people[0] // The first one is the older for now.

	// Loop through the array and check if we could find an older person.
	for index := 1; index < 10; index++ { // We skipped the first element here.
		if people[index].age > older.age { // Current's persons age vs olderest so far.
			older = people[index] // If people[index] is older, replace the value of older.
		}
	}
	return older
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

	older := Older10(array) // Call the function by passing it our array.

	// Checking
	if older.name == "Sam" {
		println("[OK]")
	} else {
		fmt.Printf("[Error] The older of the group is: %s\n", older.name)
	}
}

// * * *

func initializeArray() {
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

	// Checking
	if len(array1) == len(array2) {
		println("[OK] length")
	} else {
		fmt.Printf("[Error] len => array1: %d, array2: %d\n", len(array1), len(array2))
	}

	if array1 == array2 {
		println("[OK] comparison")
	} else {
		fmt.Printf("[Error] array1: %v\narray2: %v\n", array1, array2)
	}
}

// * * *

func multiArray() {
	// declare and initialize an array of 2 arrays of 4 ints
	doubleArray_1 := [2][4]int{[4]int{1, 2, 3, 4}, [4]int{5, 6, 7, 8}}

	// simplify the previous declaration, with the '...' syntax
	doubleArray_2 := [2][4]int{
		[...]int{1, 2, 3, 4}, [...]int{5, 6, 7, 8}}

	// super simpification!
	doubleArray_3 := [2][4]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}

	// Checking
	if doubleArray_1 == doubleArray_2 && doubleArray_2 == doubleArray_3 {
		println("[OK]")
	} else {
		fmt.Println("[Error] multi-dimensional")
	}
}

// * * *

func main() {
	println("\n== testStruct")
	testStruct()
	println("\n== testArray")
	testArray()
	println("\n== initializeArray")
	initializeArray()
	println("\n== multiArray")
	multiArray()
}

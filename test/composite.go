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
	tp_Older, tp_diff := Older(tom, paul)
	bp_Older, bp_diff := Older(bob, paul)

	fmt.Printf("Of %s and %s, %s is older by %d years\n",
		tom.name, bob.name, tb_Older.name, tb_diff)

	fmt.Printf("Of %s and %s, %s is older by %d years\n",
		tom.name, paul.name, tp_Older.name, tp_diff)

	fmt.Printf("Of %s and %s, %s is older by %d years\n",
		bob.name, paul.name, bp_Older.name, bp_diff)
}

func main() {
	println("\n== testStruct()\n")
	testStruct()

	
}
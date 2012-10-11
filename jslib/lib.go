// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

// Package g handles the features and Go types in JavaScript.

package g

// The specific type that it represents.
const (
	invalidT int = iota
	arrayT
	mapT
	sliceT
)

// == Array
//

func init() {
	// Use the toString() method when Array.isArray isn't implemented:
	// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Array/isArray#Compatibility
	if !Array.isArray {
		Array.isArray = func(arg interface{}) {
			return Object.prototype.toString.call(arg) == "[object Array]"
		}
	}
}

// The array can not be compared with nil.
// The capacity is the same than length.

// arrayType represents a fixed array type.
type arrayType struct {
	v []interface{} // array's value

	len_ map[int]int
}

// len returns the length for the given dimension.
func (a arrayType) len(dim int) int {
	if dim == nil {
		return a.len_[0]
	}
	return a.len_[len(arguments)]
}

// cap returns the capacity for the given dimension.
func (a arrayType) cap(dim int) int {
	if dim == nil {
		return a.len_[0]
	}
	return a.len_[len(arguments)]
}

// str returns the array (of bytes or runes) like a string.
func (a arrayType) str() string {
	return a.v.join("")
}

// typ returns the type.
func (a arrayType) typ() int { return arrayT }

// MkArray initializes an array of dimension "dim" to value "zero",
// merging the elements of "data" if any.
func MkArray(dim []int, zero interface{}, data []interface{}) *arrayType {
	a := new(arrayType)

	if data != nil {
		if !equalDim(dim, getDimArray(data)) {
			a.v = initArray(dim, zero)
			mergeArray(a.v, data)
		} else {
			a.v = data
		}
	} else {
		a.v = initArray(dim, zero)
	}

	for i, v := range dim {
		a.len_[i] = v
	}

	return a
}

// * * *

// equalDim reports whether d1 and d2 are equal.
func equalDim(d1, d2 []int) bool {
	if len(d1) != len(d2) {
		return false
	}
	for i, v := range d1 {
		if v != d2[i] {
			return false
		}
	}
	return true
}

// getDimArray returns the dimension of an array.
func getDimArray(a []interface{}) (dim []int) {
	for {
		dim.push(len(a))

		if Array.isArray(a[0]) {
			a = a[0]
		} else {
			break
		}
	}
	return
}

// initArray returns an array of dimension given in "dim" initialized to "zero".
func initArray(dim []int, zero interface{}) (a []interface{}) {
	if len(dim) == 0 {
		return zero
	}
	nextArray := initArray(dim.slice(1), zero)

	for i := 0; i < dim[0]; i++ {
		a[i] = nextArray
	}
	return
}

// mergeArray merges src in array dst.
func mergeArray(dst, src []interface{}) {
	for i, srcVal := range src {
		if Array.isArray(srcVal) {
			mergeArray(dst[i], srcVal)
		} else {
			isHashMap := false

			// The position is into a hash map, if any
			if typeof(srcVal) == "object" {
				for k, v := range srcVal {
					if srcVal.hasOwnProperty(k) { // identify a hashmap
						isHashMap = true
						i = k
						dst[i] = v
					}
				}
			}
			if !isHashMap {
				dst[i] = srcVal
			}
		}
	}
}

// == Slice
//

// sliceType represents a slice type.
type sliceType struct {
	arr interface{}   // the array where this slice is got, if any
	v   []interface{} // elements from scratch (make) or appended to the array

	low  int // indexes for the array
	high int
	len  int // total of elements
	cap  int

	nil_ bool // for variables declared like slices
}

// typ returns the type.
func (s sliceType) typ() int { return sliceT }

func (s sliceType) isNil() bool {
	if s.len != 0 || s.cap != 0 {
		return false
	}
	return s.nil_
}

// MkSlice initializes a slice with the zero value.
func MkSlice(zero interface{}, len, cap int) *sliceType {
	s := new(sliceType)

	if zero == nil {
		s.nil_ = true
		return s
	}

	s.len = len

	for i := 0; i < len; i++ {
		s.v[i] = zero
	}
	if cap != nil {
		s.cap = cap
	} else {
		s.cap = len
	}
	return s
}

// Slice creates a new slice with the elements in "data".
func Slice(zero interface{}, data []interface{}) *sliceType {
	s := new(sliceType)

	if len(arguments) == 0 {
		s.nil_ = true
		return s
	}

	for i, srcVal := range data {
		isHashMap := false

		// The position is into a hash map, if any
		if typeof(srcVal) == "object" {
			for k, v := range srcVal {
				if srcVal.hasOwnProperty(k) { // identify a hashmap
					isHashMap = true

					for i; i < k; i++ {
						s.v[i] = zero
					}
					s.v[i] = v
				}
			}
		}
		if !isHashMap {
			s.v[i] = srcVal
		}
	}

	s.len = len(s.v)
	s.cap = s.len
	return s
}

// SliceFrom creates a new slice from an array using the indexes low and high.
func SliceFrom(src interface{}, low, high int) *sliceType {
	s := new(sliceType)
	s.set(src, low, high)

	/*if src.low != nil { // slice
		
	} else { // array
		s.arr = src
		s.low = low
		s.high = high
		s.len = high - low
		s.cap = src.cap() - low
	}*/

	return s
}

// set sets the elements of a slice.
func (s sliceType) set(src interface{}, low, high int) {
	s.low, s.high = low, high

	if src.low != nil { // slice
		if len(src.v) != 0 {
			s.v = src.v.slice(low, high)
		} else {
			s.v = src.arr.v.slice(low, high)
		}
		s.len = len(s.v)
		s.cap = src.cap - low
	} else { // array
		s.arr = src
		s.v = src.v.slice(low, high)
		s.len = high - low
		s.cap = src.cap() - low
	}

	/*if src.v != nil { // from make
		s.v = src.v.slice(low, high)
		s.cap = src.cap - low
		s.len = len(s.v)

	} else { // from array
		s.arr = i
		s.cap = len(i) - low
		s.len = high - low
	}*/
}

// get gets the slice.
/*func (s sliceType) get() []interface{} {
	if len(s.v) != 0 {
		return s.v
	}
	//a := s.arr
	return s.arr.v.slice(s.low, s.high)
}*/

// str returns the slice (of bytes or runes) like a string.
func (s sliceType) str() string {
	/*_s := s.get()
	return _s.join("")*/
	return s.v.join("")
}

/*
func (s sliceType) setSlice(src interface{}, low, high int) {
	s.low, s.high = low, high
	s.len = high - low

	if src.arr != nil { // from slice
		s.arr = src.arr
		s.cap = src.cap - low
	} else { // from array
		s.arr = i
		s.cap = len(i) - low
	}
}

// Appends an element to the slice.
func (s sliceType) append(elt interface{}) {
	if s.len == s.cap {
		s.cap = s.len * 2
	}
	s.len++
}
*/

// == Map
//

// The length into a map is rarely used so, in JavaScript, I prefer to calculate
// the length instead of use a field.
//
// A map has not built-in function "cap".

// mapType represents a map type.
// The compiler adds the appropriate zero value for the map (which it is work out
// from the map type).
type mapType struct {
	v    map[interface{}]interface{} // map's value
	zero interface{}                 // zero value for the map's value
}

// len returns the number of elements.
func (m mapType) len() int {
	len := 0
	for key, _ := range m.v {
		if m.v.hasOwnProperty(key) {
			len++
		}
	}
	return len
}

// typ returns the type.
func (m mapType) typ() int { return mapT }

// Map creates a map storing its zero value.
func Map(zero interface{}, v map[interface{}]interface{}) *mapType {
	m := &mapType{v, zero}
	return m
}

// get returns the value for the key "k" if it exists and a boolean indicating it.
// If looking some key up in M's map gets you "nil" ("undefined" in JS),
// then return a copy of the zero value.
func (m mapType) get(k interface{}) (interface{}, bool) {
	v := m.v

	// Allow multi-dimensional index (separated by commas)
	for i := 0; i < len(arguments); i++ {
		v = v[arguments[i]]
	}

	if v == nil {
		return m.zero, false
	}
	return v, true
}

// == Utility
//

// Export adds public names from "exported" to the map "pkg".
func Export(pkg map[interface{}]interface{}, exported []interface{}) {
	for _, v := range exported {
		pkg.v = v
	}
}

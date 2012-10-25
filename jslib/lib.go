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
		Array.isArray = func(arg interface{}) bool {
			return Object.prototype.toString.call(arg) == "[object Array]"
		}
	}
}

// The array can not be compared with nil.
// The capacity is the same than length.

// arrayType represents a fixed array type.
type arrayType struct {
	v     []interface{} // array's value
	refer []interface{} // references from slices

	len_ map[int]int
}

// len returns the length for the given dimension.
func (a arrayType) len(index int) int {
	if index == nil {
		return a.len_[0]
	}
	return a.len_[len(arguments)]
}

// cap returns the capacity for the given dimension.
func (a arrayType) cap(index int) int {
	if index == nil {
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

// MkArray initializes an array of dimension "index" to value "zero",
// merging the elements of "data" if any.
func MkArray(index []int, zero interface{}, data []interface{}) *arrayType {
	a := new(arrayType)

	if data != nil {
		if !equalIndex(index, indexArray(data)) {
			a.v = initArray(index, zero)
			mergeArray(a.v, data)
		} else {
			a.v = data
		}
	} else {
		a.v = initArray(index, zero)
	}

	for i, v := range index {
		a.len_[i] = v
	}

	return a
}

// set sets the value v in the index, and updates the slices that are
// referencing to this array (if any).
func (a arrayType) set(index []int, v interface{}) {
	for i := 0; i < len(index)-1; i++ {
		a.v = a.v[index[i]]
	}
	a.v[index[i]] = v

	for _, r := range a.refer {
		r.v[index[i]] = v
	}
}

// * * *

// equalIndex reports whether index1 and index2 are equal.
func equalIndex(index1, index2 []int) bool {
	if len(index1) != len(index2) {
		return false
	}
	for i, v := range index1 {
		if v != index2[i] {
			return false
		}
	}
	return true
}

// indexArray returns the dimension of an array.
func indexArray(a []interface{}) (index []int) {
	for {
		index.push(len(a))

		if Array.isArray(a[0]) {
			a = a[0]
		} else {
			break
		}
	}
	return
}

// initArray returns an array of dimension given in "index" initialized to "zero".
func initArray(index []int, zero interface{}) (a []interface{}) {
	if len(index) == 0 {
		return zero
	}
	nextArray := initArray(index.slice(1), zero)

	for i := 0; i < index[0]; i++ {
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
	v     []interface{} // slice's value
	refer []interface{} // references from other slices

	low  int // indexes for the array
	high int
	len  int // total of elements
	cap  int

	nil_ bool // for variables declared like slices
}

func (s sliceType) isNil() bool {
	if s.len != 0 || s.cap != 0 {
		return false
	}
	return s.nil_
}

// str returns the slice (of bytes or runes) like a string.
func (s sliceType) str() string {
	return s.v.join("")
}

// typ returns the type.
func (s sliceType) typ() int { return sliceT }

// MkSlice initializes a slice with the zero value.
func MkSlice(zero interface{}, len, cap int) *sliceType {
	s := new(sliceType)

	if zero == nil {
		s.nil_ = true
		return s
	}

	// The fastest way of fill in an array is when array length is specified first.
	s.v = Array(len)
	for i := 0; i < len; i++ {
		s.v[i] = zero
	}

	if cap != nil {
		s.cap = cap
	} else {
		s.cap = len
	}

	s.len = len
	return s
}

// Slice creates a new slice with the elements in "data".
func Slice(zero interface{}, data []interface{}) *sliceType {
	s := new(sliceType)

	if zero == nil {
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

// SliceFrom creates a new slice from an array or slice using the indexes low and high.
func SliceFrom(src interface{}, low, high int) *sliceType {
	s := new(sliceType)

	if low != nil {
		s.low = low
	} else {
		s.low = 0
	}
	if high != nil {
		s.high = high
	} else {
		if src.arr != nil { // slice
			s.high = src.len
		} else {
			s.high = len(src.v)
		}
	}

	s.len = s.high - s.low

	if src.nil_ != nil { // slice
		s.cap = src.cap - s.low
//		s.low += src.low
//		s.high += src.low
	} else { // array
		s.cap = src.cap() - s.low
	}

	s.v = src.v.slice(s.low, s.high)
//	s.v = src.v.slice(low, high)

	src.refer.push(s)
	return s
}

// set sets the value v in the index, and updates the slices that are
// referencing to this slice (if any).
func (s sliceType) set(index []int, v interface{}) {
	for i := 0; i < len(index)-1; i++ {
		s.v = s.v[index[i]]
	}
	s.v[index[i]] = v

	for _, r := range s.refer {
		r.v[index[i]] = v
	}
}

// * * *

// Append implements the function "append".
func Append(src []interface{}, elt ...interface{}) (dst sliceType) {
	// Copy src to the new slice
	dst.low = src.low
	dst.high = src.high
	dst.len = src.len
	dst.cap = src.cap
	dst.nil_ = src.nil_

	dst.v = Array(src.len)
	for i, v := range src.v {
		dst.v[i] = v
	}
	//==

	// TODO: handle len() in interfaces
	// lastIdxElt := len(elt) - 1

	for _, v := range elt {
		if /*i == lastIdxElt &&*/ Array.isArray(v) { // The last field could be an ellipsis
			for _, vArr := range v {
				dst.v.push(vArr)
				if dst.len == dst.cap {
					dst.cap = dst.len * 2
				}
				dst.len++
			}
			break
		}

		dst.v.push(v)
		if dst.len == dst.cap {
			dst.cap = dst.len * 2
		}
		dst.len++
	}
	return dst
}

// Copy implements the function "copy".
func Copy(dst []interface{}, src interface{}) (n int) {
	// []T <= []T
	if src.nil_ != nil {
		for ; n < src.len; n++ {
			if n == dst.len {
				return
			}
			dst.v[n] = src.v[n]
		}
		// TODO: copy refer?
		/*if len(src.refer) != 0 {
			for _, v := range src.refer {
				dst.refer.push(v)
			}
		}*/
		return
	}

	// []byte <= string
	for ; n < len(src); n++ {
		if n == dst.len {
			break
		}
		dst.v[n] = src[n]
	}
	return
}

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

/*func Len(v interface{}) {
	
}*/

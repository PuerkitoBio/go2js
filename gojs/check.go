// Copyright 2011  The "GoJscript" Authors
//
// Use of this source code is governed by the BSD 2-Clause License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package gojs

import (
	"fmt"
	"go/ast"
	"go/token"
)

/*var types = []string{
	"bool", "string",

	// Numeric types
	"uint8", "uint16", "uint32", "uint64",
	"int8", "int16", "int32", "int64",
	"float32", "float64",
	"complex64", "complex128",
	"byte", "rune", "uint", "int", "uintptr",
}

// Checks if the literal is a type.
func isType(tok token.Token, lit string) bool {
	if tok != token.IDENT {
		return false
	}

	for _, v := range types {
		if v == lit {
			return true
		}
	}
	return false
}*/

type check struct {
	isCallExpr     bool
	isCompositeLit bool
	isNegative     bool
	fset           *token.FileSet
}

// Initializes a new type of "check".
func newCheck(fset *token.FileSet) *check {
	return &check{fset: fset}
}

// Returns a general Position.
func (c *check) Position(expr ast.Expr) token.Position {
	return c.fset.Position(expr.Pos())
}

// Checks if it has a valid type for JavaScript.
func (c *check) Type(expr ast.Expr) error {
	switch typ := expr.(type) {
	default:
		panic(fmt.Sprintf("unimplemented: %T", typ))

	case *ast.ArrayType:
		return c.Type(typ.Elt)

	case *ast.BasicLit:

	case *ast.BinaryExpr:
		if err := c.Type(typ.X); err != nil {
			return err
		}
		if err := c.Type(typ.Y); err != nil {
			return err
		}

	case *ast.CallExpr:
		c.isCallExpr = true
		ident := typ.Fun.(*ast.Ident).Name

		switch ident {
		case "make", "new":
			return c.Type(typ.Args[0])

		case "int64", "uint64":
			return fmt.Errorf("%s: conversion of type %s", c.Position(typ), ident)

		// golang.org/pkg/builtin/
		case "complex":
			return fmt.Errorf("%s: built-in function %s()", c.Position(typ), ident)
		}

	// http://golang.org/pkg/go/ast/#ChanType || godoc go/ast ChanType
	case *ast.ChanType:
		return fmt.Errorf("%s: channel type", c.Position(typ))

	case *ast.CompositeLit:
		return c.Type(typ.Type)

	case *ast.Ident:
		switch typ.Name {
		// Unsupported types
		case "int64", "uint64", "complex64", "complex128": // "uintptr"
			return fmt.Errorf("%s: %s type", c.Position(typ), typ.Name)
		}

	case *ast.InterfaceType: // ToDo: review

	case *ast.MapType:
		if err := c.Type(typ.Key); err != nil {
			return err
		}
		if err := c.Type(typ.Value); err != nil {
			return err
		}

	case *ast.ParenExpr:
		return c.Type(typ.X)

	// http://golang.org/pkg/go/ast/#StarExpr || godoc go/ast StarExpr
	//  X    Expr      // operand
	case *ast.StarExpr:
		return c.Type(typ.X)

	case *ast.StructType:

	case *ast.UnaryExpr:
		// Channel
		if typ.Op == token.ARROW {
			return fmt.Errorf("%s: channel operator", c.Position(typ))
		}

		return c.Type(typ.X)

	// The type has not been indicated
	case nil:
	}

	return nil
}
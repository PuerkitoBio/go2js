// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
// If a copy of the MPL was not distributed with this file, You can obtain one at
// http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"text/template"
)

type Kind uint8

const (
	invalidKind Kind = iota
	//arrayKind TODO: remove
	//ellipsisKind
	sliceKind
	structKind
)

// expression represents a Go expression.
type expression struct {
	tr            *translate
	*bytes.Buffer // sintaxis translated

	varName  string
	funcName string
	mapName  string
	zero     string

	kind Kind

	// TODO: change all booleans by kind
	hasError bool
	useIota  bool

	//isFunc    bool // anonymous function
	isIdent      bool
	isValue      bool // is it on the right of the assignment?
	isVarAddress bool
	isPointer    bool
	isMake       bool
	isNil        bool

	arrayHasElts bool // does array has elements?
	isEllipsis   bool
	isMultiDim   bool // multi-dimensional array

	// To handle comparisons
	isBasicLit     bool
	returnBasicLit bool

	lenArray []string // the lengths of an array
	index    []string
}

// newExpression initializes a new expression.
func (tr *translate) newExpression(iVar interface{}) *expression {
	var id string

	if iVar != nil {
		switch tVar := iVar.(type) {
		case *ast.Ident:
			id = tVar.Name
		case string:
			id = tVar
		}
	}

	return &expression{
		tr,
		new(bytes.Buffer),
		id,
		"",
		"",
		"",
		invalidKind,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		make([]string, 0),
		make([]string, 0),
	}
}

// getExpression returns the Go expression translated to JavaScript.
func (tr *translate) getExpression(expr ast.Expr) *expression {
	e := tr.newExpression(nil)

	e.translate(expr)
	return e
}

// translate translates the Go expression.
func (e *expression) translate(expr ast.Expr) {
	switch typ := expr.(type) {

	// godoc go/ast ArrayType
	//  Lbrack token.Pos // position of "["
	//  Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
	//  Elt    Expr      // element type
	case *ast.ArrayType:
		// Type checking
		if _, ok := typ.Elt.(*ast.Ident); ok {
			if e.tr.getExpression(typ.Elt).hasError {
				return
			}
		}
		if typ.Len == nil { // slice
			break
		}

		if _, ok := typ.Len.(*ast.Ellipsis); ok {
			e.zero, _ = e.tr.zeroValue(true, typ.Elt)
			e.isEllipsis = true
			break
		}

		if !e.isMultiDim {
			e.WriteString("g.MkArray([")
		} else {
			e.WriteString(",")
		}
		e.WriteString(e.tr.getExpression(typ.Len).String())
		e.tr.isArray = true

		switch t := typ.Elt.(type) {
		case *ast.ArrayType: // multi-dimensional array
			e.isMultiDim = true
			e.translate(typ.Elt)
		case *ast.Ident, *ast.StarExpr: // the type is initialized
			e.zero, _ = e.tr.zeroValue(true, typ.Elt)

			e.WriteString(fmt.Sprintf("],%s", SP+e.zero))

		default:
			panic(fmt.Sprintf("*expression.translate: type unimplemented: %T", t))
		}

	// godoc go/ast BasicLit
	//  Kind     token.Token // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
	//  Value    string      // literal string
	case *ast.BasicLit:
		if typ.Value[0] == '`' { // raw string literal
			typ.Value = template.JSEscapeString(typ.Value)
			typ.Value = `"` + typ.Value[1:len(typ.Value)-1] + `"`
		}

		// Replace new lines
		if strings.Contains(typ.Value, "\\n") {
			typ.Value = strings.Replace(typ.Value, "\\n", "<br>", -1)
		}
		// Replace tabulators
		if strings.Contains(typ.Value, "\\t") {
			typ.Value = strings.Replace(typ.Value, "\\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
		}

		e.WriteString(typ.Value)
		e.isBasicLit = true

	// http://golang.org/doc/go_spec.html#Comparison_operators
	// https://developer.mozilla.org/en/JavaScript/Reference/Operators/Comparison_Operators
	//
	// godoc go/ast BinaryExpr
	//  X     Expr        // left operand
	//  Op    token.Token // operator
	//  Y     Expr        // right operand
	case *ast.BinaryExpr:
		isComparing := false
		isOpNot := false
		op := typ.Op.String()

		switch typ.Op {
		case token.NEQ:
			isOpNot = true
			fallthrough
		case token.EQL:
			op += "="
			isComparing = true
		}

		if e.tr.isConst {
			e.translate(typ.X)
			e.WriteString(SP + op + SP)
			e.translate(typ.Y)
			break
		}

		// * * *
		x := e.tr.getExpression(typ.X)
		y := e.tr.getExpression(typ.Y)

		if isComparing {
			xStr := stripField(x.String())
			yStr := stripField(y.String())

			// Slice
			if y.isNil && e.tr.isType(sliceType, xStr) {
				if isOpNot {
					e.WriteString("!")
				}
				e.WriteString(xStr + ".isNil")
				break
			}
			if x.isNil && e.tr.isType(sliceType, yStr) {
				if isOpNot {
					e.WriteString("!")
				}
				e.WriteString(yStr + ".isNil()")
				break
			}

			// Map
			if y.isNil && e.tr.isType(mapType, xStr) {
				e.WriteString(fmt.Sprintf("%s.v%s%s%sundefined", xStr, SP, op, SP))
				break
			}
			if x.isNil && e.tr.isType(mapType, yStr) {
				e.WriteString(fmt.Sprintf("%s.v%s%s%sundefined", yStr, SP, op, SP))
				break
			}
		}

		// * * *
		stringify := false

		// JavaScript only compares basic literals.
		if isComparing && !x.isBasicLit && !x.returnBasicLit && !y.isBasicLit && !y.returnBasicLit {
			stringify = true
		}

		if stringify {
			e.WriteString("JSON.stringify(" + x.String() + ")")
		} else {
			e.WriteString(x.String())
		}
		// To know when a pointer is compared with the value nil.
		if y.isNil && !x.isPointer && !x.isVarAddress {
			e.WriteString(NIL)
		}

		e.WriteString(SP + op + SP)

		if stringify {
			e.WriteString("JSON.stringify(" + y.String() + ")")
		} else {
			e.WriteString(y.String())
		}
		if x.isNil && !y.isPointer && !y.isVarAddress {
			e.WriteString(NIL)
		}

	// godoc go/ast CallExpr
	//  Fun      Expr      // function expression
	//  Args     []Expr    // function arguments; or nil
	case *ast.CallExpr:
		// == Library
		if call, ok := typ.Fun.(*ast.SelectorExpr); ok {
			e.translate(call)

			str := fmt.Sprintf("%s", e.tr.GetArgs(e.funcName, typ.Args))
			if e.funcName != "fmt.Sprintf" && e.funcName != "fmt.Sprint" {
				str = "(" + str + ")"
			}

			e.WriteString(str)
			break
		}

		// == Conversion: []byte()
		if call, ok := typ.Fun.(*ast.ArrayType); ok {
			if call.Elt.(*ast.Ident).Name == "byte" {
				e.translate(typ.Args[0])
			} else {
				panic(fmt.Sprintf("call of conversion unimplemented: []%T()", call))
			}
			break
		}

		// == Built-in functions - golang.org/pkg/builtin/
		call := typ.Fun.(*ast.Ident).Name

		switch call {
		case "make":
			// Type checking
			if e.tr.getExpression(typ.Args[0]).hasError {
				return
			}

			switch argType := typ.Args[0].(type) {
			// For slice
			case *ast.ArrayType:
				zero, _ := e.tr.zeroValue(true, argType.Elt)

				e.WriteString(fmt.Sprintf("%s,%s%s", zero, SP,
					e.tr.getExpression(typ.Args[1]))) // length

				// capacity
				if len(typ.Args) == 3 {
					e.WriteString("," + SP + e.tr.getExpression(typ.Args[2]).String())
				}

				e.tr.slices[e.tr.funcId][e.tr.blockId][e.tr.lastVarName] = void
				e.isMake = true

			case *ast.MapType:
				e.tr.maps[e.tr.funcId][e.tr.blockId][e.tr.lastVarName] = void
				if !Bootstrap {
					e.WriteString(fmt.Sprintf("g.Map(%s,%s{})", e.tr.zeroOfMap(argType), SP))
				} else {
					e.WriteString("{}")
				}

			case *ast.ChanType:
				e.translate(typ.Fun)

			default:
				panic(fmt.Sprintf("call of 'make' unimplemented: %T", argType))
			}

		case "new":
			switch argType := typ.Args[0].(type) {
			case *ast.ArrayType:
				for _, arg := range typ.Args {
					e.translate(arg)
				}

			case *ast.Ident:
				value, _ := e.tr.zeroValue(true, argType)
				e.WriteString(value)

			default:
				panic(fmt.Sprintf("call of 'new' unimplemented: %T", argType))
			}

		// == Conversion
		case "string":
			arg := e.tr.getExpression(typ.Args[0]).String()
			_arg := stripField(arg)

			if e.tr.isType(sliceType, _arg) {
				e.WriteString(_arg + ".toString()")
			} else {
				e.WriteString(arg)
				e.returnBasicLit = true
			}

		case "uint", "uint8", "uint16", "uint32",
			"int", "int8", "int16", "int32",
			"float32", "float64", "byte", "rune":
			e.translate(typ.Args[0])
			e.returnBasicLit = true
		// ==

		case "print", "println":
			e.WriteString(fmt.Sprintf("alert(%s)", e.tr.GetArgs(call, typ.Args)))

		case "len":
			arg := e.tr.getExpression(typ.Args[0]).String()
			_arg := stripField(arg)

			if e.tr.isType(sliceType, _arg) {
				e.WriteString(_arg + ".len")
			} else if e.tr.isType(mapType, arg) {
				e.WriteString(arg + ".len()")
			} else {
				e.WriteString(arg + ".length")
			}

			e.returnBasicLit = true

		case "cap":
			arg := e.tr.getExpression(typ.Args[0]).String()
			_arg := stripField(arg)

			if e.tr.isType(sliceType, _arg) {
				if strings.HasSuffix(arg, VALUE_FIELD) {
					arg = _arg
				}
			}

			e.WriteString(arg + ".cap")
			e.returnBasicLit = true

		case "delete":
			e.WriteString(fmt.Sprintf("delete %s%s[%s]",
				e.tr.getExpression(typ.Args[0]).String(),
				VALUE_FIELD,
				e.tr.getExpression(typ.Args[1]).String()))

		case "panic":
			e.WriteString(fmt.Sprintf("throw new Error(%s)",
				e.tr.getExpression(typ.Args[0])))

		// == Not supported
		case "recover", "complex":
			e.tr.addError("%s: built-in function %s()",
				e.tr.fset.Position(typ.Fun.Pos()), call)
			e.tr.hasError = true
			return
		case "int64", "uint64":
			e.tr.addError("%s: conversion of type %s",
				e.tr.fset.Position(typ.Fun.Pos()), call)
			e.tr.hasError = true
			return

		// == Not implemented
		case "append", "close", "copy", "uintptr":
			panic(fmt.Sprintf("built-in call unimplemented: %s", call))

		// Defined functions
		default:
			args := ""

			for i, v := range typ.Args {
				if i != 0 {
					args += "," + SP
				}
				args += e.tr.getExpression(v).String()
			}

			e.WriteString(fmt.Sprintf("%s(%s)", call, args))
		}

	// godoc go/ast ChanType
	//  Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
	//  Dir   ChanDir   // channel direction
	//  Value Expr      // value type
	case *ast.ChanType:
		e.tr.addError("%s: channel type", e.tr.fset.Position(typ.Pos()))
		e.tr.hasError = true
		return

	// godoc go/ast CompositeLit
	//  Type   Expr      // literal type; or nil
	//  Lbrace token.Pos // position of "{"
	//  Elts   []Expr    // list of composite elements; or nil
	//  Rbrace token.Pos // position of "}"
	case *ast.CompositeLit:
		switch compoType := typ.Type.(type) {
		case *ast.ArrayType:
			if !e.arrayHasElts {
				e.translate(typ.Type)
			}

			if e.isEllipsis {
				e.WriteString(fmt.Sprintf("g.MkArray([%s],%s,%s",
					strconv.Itoa(len(typ.Elts)), SP+e.zero, SP))

				e.WriteString("[")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString("])")
				break
			}

			//e.kind = arrayKind TODO: remove

			// Slice
			if compoType.Len == nil {
				e.tr.slices[e.tr.funcId][e.tr.blockId][e.tr.lastVarName] = void
				e.kind = sliceKind
			}
			// Struct
			if elt, ok := compoType.Elt.(*ast.StructType); ok {
				e.tr.getStruct(elt, "", false)
				e.kind = structKind
			}
			// For arrays with elements
			if len(typ.Elts) != 0 {
				if !e.arrayHasElts && compoType.Len != nil {
					e.WriteString("," + SP)
					e.arrayHasElts = true
				}
				if e.kind == sliceKind {
					e.zero, _ = e.tr.zeroValue(true, compoType.Elt)
					e.WriteString(e.zero + "," + SP)
				}
				e.WriteString("[")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString("]")

			} /*else if e.kind == sliceKind {
				e.WriteString("[]")
			}*/

		case *ast.Ident: // Custom types
			useField := false
			e.isVarAddress = false // it is the address to a type
			e.WriteString("new " + typ.Type.(*ast.Ident).Name)

			if len(typ.Elts) != 0 {
				// Specify the fields
				if _, ok := typ.Elts[0].(*ast.KeyValueExpr); ok {
					useField = true

					e.WriteString("();")
					e.writeTypeElts(typ.Elts, typ.Lbrace)
				}
			}
			if !useField {
				e.WriteString("(")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString(")")
			}

		case *ast.MapType:
			// Type checking
			if e.tr.getExpression(typ.Type).hasError {
				return
			}
			e.tr.maps[e.tr.funcId][e.tr.blockId][e.tr.lastVarName] = void

			e.WriteString(fmt.Sprintf("g.Map(%s,%s{", e.tr.zeroOfMap(compoType), SP))
			e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
			e.WriteString("})")

		case nil:
			if e.kind == structKind {
				e.WriteString("_(")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString(")")
			} else {
				e.WriteString("[")
				e.writeElts(typ.Elts, typ.Lbrace, typ.Rbrace)
				e.WriteString("]")
			}

		default:
			panic(fmt.Sprintf("'CompositeLit' unimplemented: %T", compoType))
		}

	// godoc go/ast Ellipsis
	//  Ellipsis token.Pos // position of "..."
	//  Elt      Expr      // ellipsis element type (parameter lists only); or nil
	//case *ast.Ellipsis:

	// http://golang.org/doc/go_spec.html#Function_literals
	// https://developer.mozilla.org/en/JavaScript/Reference/Functions_and_function_scope#Function_constructor_vs._function_declaration_vs._function_expression
	// godoc go/ast FuncLit
	//
	//  Type *FuncType  // function type
	//  Body *BlockStmt // function body
	case *ast.FuncLit:
		e.translate(typ.Type)
		e.tr.getStatement(typ.Body)

	// godoc go/ast FuncType
	//  Func    token.Pos  // position of "func" keyword
	//  Params  *FieldList // (incoming) parameters; or nil
	//  Results *FieldList // (outgoing) results; or nil
	case *ast.FuncType:
		//e.isFunc = true
		e.tr.writeFunc(nil, nil, typ)

	// godoc go/ast Ident
	//  Name    string    // identifier name
	case *ast.Ident:
		name := typ.Name

		switch name {
		case "iota":
			e.WriteString(IOTA)
			e.useIota = true

		// Undefined value in array / slice
		case "_":
			if len(e.lenArray) == 0 {
				e.WriteString(name)
			}
		// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/undefined
		case "nil":
			e.WriteString("undefined")
			e.isBasicLit = true
			e.isNil = true

		// Not supported
		case "int64", "uint64", "complex64", "complex128":
			e.tr.addError("%s: %s type", e.tr.fset.Position(typ.Pos()), name)
			e.tr.hasError = true
		// Not implemented
		case "uintptr":
			e.tr.addError("%s: unimplemented type %q", e.tr.fset.Position(typ.Pos()), name)
			e.tr.hasError = true

		default:
			name = validIdent(typ.Name)

			if e.isPointer { // `*x` => `x.POINTER_FIELD`
				name += POINTER_FIELD
			} else if e.isVarAddress { // `&x` => `x`
				e.tr.addPointer(name)
			} else {
				if !e.tr.isVar {
					isSlice := false

					if e.tr.isType(sliceType, name) {
						isSlice = true
					}
					if name == e.tr.recvVar {
						name = "this" + TYPE_FIELD
					}
					if isSlice {
						name += VALUE_FIELD // slice field
					}

					if _, ok := e.tr.vars[e.tr.funcId][e.tr.blockId][name]; ok {
						name += tagPointer(false, 'P', e.tr.funcId, e.tr.blockId, name)
					}
				} else {
					e.isIdent = true
				}
			}

			e.WriteString(name)
		}

	// godoc go/ast IndexExpr
	// Represents an expression followed by an index.
	//  X      Expr      // expression
	//  Lbrack token.Pos // position of "["
	//  Index  Expr      // index expression
	//  Rbrack token.Pos // position of "]"
	case *ast.IndexExpr:
		// == Store indexes
		e.index = append(e.index, e.tr.getExpression(typ.Index).String())

		// Could be multi-dimensional
		if _, ok := typ.X.(*ast.IndexExpr); ok {
			e.translate(typ.X)
			return
		}
		// ==

		x := e.tr.getExpression(typ.X).String()
		index := ""
		indexArgs := ""

		for i := len(e.index) - 1; i >= 0; i-- { // inverse order
			idx := e.index[i]
			index += "[" + idx + "]"

			if indexArgs != "" {
				indexArgs += "," + SP
			}
			indexArgs += idx
		}

		if e.tr.isType(mapType, x) {
			e.mapName = x

			if e.tr.isVar && !e.isValue {
				e.WriteString(x + VALUE_FIELD + index)
			} else {
				e.WriteString(x + ".get(" + indexArgs + ")[0]")
			}
		} else if e.tr.isType(sliceType, x) {
			e.WriteString(x + VALUE_FIELD + index)
		} else {
			e.WriteString(x + index)
		}

	// godoc go/ast InterfaceType
	//  Interface  token.Pos  // position of "interface" keyword
	//  Methods    *FieldList // list of methods
	//  Incomplete bool       // true if (source) methods are missing in the Methods list
	case *ast.InterfaceType: // TODO: review

	// godoc go/ast KeyValueExpr
	//  Key   Expr
	//  Colon token.Pos // position of ":"
	//  Value Expr
	case *ast.KeyValueExpr:
		key := e.tr.getExpression(typ.Key).String()
		exprValue := e.tr.getExpression(typ.Value)
		value := exprValue.String()

		if value[0] == '[' { // multi-dimensional index
			value = "{" + value[1:len(value)-1] + "}"
		}

		if !e.tr.isArray && e.kind != sliceKind {
			e.WriteString(key + ":" + SP + value)
		} else {
			e.WriteString(fmt.Sprintf("{%s:%s}", key, value))
		}

	// godoc go/ast MapType
	//  Map   token.Pos // position of "map" keyword
	//  Key   Expr
	//  Value Expr
	case *ast.MapType:
		// For type checking
		e.tr.getExpression(typ.Key)
		e.tr.getExpression(typ.Value)

	// godoc go/ast ParenExpr
	//  Lparen token.Pos // position of "("
	//  X      Expr      // parenthesized expression
	//  Rparen token.Pos // position of ")"
	case *ast.ParenExpr:
		e.translate(typ.X)

	// godoc go/ast SelectorExpr
	//   X   Expr   // expression
	//   Sel *Ident // field selector
	case *ast.SelectorExpr:
		isPkg := false
		x := ""

		switch t := typ.X.(type) {
		case *ast.SelectorExpr:
			e.translate(typ.X)
		case *ast.Ident:
			x = t.Name
		case *ast.IndexExpr:
			e.translate(t)
			e.WriteString("." + typ.Sel.Name) // TODO: validIdent?
			return
		default:
			panic(fmt.Sprintf("'SelectorExpr': unimplemented: %T", t))
		}

		if x == e.tr.recvVar {
			x = "this"
		}
		goName := x + "." + validIdent(typ.Sel.Name)

		// Check is the selector is a package
		for _, v := range validImport {
			if v == x {
				isPkg = true
				break
			}
		}

		// Check if it can be translated to its equivalent in JavaScript.
		if isPkg {
			jsName, ok := Function[goName]
			if !ok {
				jsName, ok = Constant[goName]
			}

			if !ok {
				e.tr.addError(fmt.Errorf("%s: %q not supported in JS",
					e.tr.fset.Position(typ.Sel.Pos()), goName))
				e.tr.hasError = true
				break
			}

			e.funcName = goName
			e.WriteString(jsName)
		} else {
			/*if _, ok := e.tr.zeroType[x]; !ok {
				panic("selector: " + x)
			}*/

			e.WriteString(goName)
		}

	// godoc go/ast SliceExpr
	//  X      Expr      // expression
	//  Lbrack token.Pos // position of "["
	//  Low    Expr      // begin of slice range; or nil
	//  High   Expr      // end of slice range; or nil
	//  Rbrack token.Pos // position of "]"
	case *ast.SliceExpr:
		slice := "0"
		x := typ.X.(*ast.Ident).Name

		if typ.Low != nil {
			slice = typ.Low.(*ast.BasicLit).Value // e.tr.getExpression(typ.Low).String()
		}
		if typ.High != nil {
			slice += "," + SP + typ.High.(*ast.BasicLit).Value // e.tr.getExpression(typ.High).String()
		}

		if e.tr.isVar {
			e.WriteString(x + "," + SP + slice)
		} else {
			e.WriteString(fmt.Sprintf("g.Slice(%s,%s)", x, SP+slice))
		}

		e.kind = sliceKind

	// godoc go/ast StarExpr
	//  Star token.Pos // position of "*"
	//  X    Expr      // operand
	case *ast.StarExpr:
		e.isPointer = true
		e.translate(typ.X)

	// godoc go/ast StructType
	//  Struct     token.Pos  // position of "struct" keyword
	//  Fields     *FieldList // list of field declarations
	//  Incomplete bool       // true if (source) fields are missing in the Fields list
	case *ast.StructType:

	// godoc go/ast UnaryExpr
	//  OpPos token.Pos   // position of Op
	//  Op    token.Token // operator
	//  X     Expr        // operand
	case *ast.UnaryExpr:
		writeOp := true
		op := typ.Op.String()

		switch typ.Op {
		case token.XOR: // bitwise complement
			op = "~"
		case token.AND: // address operator
			e.isVarAddress = true
			writeOp = false
		case token.ARROW: // channel
			e.tr.addError("%s: channel operator", e.tr.fset.Position(typ.OpPos))
			e.tr.hasError = true
			return
		}

		if writeOp {
			e.WriteString(op)
		}
		e.translate(typ.X)

	// The type has not been indicated
	case nil:

	default:
		panic(fmt.Sprintf("unimplemented: %T", expr))
	}
}

// == Utility
//

// writeElts writes the list of composite elements.
func (e *expression) writeElts(elts []ast.Expr, Lbrace, Rbrace token.Pos) {
	firstPos := e.tr.getLine(Lbrace)
	posOldElt := firstPos
	posNewElt := 0

	for i, el := range elts {
		posNewElt = e.tr.getLine(el.Pos())

		if i != 0 {
			e.WriteString(",")
		}
		if posNewElt != posOldElt {
			e.WriteString(strings.Repeat(NL, posNewElt-posOldElt))
			e.WriteString(strings.Repeat(TAB, e.tr.tabLevel+1))
		} else if i != 0 { // in the same line
			e.WriteString(SP)
		}

		// It is necessary to create a new expression for each element, to avoid
		// add indexes in case that there are expressiones with them.
		exprElt := e.tr.newExpression(nil)

		exprElt.kind = e.kind
		exprElt.hasError = e.hasError
		exprElt.useIota = e.useIota
		exprElt.arrayHasElts = e.arrayHasElts
		exprElt.isEllipsis = e.isEllipsis
		exprElt.isMultiDim = e.isMultiDim
		exprElt.isValue = e.isValue

		exprElt.translate(el)
		e.WriteString(exprElt.String())

		posOldElt = posNewElt
	}

	// The right brace
	posNewElt = e.tr.getLine(Rbrace)
	if posNewElt != posOldElt {
		e.WriteString(strings.Repeat(NL, posNewElt-posOldElt))
		e.WriteString(strings.Repeat(TAB, e.tr.tabLevel))
	}

	e.tr.line += posNewElt - firstPos // update the global position
}

// writeTypeElts writes the list of elements for a custom type.
func (e *expression) writeTypeElts(elts []ast.Expr, Lbrace token.Pos) {
	firstPos := e.tr.getLine(Lbrace)
	posOldElt := firstPos
	posNewElt := 0
	useBracket := false

	for i, el := range elts {
		posNewElt = e.tr.getLine(el.Pos())
		kv := el.(*ast.KeyValueExpr)
		key := e.tr.getExpression(kv.Key).String()

		if i == 0 {
			if strings.HasPrefix(key, `"`) {
				useBracket = true
			} else {
				useBracket = false
			}
		}
		if useBracket {
			key = "[" + key + "]"
		} else {
			key = "." + key
		}

		if i != 0 {
			e.WriteString(",")
		}
		if posNewElt != posOldElt {
			e.WriteString(strings.Repeat(NL, posNewElt-posOldElt))
			e.WriteString(strings.Repeat(TAB, e.tr.tabLevel))
		} else { // in the same line
			e.WriteString(SP)
		}

		e.WriteString(fmt.Sprintf("%s%s=%s",
			e.tr.lastVarName,
			key+SP,
			SP+e.tr.getExpression(kv.Value).String(),
		))

		posOldElt = posNewElt
	}
	e.tr.line += posNewElt - firstPos // update the global position
}

// * * *

// stripField strips the field name VALUE_FIELD.
func stripField(name string) string {
	if strings.HasSuffix(name, VALUE_FIELD) {
		return name[:len(name)-2]
	}
	return name
}

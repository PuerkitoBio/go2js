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
	"strings"
)

// Maximum number of expressions to get.
// The expressions are the values after of "=".
const MAX_EXPRESSION = 10

// Imports
//
// http://golang.org/doc/go_spec.html#Import_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/import
func (tr *transform) getImport(spec []ast.Spec) {

	// http://golang.org/pkg/go/ast/#ImportSpec || godoc go/ast ImportSpec
	//  Name    *Ident        // local package name (including "."); or nil
	//  Path    *BasicLit     // import path
	for _, v := range spec {
		iSpec := v.(*ast.ImportSpec)
		path := iSpec.Path.Value
		pathDir := strings.SplitN(path, "/", 2)[0]

		if !strings.Contains(pathDir, ".") {
			tr.addError("%s: import from core library", path)
			continue
		}

		//import objectName.*; 
		//fmt.Println(iSpec.Name, pathDir)
	}
}

// Constants
//
// http://golang.org/doc/go_spec.html#Constant_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/const
func (tr *transform) getConst(spec []ast.Spec) {
	iotas := make([]int, MAX_EXPRESSION)
	lastValues := make([]string, MAX_EXPRESSION)

	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	//  Names   []*Ident      // value names (len(Names) > 0)
	//  Type    Expr          // value type; or nil
	//  Values  []Expr        // initial values; or nil
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		if len(vSpec.Values) > MAX_EXPRESSION {
			panic("length of 'iotas' is lesser than 'vSpec.Values'")
		}

		// Checking
		if err := newCheck(tr.fset).Type(vSpec.Type); err != nil {
			tr.addError(err)
			continue
		}

		names, skipName := tr.getName(vSpec)

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		//  type Expr interface
		values := make([]string, 0)

		if len(vSpec.Values) != 0 {
			for i, v := range vSpec.Values {
				var dst string

				// Checking
				if err := newCheck(tr.fset).Type(v); err != nil {
					tr.addError(err)
					continue
				}

				src := newExpression(names[i])
				src.transform(v)

				if src.useIota {
					dst = fmt.Sprintf(src.String(), iotas[i])
					iotas[i]++
				} else {
					dst = src.String()
				}

				if !skipName[i] {
					values = append(values, dst)
					lastValues[i] = src.String()
				}
			}
		} else { // get last value of iota
			for i := 0; i < len(names); i++ {
				dst := fmt.Sprintf(lastValues[i], iotas[0])
				values = append(values, dst)
			}
			iotas[0]++
		}

		// Skip write buffer, if any error
		if tr.hasError {
			continue
		}

		// === Write
		// TODO: calculate expression using "exp/types"
		tr.addLine(vSpec.Pos())

		isFirst := true
		for i, v := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				tr.WriteString(fmt.Sprintf("const %s%s=%s%s", v, SP, SP, values[i]))
			} else {
				tr.WriteString(fmt.Sprintf(",%s%s%s=%s%s", SP, v, SP, SP, values[i]))
			}
		}

		// It is possible that there is only a blank identifier
		if isFirst {
			continue
		}

		tr.WriteString(";")
	}
}

// Variables
//
// http://golang.org/doc/go_spec.html#Variable_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/var
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/let
//
// TODO: use let for local variables
func (tr *transform) getVar(spec []ast.Spec) {
	// http://golang.org/pkg/go/ast/#ValueSpec || godoc go/ast ValueSpec
	for _, s := range spec {
		vSpec := s.(*ast.ValueSpec)

		// Checking
		if err := newCheck(tr.fset).Type(vSpec.Type); err != nil {
			tr.addError(err)
			continue
		}

		names, skipName := tr.getName(vSpec)

		// === Values
		// http://golang.org/pkg/go/ast/#Expr || godoc go/ast Expr
		values := make([]string, 0)

		for i, v := range vSpec.Values {
			// Checking
			if err := newCheck(tr.fset).Type(v); err != nil {
				tr.addError(err)
				continue
			}

			// Skip when it is not a function
			if skipName[i] {
				if _, ok := v.(*ast.CallExpr); !ok {
					continue
				}
			}

			if !skipName[i] {
				src := newExpression(names[i])
				src.transform(v)

				values = append(values, src.String())
			}
		}

		if tr.hasError {
			continue
		}

		// === Write
		// TODO: calculate expression using "exp/types"

		tr.addLine(vSpec.Pos())

		isFirst := true
		for i, n := range names {
			if skipName[i] {
				continue
			}

			if isFirst {
				isFirst = false
				tr.WriteString("var " + n)
			} else {
				tr.WriteString("," + SP + n)
			}

			if len(values) != 0 {
				tr.WriteString(SP + "=" + SP + values[i])
			}
		}

		last := tr.Bytes()[tr.Len()-1] // last character

		if last != '}' && last != ';' {
			tr.WriteString(";")
		}
	}
}

// Types
//
// http://golang.org/doc/go_spec.html#Type_declarations
func (tr *transform) getType(spec []ast.Spec) {
	// Format fields
	format := func(fields []string) (args, allFields string) {
		for i, f := range fields {
			if i == 0 {
				args = f
			} else {
				args += "," + SP + f
				allFields += SP
			}

			allFields += fmt.Sprintf("this.%s=%s;", f, f)
		}
		return
	}

	// http://golang.org/pkg/go/ast/#TypeSpec || godoc go/ast TypeSpec
	//  Name    *Ident        // type name
	//  Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
	for _, s := range spec {
		tSpec := s.(*ast.TypeSpec)
		fields := make([]string, 0) // names of fields
		//!anonField := make([]bool, 0) // anonymous field

		// Checking
		if err := newCheck(tr.fset).Type(tSpec.Type); err != nil {
			tr.addError(err)
			continue
		}

		switch typ := tSpec.Type.(type) {
		default:
			panic(fmt.Sprintf("unimplemented: %T", typ))

		case *ast.Ident:

		// http://golang.org/pkg/go/ast/#StructType || godoc go/ast StructType
		//  Struct     token.Pos  // position of "struct" keyword
		//  Fields     *FieldList // list of field declarations
		//  Incomplete bool       // true if (source) fields are missing in the Fields list
		case *ast.StructType:
			if typ.Incomplete {
				panic("list of fields incomplete ???")
			}

			// http://golang.org/pkg/go/ast/#FieldList || godoc go/ast FieldList
			//  List    []*Field  // field list; or nil
			for _, field := range typ.Fields.List {
				if _, ok := field.Type.(*ast.FuncType); ok {
					tr.addError("%s: function type in struct",
						tr.fset.Position(field.Pos()))
					continue
				}

				// http://golang.org/pkg/go/ast/#Field || godoc go/ast Field
				//  Names   []*Ident      // field/method/parameter names; or nil if anonymous field
				//  Type    Expr          // field/method/parameter type
				//  Tag     *BasicLit     // field tag; or nil

				// Checking
				if err := newCheck(tr.fset).Type(field.Type); err != nil {
					tr.addError(err)
					continue
				}
				if field.Names == nil {
					tr.addError("%s: anonymous field in struct",
						tr.fset.Position(field.Pos()))
					continue
				}

				for _, n := range field.Names {
					name := n.Name

					if name == "_" {
						continue
					}

					fields = append(fields, name)
					//!anonField = append(anonField, false)
				}
			}
		}

		if tr.hasError {
			continue
		}

		// === Write
		name := tSpec.Name.Name
		args, allFields := format(fields)
		tr.addLine(tSpec.Pos())

		tr.WriteString(fmt.Sprintf("function %s(%s)%s{", name, args, SP))

		if len(allFields) != 0 {
			tr.WriteString(allFields)
			tr.WriteString("}")
		} else {
			tr.WriteString("}") //! empty struct
		}
		// ===

		// To export
		tr.checkPublic(name)
	}
}

// Functions
//
// http://golang.org/doc/go_spec.html#Function_declarations
// https://developer.mozilla.org/en/JavaScript/Reference/Statements/function
func (tr *transform) getFunc(decl *ast.FuncDecl) {
	// http://golang.org/pkg/go/ast/#FuncDecl || godoc go/ast FuncDecl
	//  Recv *FieldList    // receiver (methods); or nil (functions)
	//  Name *Ident        // function/method name
	//  Type *FuncType     // position of Func keyword, parameters and results
	//  Body *BlockStmt    // function body; or nil (forward declaration)

	// http://golang.org/pkg/go/ast/#FieldList || godoc go/ast FieldList
	//  List    []*Field  // field list; or nil
	//
	// http://golang.org/pkg/go/ast/#Field || godoc go/ast Field
	//  Names   []*Ident      // field/method/parameter names; or nil if anonymous field
	//  Type    Expr          // field/method/parameter type
	//  Tag     *BasicLit     // field tag; or nil

	// Check empty functions
	if len(decl.Body.List) == 0 {
		return
	}

	tr.addLine(decl.Pos())
	tr.WriteString(fmt.Sprintf("function %s(%s)%s",
		decl.Name, getParams(decl.Type), SP))
	tr.getStatement(decl.Body)
	tr.checkPublic(decl.Name.Name) // to export
}

//
// === Utility

// Gets the identifiers.
//
// http://golang.org/pkg/go/ast/#Ident || godoc go/ast Ident
//  Name    string    // identifier name
func (tr *transform) getName(spec *ast.ValueSpec) (names []string, skipName []bool) {
	skipName = make([]bool, len(spec.Names)) // for blank identifiers "_"

	for i, v := range spec.Names {
		if v.Name == "_" {
			skipName[i] = true
			continue
		}
		names = append(names, v.Name)

		// To export
		tr.checkPublic(v.Name)
	}

	return
}

// Gets the parameters.
//
// http://golang.org/pkg/go/ast/#FuncType || godoc go/ast FuncType
//  Params  *FieldList // (incoming) parameters; or nil
//  Results *FieldList // (outgoing) results; or nil
func getParams(f *ast.FuncType) string {
	s := ""

	for i, v := range f.Params.List {
		if i != 0 {
			s += "," + SP
		}
		s += v.Names[0].Name
	}

	return s
}
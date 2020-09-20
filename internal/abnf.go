// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal is a library for EBNF grammars. The input is text ([]byte)
// satisfying the following grammar (represented itself in EBNF):
//
//		rulelist       =  1*( rule / (*c-wsp c-nl) )
//		rule           =  rulename defined-as elements c-nl
//								; continues if next line starts
// 								;  with white space
//      rulename       =  ALPHA *(ALPHA / DIGIT / "-")
//		defined-as     =  *c-wsp ("=" / "=/") *c-wsp
//                              ; basic rules definition and
//                              ;  incremental alternatives
//      elements       =  alternation *c-wsp
//      c-wsp          =  WSP / (c-nl WSP)
//      c-nl           =  comment / CRLF
//                             ; comment or newline
//      comment        =  ";" *(WSP / VCHAR) CRLF
//      alternation    =  concatenation
//                        *(*c-wsp "/" *c-wsp concatenation)
//      concatenation  =  repetition *(1*c-wsp repetition)
//      repetition     =  [repeat] element
//      repeat         =  1*DIGIT / (*DIGIT "*" *DIGIT)
//      element        =  rulename / group / option /
//                        char-val / num-val / prose-val
//      group          =  "(" *c-wsp alternation *c-wsp ")"
//      option         =  "[" *c-wsp alternation *c-wsp "]"
//      char-val       =  DQUOTE *(%x20-21 / %x23-7E) DQUOTE
//                             ; quoted string of SP and VCHAR
//                             ;  without DQUOTE
//      num-val        =  "%" (bin-val / dec-val / hex-val)
//      bin-val        =  "b" 1*BIT
//                        [ 1*("." 1*BIT) / ("-" 1*BIT) ]
//                             ; series of concatenated bit values
//                             ;  or single ONEOF range
//      dec-val        =  "d" 1*DIGIT
//                        [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]
//      hex-val        =  "x" 1*HEXDIG
//                        [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
//      prose-val      =  "<" *(%x20-3D / %x3F-7E) ">"
//                             ; bracketed string of SP and VCHAR
//                             ;  without angles
//                             ; prose description, to be used as
//                             ;  last resort
//
// A name is a Go identifier, a token is a Go string, and comments
// and white space follow the same rules as for the Go language.
// Production names starting with an uppercase Unicode letter denote
// non-terminal productions (i.e., productions which allow white-space
// and comments between tokens); all other production names denote
// lexical productions.
//
package internal

import (
	"fmt"
	"text/scanner"
	"unicode"
	"unicode/utf8"
)

// ----------------------------------------------------------------------------
// Error handling

type errorList []error

func (list errorList) Err() error {
	if len(list) == 0 {
		return nil
	}
	return list
}

func (list errorList) Error() string {
	switch len(list) {
	case 0:
		return "no errors"
	case 1:
		return list[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", list[0], len(list)-1)
}

func newError(pos scanner.Position, msg string) error {
	return fmt.Errorf("%s: %s", pos, msg)
}

// ----------------------------------------------------------------------------
// Internal representation

type (
	// An Expression node represents a production expression.
	Expression interface {
		// Pos is the position of the first character of the syntactic construct
		Pos() scanner.Position
	}

	// An Alternative node represents a non-empty list of alternative expressions.
	Alternative []Expression // x | y | z

	// A Sequence node represents a non-empty list of sequential expressions.
	Sequence []Expression // x y z

	// A Name node represents a production name.
	Name struct {
		StringPos scanner.Position
		String    string
	}

	// A Token node represents a literal.
	Token struct {
		StringPos scanner.Position
		String    string
	}

	// A Range node represents a range of characters.
	Range struct {
		Begin, End *Token     // begin ... end
		Exclude    Expression // - token or range
	}

	// A Group node represents a grouped expression.
	Group struct {
		Lparen scanner.Position
		Body   Expression // (body)
	}

	// An Option node represents an optional expression.
	Option struct {
		Lbrack scanner.Position
		Body   Expression // [body]
	}

	// A Repetition node represents a repeated expression.
	Repetition struct {
		Lbrace scanner.Position
		Body   Expression // {body}
	}

	// A Production node represents an EBNF production.
	Production struct {
		Name *Name
		Expr Expression
	}

	// A Bad node stands for pieces of source code that lead to a parse error.
	Bad struct {
		TokPos scanner.Position
		Error  string // parser error message
	}

	// A Grammar is a set of EBNF productions. The map
	// is indexed by production name.
	//
	Grammar struct {
		Index       map[string]*Production
		Productions []*NamedProduction
	}

	// A NamedProduction is an ordered Production.
	NamedProduction struct {
		Name       string
		Production *Production
	}
)

func (x Alternative) Pos() scanner.Position { return x[0].Pos() } // the parser always generates non-empty Alternative
func (x Sequence) Pos() scanner.Position    { return x[0].Pos() } // the parser always generates non-empty Sequences
func (x *Name) Pos() scanner.Position       { return x.StringPos }
func (x *Token) Pos() scanner.Position      { return x.StringPos }
func (x *Range) Pos() scanner.Position      { return x.Begin.Pos() }
func (x *Group) Pos() scanner.Position      { return x.Lparen }
func (x *Option) Pos() scanner.Position     { return x.Lbrack }
func (x *Repetition) Pos() scanner.Position { return x.Lbrace }
func (x *Production) Pos() scanner.Position { return x.Name.Pos() }
func (x *Bad) Pos() scanner.Position        { return x.TokPos }

// ----------------------------------------------------------------------------
// Grammar verification

func isLexical(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return !unicode.IsUpper(ch)
}

type verifier struct {
	errors   errorList
	worklist []*Production
	reached  Grammar // set of productions reached from (and including) the root production
	grammar  Grammar
}

func (v *verifier) error(pos scanner.Position, msg string) {
	v.errors = append(v.errors, newError(pos, msg))
}

func (v *verifier) push(prod *Production) {
	name := prod.Name.String
	if _, found := v.reached.Index[name]; !found {
		v.worklist = append(v.worklist, prod)
		v.reached.Index[name] = prod
	}
}

func (v *verifier) verifyChar(x *Token) rune {
	s := x.String
	if utf8.RuneCountInString(s) != 1 {
		v.error(x.Pos(), "single char expected, found "+s)
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(s)
	return ch
}

func (v *verifier) verifyExpr(expr Expression, lexical bool) {
	switch x := expr.(type) {
	case nil:
		// empty expression
	case Alternative:
		for _, e := range x {
			v.verifyExpr(e, lexical)
		}
	case Sequence:
		for _, e := range x {
			v.verifyExpr(e, lexical)
		}
	case *Name:
		// a production with this name must exist;
		// add it to the worklist if not yet processed
		if prod, found := v.grammar.Index[x.String]; found {
			v.push(prod)
		} else {
			v.error(x.Pos(), "missing production "+x.String)
		}
		// within a lexical production references
		// to non-lexical productions are invalid
		if lexical && !isLexical(x.String) {
			v.error(x.Pos(), "reference to non-lexical production "+x.String)
		}
	case *Token:
		// nothing to do for now
	case *Range:
		i := v.verifyChar(x.Begin)
		j := v.verifyChar(x.End)
		if i >= j {
			v.error(x.Pos(), "decreasing character range")
		}
	case *Group:
		v.verifyExpr(x.Body, lexical)
	case *Option:
		v.verifyExpr(x.Body, lexical)
	case *Repetition:
		v.verifyExpr(x.Body, lexical)
	case *Bad:
		v.error(x.Pos(), x.Error)
	default:
		panic(fmt.Sprintf("internal error: unexpected type %T", expr))
	}
}

func (v *verifier) verify(grammar Grammar, start string) {
	// find root production
	root, found := grammar.Index[start]
	if !found {
		var noPos scanner.Position
		v.error(noPos, "no start production "+start)
		return
	}

	// initialize verifier
	v.worklist = v.worklist[0:0]
	v.reached = Grammar{Index: map[string]*Production{}}
	v.grammar = grammar

	// work through the worklist
	v.push(root)
	for {
		n := len(v.worklist) - 1
		if n < 0 {
			break
		}
		prod := v.worklist[n]
		v.worklist = v.worklist[0:n]
		v.verifyExpr(prod.Expr, isLexical(prod.Name.String))
	}

	// check if all productions were reached
	if len(v.reached.Index) < len(v.grammar.Index) {
		for name, prod := range v.grammar.Index {
			if _, found := v.reached.Index[name]; !found {
				v.error(prod.Pos(), name+" is unreachable")
			}
		}
	}
}

// Verify checks that:
//	- all productions used are defined
//	- all productions defined are used when beginning at start
//	- lexical productions refer only to other lexical productions
//
// Position information is interpreted relative to the file set fset.
//
func Verify(grammar Grammar, start string) error {
	var v verifier
	v.verify(grammar, start)
	return v.errors.Err()
}
package rules

import "errors"

var ErrorNoMatch = errors.New("No match")

type Rule interface {
	// Definition is how the rule was written down
	Definition() string

	// Children are all the rules that make up this rule. If the rule is
	// terminal than the this returns nil.
	Children() []*Rule
}

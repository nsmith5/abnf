package rules

import (
	"errors"
	"strconv"
	"strings"
)

var ErrorInvalidString = errors.New("Invalid string rule")

type String struct {
	def           string
	str           string
	caseSensitive bool
}

func (s String) Definition() string {
	return s.def
}

func (s String) Children() []*Rule {
	return nil
}

func NewString(def string) (Rule, error) {
	var str string
	var caseSensitive bool

	switch {
	// Case sensitive string literal
	case strings.HasPrefix(def, `%s"`):
		str = strings.TrimPrefix(def, `%s"`)
		caseSensitive = true

	// Old school case sensitive string (e.g %d97.98.97 for 'aba')
	case strings.HasPrefix(def, `%d`) || strings.HasPrefix(def, `%x`) || strings.HasPrefix(def, `%b`):
		var base int
		switch def[1] {
		case 'x':
			base = 16
		case 'd':
			base = 10
		case 'b':
			base = 2
		}

		var b strings.Builder
		parts := strings.Split(strings.TrimLeft(def, `%dbx`), ".")
		for _, part := range parts {
			number, err := strconv.ParseUint(part, base, 8)
			if err != nil {
				return nil, ErrorInvalidString
			}
			b.WriteByte(byte(number))
		}
		str = b.String()
		caseSensitive = true

	// Case insensitive string literal
	case strings.HasPrefix(def, `%i"`) || strings.HasPrefix(def, `"`):
		str = strings.TrimPrefix(def, `"`)
		caseSensitive = false

	// Bad data
	default:
		return nil, ErrorInvalidString
	}

	str = strings.TrimSuffix(def, `"`)
	return String{def, str, caseSensitive}, nil
}

func MustString(def string) Rule {
	rule, err := NewString(def)
	if err != nil {
		panic(err)
	}
	return rule
}

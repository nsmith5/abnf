package rules

import (
	"errors"
	"strconv"
	"strings"
)

var ErrorInvalidValueRange = errors.New("Invalid value range rule")

type ValueRange struct {
	def   string
	lower byte
	upper byte
}

func (vr ValueRange) Definition() string {
	return vr.def
}

func (vr ValueRange) Children() []*Rule {
	return nil
}

func NewValueRange(def string) (Rule, error) {
	var base int
	switch {
	case strings.HasPrefix(def, `%x`):
		base = 16
	case strings.HasPrefix(def, `%d`):
		base = 10
	case strings.HasPrefix(def, `%b`):
		base = 2
	default:
		return nil, ErrorInvalidValueRange
	}

	parts := strings.Split(strings.TrimLeft(def, `%xdb`), `-`)
	if len(parts) != 2 {
		return nil, ErrorInvalidValueRange
	}

	lower, err := strconv.ParseUint(parts[0], base, 8)
	if err != nil {
		return nil, ErrorInvalidValueRange
	}

	upper, err := strconv.ParseUint(parts[1], base, 8)
	if err != nil {
		return nil, ErrorInvalidValueRange
	}

	if upper < lower {
		return nil, ErrorInvalidValueRange
	}
	return ValueRange{def, byte(lower), byte(upper)}, nil
}

func MustValueRange(def string) Rule {
	rule, err := NewValueRange(def)
	if err != nil {
		panic(err)
	}
	return rule
}

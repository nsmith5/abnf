package rules

import (
	"errors"
	"strconv"
	"strings"
)

var ErrorInvalidCharacter = errors.New("Invalid character rule")

type Character struct {
	def  string
	char byte
}

func (c Character) Definition() string {
	return c.def
}

func (c Character) Children() []*Rule {
	return nil
}

func (c Character) Value() byte {
	return c.char
}

func NewCharacter(def string) (Rule, error) {
	var number uint64
	var err error

	switch {
	case strings.HasPrefix(def, `%x`):
		hex := strings.TrimPrefix(def, `%x`)
		number, err = strconv.ParseUint(hex, 16, 8)

	case strings.HasPrefix(def, `%d`):
		decimal := strings.TrimPrefix(def, `%d`)
		number, err = strconv.ParseUint(decimal, 10, 8)

	case strings.HasPrefix(def, `%b`):
		binary := strings.TrimPrefix(def, `%b`)
		number, err = strconv.ParseUint(binary, 2, 8)

	default:
		return nil, ErrorInvalidCharacter
	}
	if err != nil {
		return nil, ErrorInvalidCharacter
	}
	return Character{def, byte(number)}, nil
}

func MustCharacter(def string) Rule {
	rule, err := NewCharacter(def)
	if err != nil {
		panic(err)
	}
	return rule
}

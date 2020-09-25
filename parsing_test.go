package abnf

import (
	"testing"

	"github.com/nsmith5/abnf/rules"
)

func TestParse(t *testing.T) {
	var input [1 << 14]byte
	rule := rules.NewABNF()
	Parse(&rule, input[:])
}

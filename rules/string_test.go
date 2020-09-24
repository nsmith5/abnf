package rules

import (
	"testing"
)

func TestNewStringRule(t *testing.T) {
	for _, test := range []string{`"literal"`, `%d97.98.97`, `%s"aBa"`, `%i"aba"`} {
		_, err := NewString(test)
		if err != nil {
			t.Error(err)
		}
	}
}

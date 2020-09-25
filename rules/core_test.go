package rules

import (
	"fmt"
	"testing"
)

func TestCore(t *testing.T) {
	for _, rule := range Core {
		switch (*rule).(type) {
		case nil:
			t.Error("No Core rules should be nil type")
		}
		// Check that all rules in core rules are not nil
		fmt.Printf("(%v, %T)\n", *rule, *rule)
	}
}

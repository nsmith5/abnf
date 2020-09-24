package rules

import (
	"testing"
)

func TestValueRange(t *testing.T) {
	for _, test := range []string{`%x91-92`, `%d91-92`, `%b01101100-01101111`} {
		_, err := NewValueRange(test)
		if err != nil {
			t.Error(err)
		}
	}

	_, err := NewValueRange(`%d97-96`)
	if err != ErrorInvalidValueRange {
		t.Error("Expected invalid range")
	}
}

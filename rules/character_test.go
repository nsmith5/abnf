package rules

import (
	"testing"
)

func TestBadCharacter(t *testing.T) {
	// Just garbage
	_, err := NewCharacter(`lasdfjlasdf`)
	if err != ErrorInvalidCharacter {
		t.Error(err)
	}

	// Bad binary rule
	_, err = NewCharacter(`%b1098123`)
	if err != ErrorInvalidCharacter {
		t.Error(err)
	}

	// Bad hex rule
	_, err = NewCharacter(`%xdeaf012309`)
	if err != ErrorInvalidCharacter {
		t.Error(err)
	}

	// Bad decimal
	_, err = NewCharacter(`%d256`)
	if err != ErrorInvalidCharacter {
		t.Error(err)
	}
}

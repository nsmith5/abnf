package abnf

import "testing"

func TestALPHA(t *testing.T) {
	for _, v := range []string{"a", "z", "v", "A", "Z", "V", "Q"} {
		m := ALPHA([]byte(v))
		if m == nil {
			t.Error("All characters a-z and A-Z should match ALPHA")
		}
	}

	for _, v := range []string{"-", "8", "."} {
		m := ALPHA([]byte(v))
		if m != nil {
			t.Error("These characters shouldn't match ALPHA")
		}
	}
}

func TestBIT(t *testing.T) {
	for _, v := range []string{"0", "1"} {
		m := BIT([]byte(v))
		if m == nil {
			t.Error("Failed to match BIT")
		}
	}
}

func TestCHAR(t *testing.T) {
	var buf [1]byte
	var v byte
	for v = 0x01; v <= 0x7F; v++ {
		buf[0] = v
		m := CHAR(buf[:])
		if m == nil {
			t.Error("Failed to match CHAR")
		}
	}
}

func TestCRLF(t *testing.T) {
	m := CR([]byte("\r"))
	if m == nil {
		t.Error("Failed to match carriage return")
	}

	m = CRLF([]byte("\r\n"))
	if m == nil {
		t.Error("Failed to match carriage return line feed")
	}

	m = LF([]byte("\n"))
	if m == nil {
		t.Error("Failed to match line feed")
	}

	m = SP([]byte(" "))
	if m == nil {
		t.Error("Failed to match space")
	}
	m = HTAB([]byte("\t"))
	if m == nil {
		t.Error("Failed to match tab")
	}
}

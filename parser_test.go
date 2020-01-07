package abnf

import (
	"testing"
)

type TestCase struct {
	RuleName string
	Rule     Rule
	Examples []string
}

var validTests = []TestCase{
	TestCase{
		"ruleName",
		ruleName,
		[]string{`this`, `tha1a1231`, `fasd-123-adf-`},
	},
	TestCase{
		"definedAs",
		definedAs,
		[]string{` = `},
	},
	TestCase{
		`repetition`,
		repetition,
		[]string{
			`*"thisandthat"`,
			`1*100"thing"`,
			`*300%d10`,
		},
	},
	TestCase{
		`repeat`,
		repeat,
		[]string{
			`1234567890`,
			`1*`,
			`1*100`,
			`*4`,
		},
	},
	TestCase{
		`element`,
		element,
		[]string{
			`valid-rule-name`, // rule name
			`( %x01 )`,        // group
			`[ %x01 ]`,        // option
			`"charval"`,       // char value
			`%x12032`,         // numerical value
			`<abc>`,           // prose value
		},
	},
	TestCase{
		`commentOrWhiteSpace`,
		commentOrWhitespace,
		[]string{" ", "\t", "\r\n", "; this and that \r\n"},
	},
	TestCase{
		`commentOrNewline`,
		commentOrNewline,
		[]string{"\r\n", ";a\r\n"},
	},
	TestCase{`comment`, comment, []string{"; this and that\r\n"}},
	TestCase{
		`charVal`,
		charVal,
		[]string{
			`"abcdefghijklmnopqrstuvwzyz"`,
			`"ABCDEFGHIJKLMNOPQRSTUVWXYZ"`,
		},
	},
	TestCase{`numVal`, numVal, []string{`%b010001`, `%d1234123.123123.123`, `%x0344e.x45234`}},
	TestCase{`binVal`, binVal, []string{`b011101001101`}},
	TestCase{
		`decVal`,
		decVal,
		[]string{`d1234567890`, `d1235123`},
	},
	TestCase{
		`hexVal`,
		hexVal,
		[]string{`x1A`, `X0123456789ABCDEF`, `xabcdef`},
	},
	TestCase{
		`proseVal`,
		proseVal,
		[]string{`<abc>`},
	},
}

var invalidTests = []TestCase{
	TestCase{
		`ruleName`,
		ruleName,
		[]string{"1shitat", "-sdfasdf"},
	},
}

func TestValid(t *testing.T) {
	var m *Match
	for _, c := range validTests {
		for _, e := range c.Examples {
			m = c.Rule([]byte(e))
			if m == nil {
				t.Errorf("Failed to match rule %s against valid example '%s'", c.RuleName, e)
			}
		}
	}
}

func TestInValid(t *testing.T) {
	var m *Match
	for _, c := range invalidTests {
		for _, e := range c.Examples {
			m = c.Rule([]byte(e))
			if m != nil {
				t.Errorf("Shouldn't have matched rule %s against invalid example '%s'", c.RuleName, e)
			}
		}
	}
}

func TestRule(t *testing.T) {
	m := RuleList([]byte("this = %x61\r\n"))
	if m == nil {
		t.Error("What!")
	}
}

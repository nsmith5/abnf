package rules

type Concatenation struct {
	def   string
	rules []*Rule
}

func NewConcatenation(def string, rules ...*Rule) Rule {
	return Concatenation{def, rules}
}

func (c Concatenation) Definition() string {
	return c.def
}

func (c Concatenation) Children() []*Rule {
	return c.rules
}

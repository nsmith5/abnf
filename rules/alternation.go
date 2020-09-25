package rules

type Alternation struct {
	def   string
	rules []*Rule
}

func NewAlternation(def string, rules ...*Rule) Rule {
	return Alternation{def, rules}
}

func (a Alternation) Definition() string {
	return a.def
}

func (a Alternation) Children() []*Rule {
	return a.rules
}

package rules

type Repetition struct {
	def  string
	min  int
	max  int
	rule *Rule
}

func NewVariableRepetition(def string, min, max int, rule *Rule) Rule {
	return Repetition{def, min, max, rule}
}

func NewRepetition(def string, count int, rule *Rule) Rule {
	return Repetition{def, count, count, rule}
}

func NewOption(def string, rule *Rule) Rule {
	return Repetition{def, 0, 1, rule}
}

func (r Repetition) Definition() string {
	return r.def
}

func (r Repetition) Children() []*Rule {
	return []*Rule{r.rule}
}

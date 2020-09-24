package rules

// ABNF Grammar
var ABNF Rule

func init() {
	rules := [...]string{
		"rulelist",
		"rule",
		"rulename",
		"defined-as",
		"elements",
		// "c-wsp", No thank you!
		// "c-nl", Nope!
		// "comment", No comments
		"alternation",
		"concatenation",
		"repetition",
		"repeat",
		"element",
		"group",
		"option",
		"char-val",
		"num-val",
		"bin-val",
		"dec-val",
		"hex-val",
		// "prose-val", Not supported
	}

	var table = make(map[string]*Rule)

	// Add Core rules
	for k, v := range Core {
		table[k] = v
	}
	// Add all rule names
	for _, rule := range rules {
		table[rule] = new(Rule)
	}

	*table[`rulelist`] = NewVariableRepetition(`1*rule`, 1, -1, table[`rule`])
	*table[`rule`] = NewConcatenation(
		`rulename defined-as elements CRLF`,
		table[`rulename`],
		table[`defined-as`],
		table[`elements`],
		table[`CRLF`],
	)

	// Fuck how will I turn this into a function. Look at all the intermediate
	// steps.  Perhaps intermediate symbol tables? I'm sort of doing that right
	// here right?  Anyways, something to think about in the synthesis code I
	// guess.
	{
		a := MustString(`"-"`)
		b := NewAlternation(`ALPHA / DIGIT / "-"`, table[`ALPHA`], table[`DIGIT`], &a)
		c := NewVariableRepetition(`*(ALPHA / DIGIT / "-")`, 0, -1, &b)
		*table[`rulename`] = NewConcatenation(`ALPHA *(ALPHA / DIGIT / "-")`, table[`ALPHA`], &c)
	}

	{
		a := MustString(`"="`)
		*table[`defined-as`] = NewConcatenation(`SP "=" SP`, table[`SP`], &a, table[`SP`])
	}

	{
		a := MustString(`"/"`)
		b := NewConcatenation(`SP "/" SP concatenation`, table[`SP`], &a, table[`SP`], table[`concatenation`])
		c := NewVariableRepetition(`*(SP "/" SP concatenation)`, 0, -1, &b)
		*table[`alternation`] = NewConcatenation(
			`concatenation *(SP "/" SP concatenation)`,
			table[`concatenation`],
			&c,
		)
	}

	{
		a := NewConcatenation(`SP repetition`, table[`SP`], table[`concatenation`])
		b := NewVariableRepetition(`*(SP repetition)`, 0, -1, &a)
		*table[`concatenation`] = NewConcatenation(`repetition *(SP repetition)`, table[`repetition`], &b)
	}

	{
		a := NewOption(`[repeat]`, table[`repeat`])
		*table[`repeat`] = NewConcatenation(`[repeat] element`, &a, table[`element`])
	}
}

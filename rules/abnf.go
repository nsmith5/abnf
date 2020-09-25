package rules

func NewABNF() Rule {
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
		table[`alternation`],
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
		*table[`repetition`] = NewConcatenation(`[repeat] element`, &a, table[`element`])
	}

	{
		bb := MustString(`"*"`)
		ba := NewVariableRepetition(`*DIGIT`, 0, -1, table[`DIGIT`])
		b := NewConcatenation(`*DIGIT "*" *DIGIT`, &ba, &bb, &ba)
		a := NewVariableRepetition(`1*DIGIT`, 1, -1, table[`DIGIT`])
		*table[`repeat`] = NewAlternation(`1*DIGIT / (*DIGIT "*" *DIGIT)`, &a, &b)
	}

	*table[`element`] = NewAlternation(`rulename / group / option / char-val / num-val`, table[`rulename`], table[`group`], table[`char-val`], table[`num-val`])

	{
		a := MustString(`"("`)
		b := MustString(`")"`)
		*table[`group`] = NewConcatenation(`"(" SP alternation SP ")"`, &a, table[`SP`], table[`alternation`], table[`SP`], &b)
	}

	{
		a := MustString(`"["`)
		b := MustString(`"]"`)
		*table[`option`] = NewConcatenation(`"[" SP alternation SP "]"`, &a, table[`SP`], table[`alternation`], table[`SP`], &b)
	}

	{
		d := MustValueRange(`%x23-7E`)
		c := MustValueRange(`%x20-21`)
		b := NewAlternation(`%x20-21 / %x23-7E`, &c, &d)
		a := NewVariableRepetition(`*(%x20-21 / %x23-7E)`, 0, -1, &b)
		*table[`char-val`] = NewConcatenation(`DQUOTE *(%x20-21 / %x23-7E) DQUOTE`, table[`DQUOTE`], &a, table[`DQUOTE`])
	}

	{
		b := NewAlternation(`bin-val / dec-val / hex-val`, table[`bin-val`], table[`dec-val`], table[`hex-val`])
		a := MustString(`"%"`)
		*table[`num-val`] = NewConcatenation(`"%" (bin-val / dec-val / hex-val)`, &a, &b)
	}

	{
		cabb := NewVariableRepetition(`1*BIT`, 1, -1, table[`BIT`])
		caba := MustString(`"-"`)
		cab := NewConcatenation(`"-" 1*BIT`, &caba, &cabb)
		caaab := NewVariableRepetition(`1*BIT`, 1, -1, table[`BIT`])
		caaaa := MustString(`"."`)
		caaa := NewConcatenation(`"." 1*BIT`, &caaaa, &caaab)
		caa := NewVariableRepetition(`1*("." 1*BIT)`, 1, -1, &caaa)
		ca := NewAlternation(`1*("." 1*BIT) / ("-" 1*BIT)`, &caa, &cab)
		c := NewOption(`[ 1*("." 1*BIT) / ("-" 1*BIT) ]`, &ca)
		b := NewVariableRepetition(`1*BIT`, 1, -1, table[`BIT`])
		a := MustString(`"b"`)
		*table[`bin-val`] = NewConcatenation(`"b" 1*BIT [ 1*("." 1*BIT) / ("-" 1*BIT) ]`, &a, &b, &c)
	}

	{
		cabb := NewVariableRepetition(`1*DIGIT`, 1, -1, table[`DIGIT`])
		caba := MustString(`"-"`)
		cab := NewConcatenation(`"-" 1*DIGIT`, &caba, &cabb)
		caaab := NewVariableRepetition(`1*DIGIT`, 1, -1, table[`DIGIT`])
		caaaa := MustString(`"."`)
		caaa := NewConcatenation(`"." 1*DIGIT`, &caaaa, &caaab)
		caa := NewVariableRepetition(`1*("." 1*DIGIT)`, 1, -1, &caaa)
		ca := NewAlternation(`1*("." 1*DIGIT) / ("-" 1*DIGIT)`, &caa, &cab)
		c := NewOption(`[ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]`, &ca)
		b := NewVariableRepetition(`1*DIGIT`, 1, -1, table[`DIGIT`])
		a := MustString(`"d"`)
		*table[`dec-val`] = NewConcatenation(`"d" 1*DIGIT [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]`, &a, &b, &c)
	}

	{
		cabb := NewVariableRepetition(`1*HEXDIG`, 1, -1, table[`HEXDIG`])
		caba := MustString(`"-"`)
		cab := NewConcatenation(`"-" 1*HEXDIG`, &caba, &cabb)
		caaab := NewVariableRepetition(`1*HEXDIG`, 1, -1, table[`HEXDIG`])
		caaaa := MustString(`"."`)
		caaa := NewConcatenation(`"." 1*HEXDIG`, &caaaa, &caaab)
		caa := NewVariableRepetition(`1*("." 1*HEXDIG)`, 1, -1, &caaa)
		ca := NewAlternation(`1*("." 1*HEXDIG) / ("-" 1*HEXDIG)`, &caa, &cab)
		c := NewOption(`[ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]`, &ca)
		b := NewVariableRepetition(`1*HEXDIG`, 1, -1, table[`HEXDIG`])
		a := MustString(`"x"`)
		*table[`hex-val`] = NewConcatenation(`"x" 1*HEXDIG [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]`, &a, &b, &c)
	}

	return *table[`rulelist`]
}

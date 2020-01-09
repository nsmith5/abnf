package abnf

// group = "(" *c-wsp alternation *c-wsp ")"
func group(input []byte) *Match {
	return Sequence(
		`group`,
		String("(", "(", false),
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		alternation,
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		String(")", ")", false),
	)(input)
}

// option = "[" *c-wsp alternation *c-wsp "]"
func option(input []byte) *Match {
	return Sequence(
		`option`,
		String(`[`, `[`, false),
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		alternation,
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		String(`]`, `]`, false),
	)(input)
}

// alternation = concatenation *(*c-wsp "/" *c-wsp concatenation)
func alternation(input []byte) *Match {
	return Sequence(
		"alternation",
		concatenation,
		Repeat(
			`*(*c-wsp "/" *c-wsp concatenation)`,
			-1,
			-1,
			Sequence(
				`*c-wsp "/" *c-wsp concatenation`,
				Repeat("*c-wsp", -1, -1, commentOrWhitespace),
				String("/", "/", false),
				Repeat("*c-wsp", -1, -1, commentOrWhitespace),
				concatenation,
			),
		),
	)(input)
}

// concatenation = repetition *(1*c-wsp repetition)
func concatenation(input []byte) *Match {
	return Sequence(
		`concatenation`,
		repetition,
		Repeat(
			`*(1*c-wsp repetition)`,
			-1,
			-1,
			Sequence(
				`1*c-wsp repetition`,
				Repeat(`1*c-wsp`, 1, -1, commentOrWhitespace),
				repetition,
			),
		),
	)(input)
}

// repetition = [repeat] element
func repetition(input []byte) *Match {
	return Sequence(`repetition`, Option(`[repeat]`, repeat), element)(input)
}

// repeat = 1*DIGIT / (*DIGIT "*" *DIGIT)
func repeat(input []byte) *Match {
	return Or(
		`repeat`,
		Repeat(`1*DIGIT`, 1, -1, DIGIT),
		Sequence(
			`*DIGIT "*" *DIGIT`,
			Repeat(`*DIGIT`, -1, -1, DIGIT),
			String("*", "*", true),
			Repeat(`*DIGIT`, -1, -1, DIGIT),
		),
	)(input)
}

// element = rulename / group / option / char-val / num-val / prose-val
func element(input []byte) *Match {
	return Or(`element`, ruleName, group, option, charVal, numVal, proseVal)(input)
}

// Rules of the ABNF grammar
var (
	// rulelist = 1*( rule / (*WSP c-nl) )
	// NOTE: this takes into account Errata 3076
	// https://www.rfc-editor.org/errata/eid3076
	RuleList = Repeat(
		`rule-list`,
		1,
		-1,
		Or(
			`rule-list`,
			rule,
			Sequence(
				`(*c-wsp c-nil)`,
				Repeat(`*WSP`, -1, -1, WSP),
				commentOrNewline,
			),
		),
	)

	// rule = rulename defined-as elements c-nl
	rule = Sequence("rule", ruleName, definedAs, elements, commentOrNewline)

	// rulename = ALPHA *(ALPHA / DIGIT / "-")
	ruleName = Sequence(
		`rulename`,
		ALPHA,
		Repeat(
			`*(ALPHA / DIGIT / "-")`,
			-1,
			-1,
			Or(`ALPHA / DIGIT / "-"`, ALPHA, DIGIT, String("-", "-", true)),
		),
	)

	// defined-as = *c-wsp ("=" / "=/") *c-wsp
	definedAs = Sequence(
		`defined-as`,
		Repeat("*c-wsp", -1, -1, commentOrWhitespace),
		Or(
			`("=" / "=/")`,
			String("=", "=", true),
			String("=/", "=/", true),
		),
		Repeat("*c-wsp", -1, -1, commentOrWhitespace),
	)

	// elements = alternation *WSP
	// NOTE: this takes into account Errata 2968
	// https://www.rfc-editor.org/errata/eid2968
	elements = Sequence("elements", alternation, Repeat("*WSP", -1, -1, WSP))

	// c-wsp = WSP / (c-nl WSP)
	commentOrWhitespace = Or("c-wsp", WSP, commentOrNewline)

	// c-nl = comment / CRLF
	commentOrNewline = Or("c-nl", comment, CRLF)

	//  comment = ";" *(WSP / VCHAR) CRLF
	comment = Sequence(
		"comment",
		String(";", ";", true),
		Repeat(
			`*(WSP / VCHAR)`,
			-1,
			-1,
			Or(
				`(WSP / VCHAR)`,
				WSP,
				VCHAR,
			),
		),
		CRLF,
	)
	// char-val = DQUOTE *(%x20-21 / %x23-7E) DQUOTE
	charVal = Sequence(
		`char-val`,
		DQUOTE,
		Repeat(
			`*(%x20-21 / %x23-7E)`,
			-1,
			-1,
			Or(
				`%x20-21 / %x23-7E`,
				ByteRange(`%x20-21`, 0x20, 0x21),
				ByteRange(`%x23-7E`, 0x23, 0x7E),
			),
		),
		DQUOTE,
	)

	//num-val =  "%" (bin-val / dec-val / hex-val)
	numVal = Sequence(
		`num-val`,
		String(`%`, `%`, true),
		Or(`bin-val / dec-val / hex-val`, binVal, decVal, hexVal),
	)

	// bin-val  =  "b" 1*BIT
	//             [ 1*("." 1*BIT) / ("-" 1*BIT) ]
	//             ; series of concatenated bit values
	//             ;  or single ONEOF range
	binVal = Sequence(
		`bin-val`,
		String(`b`, `b`, false),
		Repeat(`1*BIT`, 1, -1, BIT),
		Option(
			`[1*("." 1*BIT) / ("-" 1*BIT)]`,
			Or(
				`1*("." 1*BIT) / ("-" 1*BIT)`,
				Repeat(
					`1*("." 1*BIT)`,
					1,
					-1,
					Sequence(`"." 1*BIT`, String(`.`, `.`, true), Repeat(`1*BIT`, 1, -1, BIT)),
				),
				Sequence(`"-" 1*BIT`, String(`-`, `-`, true), Repeat(`1*BIT`, 1, -1, BIT)),
			),
		),
	)

	// dec-val  =  "d" 1*DIGIT
	//             [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]
	decVal = Sequence(
		`dec-val`,
		String(`d`, `d`, false),
		Repeat(`1*DIGIT`, 1, -1, DIGIT),
		Option(
			`[1*("." 1*DIGIT) / ("-" 1*DIGIT)]`,
			Or(
				`1*("." 1*DIGIT) / ("-" 1*DIGIT)`,
				Repeat(
					`1*("." 1*DIGIT)`,
					1,
					-1,
					Sequence(`"." 1*DIGIT`, String(`.`, `.`, true), Repeat(`1*DIGIT`, 1, -1, DIGIT)),
				),
				Sequence(`"-" 1*DIGIT`, String(`-`, `-`, true), Repeat(`1*DIGIT`, 1, -1, DIGIT)),
			),
		),
	)

	// hex-val  =  "x" 1*HEXDIG
	//             [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
	hexVal = Sequence(
		`hex-val`,
		String(`x`, `x`, false),
		Repeat(`1*HEXDIG`, 1, -1, HEXDIG),
		Option(
			`[1*("." 1*HEXDIG) / ("-" 1*HEXDIG)]`,
			Or(
				`1*("." 1*HEXDIG) / ("-" 1*HEXDIG)`,
				Repeat(
					`1*("." 1*HEXDIG)`,
					1,
					-1,
					Sequence(`"." 1*HEXDIG`, String(`.`, `.`, true), Repeat(`1*HEXDIG`, 1, -1, HEXDIG)),
				),
				Sequence(`"-" 1*HEXDIG`, String(`-`, `-`, true), Repeat(`1*HEXDIG`, 1, -1, HEXDIG)),
			),
		),
	)

	//	prose-val =  "<" *(%x20-3D / %x3F-7E) ">"
	//               ; bracketed string of SP and VCHAR
	//               ;  without angles
	//               ; prose description, to be used as
	//               ;  last resort
	proseVal = Sequence(
		`prose-val`,
		String(`<`, `<`, true),
		Repeat(
			`*(%x20-3D / %x3F-7E)`,
			-1,
			-1,
			Or(
				`%x20-3D / %x3F-7E`,
				ByteRange(`%x20-3D`, 0x20, 0x3D),
				ByteRange(`%x3F-7E`, 0x3F, 0x7E),
			),
		),
		String(`>`, `>`, true),
	)
)

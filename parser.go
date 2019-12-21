package main

// CORE rules for the ABNF grammar
var (
	ALPHA  = Or("ALPHA", ByteRange("A-Z", 'A', 'Z'), ByteRange("a-z", 'a', 'z'))
	BIT    = Or("BIT", Byte("0", '0'), Byte("1", '0'))
	CHAR   = ByteRange("CHAR", 0x01, 0x7F)
	CR     = Byte("CR", 0x0D)
	CRLF   = Sequence("CRLF", CR, LF)
	CTL    = Or("CTL", ByteRange("NULL - US", 0x00, 0x1F), Byte("DEL", 0x7F))
	DIGIT  = ByteRange("DIGIT", '0', '9')
	DQUOTE = Byte("DQUOTE", '"')
	HEXDIG = Or(
		"HEXDIG",
		DIGIT,
		String("A", "A", false),
		String("B", "B", false),
		String("C", "C", false),
		String("D", "D", false),
		String("E", "E", false),
		String("F", "F", false),
	)
	HTAB  = Byte("HTAB", 0x09)
	LF    = Byte("LF", 0x0A)
	LWSP  = Repeat("LWSP", -1, -1, Or("LWSP", WSP, Sequence("LWSP", CRLF, WSP)))
	OCTET = ByteRange("OCTECT", 0x00, 0xFF)
	SP    = Byte("SP", ' ')
	VCHAR = ByteRange("VCHAR", 0x21, 0x7E)
	WSP   = Or("WSP", SP, HTAB)
)

// Rules of the ABNF grammar
var (
	// rulelist = 1*( rule / (*c-wsp c-nl) )
	ruleList = Repeat(
		`rule-list`,
		1,
		-1,
		Or(
			`rule-list`,
			rule,
			Sequence(
				`(*c-wsp c-nil)`,
				Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
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
			`rulename`,
			-1,
			-1,
			Or("rulename", ALPHA, DIGIT, String("-", "-", false)),
		),
	)

	// defined-as = *c-wsp ("=" / "=/") *c-wsp
	definedAs = Sequence(
		`defined-as`,
		Repeat("*c-wsp", -1, -1, commentOrWhitespace),
		Or(
			`("=" / "=/")`,
			String("=", "=", false),
			String("=/", "=/", false),
		),
	)

	// elements = alternation *c-wsp
	elements = Sequence("elements", alternation, Repeat("*c-wsp", -1, -1, commentOrWhitespace))

	// c-wsp = WSP / (c-nl WSP)
	commentOrWhitespace = Or("c-wsp", WSP, commentOrNewline)

	// c-nl = comment / CRLF
	commentOrNewline = Or("c-nl", comment, CRLF)

	//  comment = ";" *(WSP / VCHAR) CRLF
	comment = Sequence(
		"comment",
		String(";", ";", false),
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

	// alternation = concatenation *(*c-wsp "/" *c-wsp concatenation)
	alternation = Sequence(
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
	)

	// concatenation = repetition *(1*c-wsp repetition)
	concatenation = Sequence(
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
	)

	// repetition = [repeat] element
	repetition = Sequence(`repetition`, Option(`[repeat]`, repeat), element)

	// repeat = 1*DIGIT / (*DIGIT "*" *DIGIT)
	repeat = Or(
		`repeat`,
		Repeat(`1*DIGIT`, 1, -1, DIGIT),
		Sequence(
			`*DIGIT "*" *DIGIT`,
			Repeat(`*DIGIT`, -1, -1, DIGIT),
			String("*", "*", false),
			Repeat(`*DIGIT`, -1, -1, DIGIT),
		),
	)

	// element = rulename / group / option / char-val / num-val / prose-val
	element = Or(`element`, ruleName, group, option, charVal, numVal, proseVal)

	// group = "(" *c-wsp alternation *c-wsp ")"
	group = Sequence(
		`group`,
		String("(", "(", false),
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		alternation,
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		String(")", ")", false),
	)

	// option = "[" *c-wsp alternation *c-wsp "]"
	option = Sequence(
		`option`,
		String(`[`, `[`, false),
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		alternation,
		Repeat(`*c-wsp`, -1, -1, commentOrWhitespace),
		String(`]`, `]`, false),
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
		String(`%`, `%`, false),
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
					Sequence(`"." 1*BIT`, String(`.`, `.`, false), Repeat(`1*BIT`, 1, -1, BIT)),
				),
				Sequence(`"-" 1*BIT`, String(`-`, `-`, false), Repeat(`1*BIT`, 1, -1, BIT)),
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
					Sequence(`"." 1*DIGIT`, String(`.`, `.`, false), Repeat(`1*DIGIT`, 1, -1, DIGIT)),
				),
				Sequence(`"-" 1*DIGIT`, String(`-`, `-`, false), Repeat(`1*DIGIT`, 1, -1, DIGIT)),
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
					Sequence(`"." 1*HEXDIG`, String(`.`, `.`, false), Repeat(`1*HEXDIG`, 1, -1, HEXDIG)),
				),
				Sequence(`"-" 1*HEXDIG`, String(`-`, `-`, false), Repeat(`1*HEXDIG`, 1, -1, HEXDIG)),
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
		String(`<`, `<`, false),
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
		String(`>`, `>`, false),
	)
)

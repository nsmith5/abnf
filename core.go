package abnf

// CORE rules for the ABNF grammar
var (
	ALPHA  = Or("ALPHA", ByteRange("A-Z", 'A', 'Z'), ByteRange("a-z", 'a', 'z'))
	BIT    = Or("BIT", Byte("0", '0'), Byte("1", '1'))
	CHAR   = ByteRange("CHAR", 0x01, 0x7F)
	CR     = Byte("CR", 0x0D)
	CRLF   = Sequence("CRLF", CR, LF)
	CTL    = Or("CTL", ByteRange("NULL - US", 0x00, 0x1F), Byte("DEL", 0x7F))
	DIGIT  = ByteRange("DIGIT", '0', '9')
	DQUOTE = Byte("DQUOTE", '"')
	HEXDIG = Or(
		"HEXDIG",
		DIGIT,
		ByteRange("A-F", 'A', 'F'),
		ByteRange("a-f", 'a', 'f'),
	)
	HTAB  = Byte("HTAB", 0x09)
	LF    = Byte("LF", 0x0A)
	LWSP  = Repeat("LWSP", -1, -1, Or("LWSP", WSP, Sequence("LWSP", CRLF, WSP)))
	OCTET = ByteRange("OCTECT", 0x00, 0xFF)
	SP    = Byte("SP", ' ')
	VCHAR = ByteRange("VCHAR", 0x21, 0x7E)
	WSP   = Or("WSP", SP, HTAB)
)

package rules

// Core is a symbol table of the core rules defined in ABNF
var Core = map[string]*Rule{}

var ruleNames = [...]string{
	"ALPHA",
	"BIT",
	"CHAR",
	"CR",
	"CRLF",
	"CTL",
	"DIGIT",
	"DQUOTE",
	"HEXDIG",
	"HTAB",
	"LF",
	// "LWSP",  Not implimented
	"OCTET",
	"SP",
	"VCHAR",
	"WSP",
}

func init() {
	// Fill up the symbol table first. This is helpful because it
	// means rules can reference other rules that haven't been defined
	// yet.
	for _, ruleName := range ruleNames {
		Core[ruleName] = new(Rule)
	}

	// ALPHA  =  %x41-5A / %x61-7A  ; A-Z / a-z
	lowercase := MustValueRange("%x41-5A")
	uppercase := MustValueRange("%x61-7A")
	*Core[`ALPHA`] = NewAlternation("%x41-5A / %x61-7A", &lowercase, &uppercase)

	// BIT = "0" / "1"
	one := MustString(`"1"`)
	zero := MustString(`"0"`)
	*Core[`BIT`] = NewAlternation(`"0" / "1"`, &zero, &one)

	// CHAR = %x01-7F
	char := MustValueRange(`%x01-7F`)
	*Core[`CHAR`] = char

	// CR = %x0D ; carriage return
	cr := MustCharacter(`%x0D`)
	*Core[`CR`] = cr

	// CRLF = CR LF ; internet line ending \r\n
	*Core[`CRLF`] = NewConcatenation(`CR LF`, Core[`CR`], Core[`LF`])

	// CTL = %x00-1F / %x7F ; control characters
	lowerctl := MustValueRange(`%x00-1F`)
	del := MustCharacter(`%x7F`)
	*Core[`CTL`] = NewAlternation(`%x00-1F / %x7F`, &lowerctl, &del)

	// DIGIT = %x30-39 ; 0-9
	digit := MustValueRange(`%x30-39`)
	*Core[`DIGIT`] = digit

	// DQUOTE = %x22
	dquote := MustCharacter(`%x22`)
	*Core[`DQUOTE`] = dquote

	// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
	a := MustString(`"A"`)
	b := MustString(`"B"`)
	c := MustString(`"C"`)
	d := MustString(`"D"`)
	e := MustString(`"E"`)
	f := MustString(`"F"`)
	*Core[`HEXDIG`] = NewAlternation(`DIGIT / "A" / "B" / "C" / "D" / "E" / "F"`, Core[`DIGIT`], &a, &b, &c, &d, &e, &f)

	// HTAB = %x09
	htab := MustCharacter(`%x09`)
	*Core[`HTAB`] = htab

	// LF = %x0A
	lf := MustCharacter(`%x0A`)
	*Core[`LF`] = lf

	// LWSP = *(WSP / CRLF WSP)
	// ABNF has cautioned against even using this one so we leave it out!

	// OCTET = %x00-FF
	octet := MustValueRange(`%x00-FF`)
	*Core[`OCTET`] = octet

	// SP = %x20
	space := MustCharacter(`%x20`)
	*Core[`SP`] = space

	// VCHAR = %x21-7E ; printable characters
	vchar := MustValueRange(`%x21-7E`)
	*Core[`VCHAR`] = vchar

	// WSP = SP / HTAB ; whitespace
	*Core[`WSP`] = NewAlternation(`SP / HTAP`, Core[`SP`], Core[`HTAB`])
}

package parse

// Lex the drawing language tokens
// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	errInvalidUnicodeEscapeMsg = "Invalid unicode escape sequence %s: must be \\uXXXX, \\uXXXXXX, \\U+XXXX, or \\U+XXXXXX"
	errInvalidEscapeMsg        = "Invalid escape sequence %s: must be \\\\, \\', \\n, \\uXXXX, \\uXXXXXX, \\U+XXXX, or \\U+XXXXXX"
	errInvalidColourMsg        = "Invalid colour %s: there must be six hex characters after the #"
	errIncompleteFloatMsg      = "Incomplete float number %s: a float cannot end with a ., e, or E"
	errNameTooLongMsg          = "Name too long %q: a name can be a max of 16 chars"
	errIllegalStringCharMsg    = "Illegal string %q: a string cannot contain ASCII control characters except for \r and \n"
	errUnexpectedEOF           = fmt.Errorf("Unexpected EOF")
)

// TokenType describes the types of tokens to lex
type TokenType uint

const (
	Eol TokenType = iota
	Percent
	AssignModulus
	OParens
	CParens
	Star
	AssignMultiply
	Plus
	AssignAdd
	Increment
	Comma
	Minus
	AssignSubtract
	Decrement
	Slash
	AssignDivide
	Colon
	LessThan
	Equals
	GreaterThan
	OBracket
	CBracket
	OBrace
	CBrace
	Eof
	Undefined

	Colour
	FloatNumber
	IntNumber
	Name
	Str
)

// Constants for tokens that are always the same sequence of runes
var (
	cEol            = LexToken{Eol, "\n"}
	cPercent        = LexToken{Percent, "%"}
	cAssignModulus  = LexToken{AssignModulus, "%="}
	cOParens        = LexToken{OParens, "("}
	cCParens        = LexToken{CParens, ")"}
	cStar           = LexToken{Star, "*"}
	cAssignMultiply = LexToken{AssignMultiply, "*="}
	cPlus           = LexToken{Plus, "+"}
	cAssignAdd      = LexToken{AssignAdd, "+="}
	cIncrement      = LexToken{Increment, "++"}
	cComma          = LexToken{Comma, ","}
	cMinus          = LexToken{Minus, "-"}
	cAssignSubtract = LexToken{AssignSubtract, "-="}
	cDecrement      = LexToken{Decrement, "--"}
	cSlash          = LexToken{Slash, "/"}
	cAssignDivide   = LexToken{AssignDivide, "/="}
	cColon          = LexToken{Colon, ":"}
	cLessThan       = LexToken{LessThan, "<"}
	cEquals         = LexToken{Equals, "="}
	cGreaterThan    = LexToken{GreaterThan, ">"}
	cOBracket       = LexToken{OBracket, "["}
	cCBracket       = LexToken{CBracket, "]"}
	cOBrace         = LexToken{OBrace, "{"}
	cCBrace         = LexToken{CBrace, "}"}
	cEof            = LexToken{Eof, ""}
	cUndefined      = LexToken{Undefined, ""}
)

// LexToken describes a single token, as a TokenType and a string of characters
type LexToken struct {
	TokenType
	Token string
}

// IntValue returns the integer value for Colour and IntNumber token types
func (l LexToken) IntValue() uint64 {
	// Remove any prefix thaat might be included in the token
	var (
		str  = l.Token
		base = 10
	)

	switch {
	case strings.HasPrefix(l.Token, "0b"):
		str = str[2:]
		base = 2

	case strings.HasPrefix(l.Token, "0x"):
		str = str[2:]
		base = 16

	case strings.HasPrefix(l.Token, "#"):
		str = str[1:]
		base = 16
	}

	// Replace all _ with empty string
	str = strings.ReplaceAll(str, "_", "")

	// Convert to a uint64
	val, err := strconv.ParseUint(str, base, 64)
	if err != nil {
		panic(err.(*strconv.NumError).Err)
	}

	return val
}

func (l LexToken) FloatValue() float32 {
	// No prefix, straightforward read of string
	val, err := strconv.ParseFloat(l.Token, 32)
	if err != nil {
		panic(err.(*strconv.NumError).Err)
	}

	return float32(val)
}

// nextRune returns the next rune from the input.
// eof results in 0; panics on any other error.
func nextRune(src io.RuneScanner) rune {
	// Get next char, we don't care how many bytes it takes
	r, _, err := src.ReadRune()

	// Error handling
	if err == io.EOF {
		return 0
	} else if err != nil {
		panic(err)
	}

	return r
}

// Helper function to determine if a char is a hex char, and if so, what is the value of it from 0 to 15
func hexVal(r rune) (uint64, bool) {
	switch {
	case ((r >= '0') && (r <= '9')):
		return uint64(r - '0'), true
	case ((r >= 'A') && (r <= 'F')):
		return uint64(r - 'A' + 10), true
	case ((r >= 'a') && (r <= 'f')):
		return uint64(r - 'a' + 10), true
	default:
		return uint64(r), false
	}
}

// Helper function to read 4 or 6 hex chars that specify a unicode char
// Have already read prefix of \u or \U+
func unicodeHex(prefix string, src io.RuneScanner) rune {
	var (
		res   uint64
		r     rune
		chars = prefix
	)

	// Has to have at least 4 hex chars
	for i := 0; i < 4; i++ {
		r = nextRune(src)
		v, haveIt := hexVal(r)
		if !haveIt {
			panic(fmt.Errorf(errInvalidUnicodeEscapeMsg, chars))
		}

		chars += string(r)
		res = res*16 + v
	}

	// May be 6 hex chars
	r = nextRune(src)
	v, haveIt := hexVal(r)
	if !haveIt {
		// Not a hex char, unread it and return unicode char
		src.UnreadRune()
		return rune(res)
	}

	// Have 5 hex chars
	chars += string(r)
	res = res*16 + v

	// Must have one more hex char
	r = nextRune(src)
	v, haveIt = hexVal(r)
	if !haveIt {
		panic(fmt.Errorf(errInvalidUnicodeEscapeMsg, chars))
	}

	// Have 6 hex chars, return unicode char
	chars += string(r)
	res = res*16 + v
	return rune(res)
}

// Helper function for literal strings
// Read a unicode char, or escape sequence
// A unicode char is any char from space onwards except for DEL
//
// Escapes can be used for quotes and newlines:
// - \\ for an actual backslash
// - \' for an escaped '
// - \n for an escaped eol
//
// Escapes can be used for non-ASCII unicode chars:
// - \u  [0-9A-Fa-f]{4}
// - \u  [0-9A-Fa-f]{6}
// - \U+ [0-9A-Fa-f]{4}
// - \U+ [0-9A-Fa-f]{6}
//
// Note that \r is not allowed, only \n
// Returns resulting char and true if it was the result of an escape sequence
// The bool allows the caller to differentiate between an escaped or unescaped quote char
func escapedChar(src io.RuneScanner) (rune, bool) {
	r := nextRune(src)
	if r == '\\' {
		switch r = nextRune(src); r {
		case '\\': // \\ = \
			return r, true
		case '\'': // \' = '
			return r, true
		case 'n': // \n = newline
			return '\n', true
		case 'u': // \u needs 4 or 6 hex chars
			return unicodeHex("\\u", src), true
		case 'U': // \U needs a + followed by 4 or 6 hex chars
			if r = nextRune(src); r != '+' {
				panic(fmt.Errorf(errInvalidUnicodeEscapeMsg, "\\U"+string(r)))
			}
			return unicodeHex("\\U+", src), true
		default:
			panic(fmt.Errorf(errInvalidEscapeMsg, fmt.Sprintf("\\%s", string(r))))
		}
	}

	return r, false
}

// Helper function to read a single quoted string
// Single quoted strings end with an unescaped single quote, and can have escaped or embedded newlines
func readString(src io.RuneScanner) LexToken {
	var str strings.Builder
	str.WriteRune('\'')

	for {
		r, escaped := escapedChar(src)
		str.WriteRune(r)
		if (r < ' ') && ((r != '\r') && (r != '\n')) {
			if r == 0 {
				panic(errUnexpectedEOF)
			}
			panic(fmt.Errorf(errIllegalStringCharMsg, str.String()))
		}

		if (r == '\'') && (!escaped) {
			// Complete string
			return LexToken{Str, str.String()}
		}
	}

	// All switch cases in above for loop return or panic, so this line can never be reached
}

// Helper function to read a binary number of 0, 1, and _
func readBinaryNumber(src io.RuneScanner) LexToken {
	var str strings.Builder

	str.WriteRune('0')
	str.WriteRune('b')

	for {
		r := nextRune(src)

		switch {
		case (r == '0') || (r == '1'):
			str.WriteRune(r)

		case r == '_': // separator, ignore it as far as the value goes
			str.WriteRune(r)

		default:
			// first char of next token
			src.UnreadRune()
			return LexToken{IntNumber, str.String()}
		}
	}
}

// Helper function to read a hex number of hex digits and _
func readHexNumber(src io.RuneScanner) LexToken {
	var str strings.Builder

	str.WriteRune('0')
	str.WriteRune('x')

	for {
		r := nextRune(src)
		_, haveIt := hexVal(r)

		switch {
		case haveIt:
			str.WriteRune(r)

		case r == '_': // separator, ignore it as far as the value goes
			str.WriteRune(r)

		default:
			// first char of next token
			src.UnreadRune()
			return LexToken{IntNumber, str.String()}
		}
	}
}

// Helper function to read a decimal number, which may be an integer or float
// the mantissa may have _
func readDecimalNumber(firstDigit rune, src io.RuneScanner) LexToken {
	var str strings.Builder

	str.WriteRune(firstDigit)

	for {
		r := nextRune(src)

		switch {
		case (r >= '0') && (r <= '9'):
			str.WriteRune(r)

		case r == '_': // separator, ignore it as far as the value goes
			str.WriteRune(r)

		case (r == '.') || (r == 'e') || (r == 'E'): // change to float mode
			str.WriteRune(r)
			return readFloatNumber(&str, r, src)

		default:
			// first char of next token
			src.UnreadRune()
			return LexToken{IntNumber, str.String()}
		}
	}
}

// Helper function to read a float number
// We started as a decimal number, then we just read  a ., e, or E
func readFloatNumber(str *strings.Builder, r rune, src io.RuneScanner) LexToken {
	var (
		// 0: after ., before first digit
		// 1: digits after .
		// 2: after e, before first exponent digit
		// 3: after first exponent digit
		mode = 0
	)

	if r != '.' {
		mode = 2 // Only other chars are e and E
	}

	for {
		r = nextRune(src)

		switch {
		case (r >= '0') && (r <= '9'):
			str.WriteRune(r)

			switch mode {
			case 0:
				// first digit after .
				mode = 1
			case 2:
				// first digit after e
				mode = 3
			}

		case (r == 'e') || (r == 'E'):
			switch mode {
			case 0:
				// After ., we need a digit, not an e
				str.WriteRune(r)
				panic(fmt.Errorf(errIncompleteFloatMsg, str.String()))

			case 1:
				// Already read digits after ., switching to exponent
				str.WriteRune(r)
				mode = 2

			case 2:
				// After an e, we need a digit, not another e
				panic(fmt.Errorf(errIncompleteFloatMsg, str.String()))

			default:
				// After e and digits, first char of next token
				src.UnreadRune()
				return LexToken{FloatNumber, str.String()}
			}

		default:
			// Not a float char
			switch mode {
			case 0:
				// After ., we need a digit
				if r != 0 {
					str.WriteRune(r)
				}
				panic(fmt.Errorf(errIncompleteFloatMsg, str.String()))

			case 1:
				// After . and digits, first char of next token
				src.UnreadRune()
				return LexToken{FloatNumber, str.String()}

			case 2:
				// After e, we need a digit
				if r != 0 {
					str.WriteRune(r)
				}
				panic(fmt.Errorf(errIncompleteFloatMsg, str.String()))

			default:
				// After e and digits, first char of next token
				src.UnreadRune()
				return LexToken{FloatNumber, str.String()}
			}
		}
	}
}

// Lex lexes the next token in the given RuneScanner.
// Some tokens are negatively identified by stopping at a character that is not part of the token,
// so the RuneScanner is used to be able to unread a single rune as the first char of the next token.
// Whitespace is skipped, except for newlines that are preserved, since they are significant in the parsing.
// All newline sequences are coalesced into a Unix newline, for simplicity.
func Lex(src io.RuneScanner) LexToken {
	// Get next rune
	r := nextRune(src)

	// EOF handling
	if r == 0 {
		return cEof
	}

	// Lex a complete token, that is longest match
	switch {
	case r == '\n':
		// unix eol
		return cEol

	case r == '\r':
		// if next rune is \n, windows \r\n
		if r = nextRune(src); r != '\n' {
			// otherwise, mac \r by itself
			src.UnreadRune()
		}
		return cEol

	case r == '#':
		// colour, needs 6 hex digits
		var str strings.Builder
		str.WriteRune('#')
		for i := 0; i < 6; i++ {
			r := nextRune(src)
			str.WriteRune(r)
			_, haveIt := hexVal(r)
			if !haveIt {
				panic(fmt.Errorf(errInvalidColourMsg, str.String()))
			}
		}
		return LexToken{Colour, str.String()}

	case r == '%':
		// Could be % or %=
		switch r = nextRune(src); r {
		case '=': // %=
			return cAssignModulus
		default: // %
			src.UnreadRune()
			return cPercent
		}

	case r == '\'':
		// string, read all until next unescaped ", interpreting escapes, and allowing embedded newlines
		return readString(src)

	case r == '(':
		return cOParens

	case r == ')':
		return cCParens

	case r == '*':
		// Could be * or *=
		switch r = nextRune(src); r {
		case '=': // *=
			return cAssignMultiply
		default: // *
			src.UnreadRune()
			return cStar
		}

	case r == '+':
		// Could be +, +=, or ++
		switch r = nextRune(src); r {
		case '=': // +=
			return cAssignAdd
		case '+': // ++
			return cIncrement
		default: // +
			src.UnreadRune()
			return cPlus
		}

	case r == ',':
		return cComma

	case r == '-':
		// Could be -, -=, or --
		switch r = nextRune(src); r {
		case '=': // -=
			return cAssignSubtract
		case '-': // --
			return cDecrement
		default: // -
			src.UnreadRune()
			return cMinus
		}

	case r == '/':
		// Could be / or /=
		switch r = nextRune(src); r {
		case '=': // /=
			return cAssignDivide
		default: // /
			src.UnreadRune()
			return cSlash
		}

	case r == ':':
		return cColon

	case r == '<':
		return cLessThan

	case r == '=':
		return cEquals

	case r == '>':
		return cGreaterThan

	case r == '[':
		return cOBracket

	case r == ']':
		return cCBracket

	case r == '{':
		return cOBrace

	case r == '}':
		return cCBrace

	case r == '0':
		r = nextRune(src)
		switch {
		case r == 'b': // binary number, read all 0, 1, and _
			return readBinaryNumber(src)

		case r == 'x': // hex number, read all hex and _
			return readHexNumber(src)

		case (r >= '0') && (r <= '9'): // decimal with leading 0
			// Unread char after leading 0
			src.UnreadRune()
			// Pass leading 0 as prefix
			return readDecimalNumber('0', src)
		}

	case (r >= '1') && (r <= '9'):
		return readDecimalNumber(r, src)

	case ((r >= 'A') && (r <= 'Z')) || ((r >= 'a') && (r <= 'z')):
		var str strings.Builder
		str.WriteRune(r)
		for {
			if r = nextRune(src); ((r >= 'A') && (r <= 'Z')) || ((r >= 'a') && (r <= 'z')) || ((r >= '0') && (r <= '9')) || (r == '_') {
				str.WriteRune(r)
			} else {
				break
			}
		}
		if str.Len() > 16 {
			panic(fmt.Errorf(errNameTooLongMsg, str.String()))
		}
		return LexToken{Name, str.String()}
	}

	return cUndefined
}

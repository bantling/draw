package parse

// Lex the drawing language tokens
// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"strings"
)

var (
	errInvalidUnicodeEscapeMsg = "Invalid unicode escape string %s: must be \\uXXXX, \\uXXXXXX, \\U+XXXX, or \\U+XXXXXX"
	errInvalidColourMsg        = "Invalid colour %s: there must be six hex characters after the #"
	errIntTooLargeMsg          = "Invalid number %s: it is too large to fit into a 64 bit integer"
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
	cEol            = LexToken{Eol, "\n", 0}
	cPercent        = LexToken{Percent, "%", 0}
	cAssignModulus  = LexToken{AssignModulus, "%=", 0}
	cOParens        = LexToken{OParens, "(", 0}
	cCParens        = LexToken{CParens, ")", 0}
	cStar           = LexToken{Star, "*", 0}
	cAssignMultiply = LexToken{AssignMultiply, "*=", 0}
	cPlus           = LexToken{Plus, "+", 0}
	cAssignAdd      = LexToken{AssignAdd, "+=", 0}
	cIncrement      = LexToken{Increment, "++", 0}
	cComma          = LexToken{Comma, ",", 0}
	cMinus          = LexToken{Minus, "-", 0}
	cAssignSubtract = LexToken{AssignSubtract, "-=", 0}
	cDecrement      = LexToken{Decrement, "--", 0}
	cSlash          = LexToken{Slash, "/", 0}
	cAssignDivide   = LexToken{AssignDivide, "/=", 0}
	cColon          = LexToken{Colon, ":", 0}
	cLessThan       = LexToken{LessThan, "<", 0}
	cEquals         = LexToken{Equals, "=", 0}
	cGreaterThan    = LexToken{GreaterThan, ">", 0}
	cOBracket       = LexToken{OBracket, "[", 0}
	cCBracket       = LexToken{CBracket, "]", 0}
	cOBrace         = LexToken{OBrace, "{", 0}
	cCBrace         = LexToken{CBrace, "}", 0}
	cEof            = LexToken{Eof, "", 0}
	cUndefined      = LexToken{Undefined, "", 0}
)

// LexToken describes a single token, as a TokenType and a string of characters
type LexToken struct {
	TokenType
	Token string
	IntValue uint64
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
		r rune
		chars = prefix
	)

	// Has to have at least 4 hex chars
	for i := 0; i < 4; i++ {
		r = nextRune(src)
		chars += string(r)
		v, haveIt := hexVal(r)
		if !haveIt {
			panic(fmt.Errorf(errInvalidUnicodeEscapeMsg, chars))
		}

		res = res*16 + v
	}

	// May be 6 hex chars
	r = nextRune(src)
	chars += string(r)
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
	chars += string(r)
	v, haveIt = hexVal(r)
	if !haveIt {
		panic(fmt.Errorf(errInvalidUnicodeEscapeMsg, chars))
	}

	// Have 6 hex chars, return unicode char
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
		}
	}

	return r, false
}

// Helper function to read a single quoted string
// Single quoted strings end with an unescaped single quote, and can have escaped or embedded newlines
func readString(src io.RuneScanner) string {
	var str strings.Builder

	for {
		r, escaped := escapedChar(src)
		str.WriteRune(r)

		switch r {
		case '\'':
			if !escaped {
				// Complete single line string
				return str.String()
			}
		case 0:
			panic(errUnexpectedEOF)
		}
	}

	// All switch cases in above for loop return or panic, so this line can never be reached
}

func readBinaryNumber(src io.RuneScanner) LexToken {
	var (
		chars = "0b"
		n1,n2 uint64
	)
	
	for {
		r := nextRune(src)
		
		switch {
		case (r == '0') || (r == '1'):
			chars += string(r)
			n2 = n1 * 2 + uint64(r - '0')
			if n2 < n1 {
				panic(fmt.Errorf(errIntTooLargeMsg, chars))
			}
			n1 = n2
		
		case r == '_': // separator, ignore it as far as the value goes
			chars += string(r)
		
		default:
			// first char of next token
			src.UnreadRune()
			return LexToken{IntNumber, chars, n1}
		}
	}
}

func readHexNumber(src io.RuneScanner) LexToken {
	var (
		chars = "0x"
		n1,n2 uint64
	)
	
	for {
		r := nextRune(src)
		i, haveIt := hexVal(r)
		
		switch {
		case haveIt:
			chars += string(i)
			n2 = n1 * 16 + i
			if n2 < n1 {
				panic(fmt.Errorf(errIntTooLargeMsg, chars))
			}
			n1 = n2
		
		case r == '_': // separator, ignore it as far as the value goes
			chars += string(r)
		
		default:
			// first char of next token
			src.UnreadRune()
			return LexToken{IntNumber, chars, n1}
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
		chars := "#"
		for i := 0; i < 6; i++ {
			c, haveIt := hexVal(nextRune(src))
			chars += string(c)
			if !haveIt {
				panic(fmt.Errorf(errInvalidColourMsg, chars))
			}
		}
		return LexToken{Colour, chars}

	case r == '%':
		// Could be % or %=
		switch r = nextRune(src); r {
		case '=': // %=
			return cAssignModulus
		default: // /
			src.UnreadRune()
			return cPercent
		}

	case r == '\'':
		// string, read all until next unescaped ", interpreting escapes, and allowing embedded newlines
		return LexToken{Str, readString(src)}

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
				return LexToken{BinaryNumber, readBinaryNumber(src)}
			
			case r == 'x': // hex number, read all hex and _
				return LexToken{HexNumber, readHexNumber(src)}
			
			case (r >= '0') && (r <= '9'): // decimal with leading 0
		}
	}

	return cUndefined
}

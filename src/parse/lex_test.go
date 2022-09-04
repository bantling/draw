package parse

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEol(t *testing.T) {
	for _, str := range []string{"\r", "\n", "\r\n"} {
		src := strings.NewReader(str)
		assert.Equal(t, cEol, Lex(src))
		assert.Equal(t, cEof, Lex(src))
	}
}

func TestPercent(t *testing.T) {
	src := strings.NewReader("%")
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("%%")
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestAssignModulus(t *testing.T) {
	src := strings.NewReader("%=")
	assert.Equal(t, cAssignModulus, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("%=%")
	assert.Equal(t, cAssignModulus, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestOParens(t *testing.T) {
	src := strings.NewReader("(")
	assert.Equal(t, cOParens, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("(%")
	assert.Equal(t, cOParens, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestCParens(t *testing.T) {
	src := strings.NewReader(")")
	assert.Equal(t, cCParens, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader(")%")
	assert.Equal(t, cCParens, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestStar(t *testing.T) {
	src := strings.NewReader("*")
	assert.Equal(t, cStar, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("*%")
	assert.Equal(t, cStar, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestAssignMultiply(t *testing.T) {
	src := strings.NewReader("*=")
	assert.Equal(t, cAssignMultiply, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("*=%")
	assert.Equal(t, cAssignMultiply, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestPlus(t *testing.T) {
	src := strings.NewReader("+")
	assert.Equal(t, cPlus, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("+%")
	assert.Equal(t, cPlus, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestAssignAdd(t *testing.T) {
	src := strings.NewReader("+=")
	assert.Equal(t, cAssignAdd, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("+=%")
	assert.Equal(t, cAssignAdd, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestIncrement(t *testing.T) {
	src := strings.NewReader("++")
	assert.Equal(t, cIncrement, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("++%")
	assert.Equal(t, cIncrement, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestComma(t *testing.T) {
	src := strings.NewReader(",")
	assert.Equal(t, cComma, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader(",%")
	assert.Equal(t, cComma, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestMinus(t *testing.T) {
	src := strings.NewReader("-")
	assert.Equal(t, cMinus, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("-%")
	assert.Equal(t, cMinus, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestAssignSubtract(t *testing.T) {
	src := strings.NewReader("-=")
	assert.Equal(t, cAssignSubtract, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("-=%")
	assert.Equal(t, cAssignSubtract, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestDecrement(t *testing.T) {
	src := strings.NewReader("--")
	assert.Equal(t, cDecrement, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("--%")
	assert.Equal(t, cDecrement, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestSlash(t *testing.T) {
	src := strings.NewReader("/")
	assert.Equal(t, cSlash, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("/%")
	assert.Equal(t, cSlash, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestAssignDivide(t *testing.T) {
	src := strings.NewReader("/=")
	assert.Equal(t, cAssignDivide, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("/=%")
	assert.Equal(t, cAssignDivide, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestColon(t *testing.T) {
	src := strings.NewReader(":")
	assert.Equal(t, cColon, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader(":%")
	assert.Equal(t, cColon, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestLessThan(t *testing.T) {
	src := strings.NewReader("<")
	assert.Equal(t, cLessThan, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("<%")
	assert.Equal(t, cLessThan, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestEquals(t *testing.T) {
	src := strings.NewReader("=")
	assert.Equal(t, cEquals, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("=%")
	assert.Equal(t, cEquals, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestGreaterThan(t *testing.T) {
	src := strings.NewReader(">")
	assert.Equal(t, cGreaterThan, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader(">%")
	assert.Equal(t, cGreaterThan, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestOBracket(t *testing.T) {
	src := strings.NewReader("[")
	assert.Equal(t, cOBracket, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("[%")
	assert.Equal(t, cOBracket, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestCBracket(t *testing.T) {
	src := strings.NewReader("]")
	assert.Equal(t, cCBracket, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("]%")
	assert.Equal(t, cCBracket, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestOBrace(t *testing.T) {
	src := strings.NewReader("{")
	assert.Equal(t, cOBrace, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("{%")
	assert.Equal(t, cOBrace, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestCBrace(t *testing.T) {
	src := strings.NewReader("}")
	assert.Equal(t, cCBrace, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("}%")
	assert.Equal(t, cCBrace, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestEof(t *testing.T) {
	src := strings.NewReader("")
	assert.Equal(t, cEof, Lex(src))
}

func TestUndefined(t *testing.T) {
	src := strings.NewReader("~")
	assert.Equal(t, cUndefined, Lex(src))
	assert.Equal(t, cEof, Lex(src))

	src = strings.NewReader("~%")
	assert.Equal(t, cUndefined, Lex(src))
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
}

func TestColour(t *testing.T) {
	src := strings.NewReader("#123456")
	tok := Lex(src)
	assert.Equal(t, LexToken{Colour, "#123456"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, uint64(0x123456), tok.IntValue())

	src = strings.NewReader("#123456%")
	tok = Lex(src)
	assert.Equal(t, LexToken{Colour, "#123456"}, tok)
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, uint64(0x123456), tok.IntValue())
}

func TestFloatNumber(t *testing.T) {
	src := strings.NewReader("12.34")
	tok := Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34), tok.FloatValue())

	src = strings.NewReader("12.34%")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34"}, tok)
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34), tok.FloatValue())

	src = strings.NewReader("12.34.")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34"}, tok)
	assert.Equal(t, cUndefined, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34), tok.FloatValue())

	src = strings.NewReader("12e26")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12e26"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12e26), tok.FloatValue())

	src = strings.NewReader("12E26")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12E26"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12e26), tok.FloatValue())

	src = strings.NewReader("12e26%")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12e26"}, tok)
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12e26), tok.FloatValue())

	src = strings.NewReader("12E26.")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12E26"}, tok)
	assert.Equal(t, cUndefined, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12e26), tok.FloatValue())

	src = strings.NewReader("12.34e26")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34e26"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34e26), tok.FloatValue())

	src = strings.NewReader("12.34E26")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34E26"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34e26), tok.FloatValue())

	src = strings.NewReader("12.34e26%")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34e26"}, tok)
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34e26), tok.FloatValue())

	src = strings.NewReader("12.34E26.")
	tok = Lex(src)
	assert.Equal(t, LexToken{FloatNumber, "12.34E26"}, tok)
	assert.Equal(t, cUndefined, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, float32(12.34e26), tok.FloatValue())

	func() {
		defer func() {
			assert.Equal(t, fmt.Errorf(errIncompleteFloatMsg, "12."), recover())
		}()

		src = strings.NewReader("12.")
		Lex(src)
		assert.Fail(t, "Must die")
	}()

	func() {
		defer func() {
			assert.Equal(t, fmt.Errorf(errIncompleteFloatMsg, "12e"), recover())
		}()

		src = strings.NewReader("12e")
		Lex(src)
		assert.Fail(t, "Must die")
	}()

	func() {
		defer func() {
			assert.Equal(t, fmt.Errorf(errIncompleteFloatMsg, "12E"), recover())
		}()

		src = strings.NewReader("12E")
		Lex(src)
		assert.Fail(t, "Must die")
	}()

	func() {
		defer func() {
			assert.Equal(t, strconv.ErrRange, recover())
		}()

		src = strings.NewReader("12e500")
		Lex(src).FloatValue()
		assert.Fail(t, "Must die")
	}()
}

func TestIntNumber(t *testing.T) {
	src := strings.NewReader("12")
	tok := Lex(src)
	assert.Equal(t, LexToken{IntNumber, "12"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, uint64(12), tok.IntValue())

	src = strings.NewReader("12%")
	tok = Lex(src)
	assert.Equal(t, LexToken{IntNumber, "12"}, tok)
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, uint64(12), tok.IntValue())

	src = strings.NewReader("18446744073709551615")
	tok = Lex(src)
	assert.Equal(t, LexToken{IntNumber, "18446744073709551615"}, tok)
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, uint64(math.MaxUint64), tok.IntValue())

	src = strings.NewReader("18446744073709551615%")
	tok = Lex(src)
	assert.Equal(t, LexToken{IntNumber, "18446744073709551615"}, tok)
	assert.Equal(t, cPercent, Lex(src))
	assert.Equal(t, cEof, Lex(src))
	assert.Equal(t, uint64(math.MaxUint64), tok.IntValue())

	func() {
		defer func() {
			assert.Equal(t, strconv.ErrRange, recover())
		}()

		src = strings.NewReader("18446744073709551616")
		tok = Lex(src)
		assert.Equal(t, LexToken{IntNumber, "18446744073709551616"}, tok)
		assert.Equal(t, cEof, Lex(src))
		tok.IntValue()
		assert.Fail(t, "Must die")
	}()

	func() {
		defer func() {
			assert.Equal(t, strconv.ErrRange, recover())
		}()

		src = strings.NewReader("18446744073709551616%")
		tok = Lex(src)
		assert.Equal(t, LexToken{IntNumber, "18446744073709551616"}, tok)
		assert.Equal(t, cPercent, Lex(src))
		assert.Equal(t, cEof, Lex(src))
		tok.IntValue()
		assert.Fail(t, "Must die")
	}()
}

func TestName(t *testing.T) {
	src := strings.NewReader("A1_")
	assert.Equal(t, LexToken{Name, "A1_"}, Lex(src))
	src = strings.NewReader("a1_")
	assert.Equal(t, LexToken{Name, "a1_"}, Lex(src))

	func() {
		str := "abcdef1234567890_"
		defer func() {
			assert.Equal(t, fmt.Errorf(errNameTooLongMsg, str), recover())
		}()

		Lex(strings.NewReader(str))
		assert.Fail(t, "Must die")
	}()
}

func TestStr(t *testing.T) {
	src := strings.NewReader("'an example STRING \\\\ \\' \\n \\u0041 \\u010000 \\U+0061 \\U+010000'")
	assert.Equal(t, LexToken{Str, "'an example STRING \\ ' \n A \U00010000 a \U00010000'"}, Lex(src))

	func() {
		defer func() {
			assert.Equal(t, fmt.Errorf(errInvalidEscapeMsg, "\\z"), recover())
		}()

		Lex(strings.NewReader("'\\z'"))
		assert.Fail(t, "Must die")
	}()
}

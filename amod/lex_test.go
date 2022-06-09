package amod

import (
	"strings"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
)

func TestNumber(t *testing.T) {
	t.Parallel()

	// https://github.com/asmaloney/gactar/issues/2
	src := `0 0.3 5 55.6 .9`

	l := lex("test", src)

	expecteds := strings.Split(src, " ")

	for i, expected := range expecteds {
		token, err := l.Next()
		if err != nil {
			t.Errorf("[index %d] error getting next token: %s", i, err.Error())
		}

		if token.Type != lexer.TokenType(lexemeNumber) {
			t.Errorf("[index %d] expected to lex '%s' as int (%d) - got type %d", i, token.Value, lexemeNumber, token.Type)
		}
		if token.Value != expected {
			t.Errorf("[index %d] expected token value: %s - got %s", i, expected, token.Value)
		}
	}
}

func TestInvalidSection(t *testing.T) {
	t.Parallel()

	l := lex("test", "==invalid==")

	token, err := l.Next()
	if err != nil {
		t.Errorf("error getting next token: %s", err.Error())
	}

	if token.Type != lexer.TokenType(lexemeChar) {
		t.Errorf("expected to lex '%s' as int (%d) - got type %d", token.Value, lexemeChar, token.Type)
	}
}

func TestUnterminatedQuote(t *testing.T) {
	t.Parallel()

	l := lex("test", `"a string`)

	_, err := l.Next()

	expected := "ERROR on line 1 at position 8: unterminated quoted string"

	if err == nil {
		t.Errorf("expected error: %q", expected)
	} else if err.Error() != expected {
		t.Errorf("expected error: %q but got %q", expected, err.Error())
	}
}

func TestForwardSlashNotComment(t *testing.T) {
	t.Parallel()

	l := lex("test", "/foo")

	_, err := l.Next()

	if err != nil {
		t.Errorf("unexpected error: %q", err.Error())
	}
}

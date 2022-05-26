package amod

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/asmaloney/gactar/issues"
)

// Log wraps issues.Log so we can provide extra convenience functions.
type Log struct {
	issues.Log
}

// ErrorT constructs our location information from tokens and uses that to add an error.
func (l *Log) ErrorT(tokens []lexer.Token, s string, a ...interface{}) {
	l.Log.Error(tokensToLocation(tokens), s, a...)
}

// ErrorT constructs our location information from a range of tokens and uses that to add an error.
func (l *Log) ErrorTR(tokens []lexer.Token, start, end int, s string, a ...interface{}) {
	l.Log.Error(tokenRangeToLocation(tokens, start, end), s, a...)
}

// tokensToLocation takes the list of lexer tokens and converts it to our own
// issues.Location struct.
func tokensToLocation(t []lexer.Token) *issues.Location {
	tokens := trimCommentsFromRange(t)

	// If we are in the middle of a pattern, we might have pattern spaces,
	// so remove from beginning & end.
	if tokens[0].Type == lexer.TokenType(lexemePatternSpace) {
		tokens = tokens[1:]
	}
	if tokens[len(tokens)-1].Type == lexer.TokenType(lexemePatternSpace) {
		tokens = tokens[:len(tokens)-1]
	}

	// first & last may end being the same - that's ok
	firstToken := tokens[0]
	lastToken := tokens[len(tokens)-1]

	// Because the parser strips quotes (see var amodParser), we need to
	// account for them here.
	lastTokenLen := len(lastToken.Value)
	if lastToken.Type == lexer.TokenType(lexemeString) {
		lastTokenLen += 2
	}

	return &issues.Location{
		Line:        firstToken.Pos.Line,
		ColumnStart: firstToken.Pos.Column,
		ColumnEnd:   lastToken.Pos.Column + lastTokenLen,
	}
}

func tokenRangeToLocation(t []lexer.Token, start, end int) *issues.Location {
	if start < 0 || end < 1 || start == end || end < start {
		fmt.Printf("Internal error (tokenRangeToLocation): start (%d) and/or end (%d) incorrect. Using full range.\n", start, end)
		return tokensToLocation(t)
	}

	tokens := trimCommentsFromRange(t)

	numTokens := len(tokens)
	if end > numTokens-1 {
		fmt.Printf("Internal error (tokenRangeToLocation): end (%d - 0-indexed) greater than tokens len (%d). Using full range.\n", end, numTokens)
		return tokensToLocation(tokens)
	}

	restricted := tokens[start:end]

	return tokensToLocation(restricted)
}

// trimCommentsFromRange will remove any comment tokens from the beginning and end of the range.
// This is necessary because participle will include them with the Tokens in a struct.
func trimCommentsFromRange(t []lexer.Token) (tokens []lexer.Token) {
	begin := 0
	for _, token := range t {
		if token.Type == lexer.TokenType(lexemeComment) {
			begin++
			continue
		}

		break
	}

	end := len(t)
	for i := end - 1; i >= 0; i-- {
		token := t[i]

		if token.Type == lexer.TokenType(lexemeComment) {
			end--
			continue
		}

		break
	}

	tokens = t[begin:end]
	return
}

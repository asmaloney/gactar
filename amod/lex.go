package amod

// Mostly based on Rob Pike's talk:
// 	https://www.youtube.com/watch?v=HxaD_trXwRE
// Not sure I implemented precisely what he's advocating since I ended
// up with a central switch anyways. I don't see how it can be avoided.

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/alecthomas/participle/v2/lexer"
)

type lexer_def struct {
	lexer.Definition
}

// LexerDefinition provides the interface for the participle parser
var LexerDefinition lexer.Definition = lexer_def{}

// LexError is returned for errors we have while lexing.
type LexError struct {
	Line     int
	Position int
	Value    string
}

func (err LexError) Error() string {
	return fmt.Sprintf("ERROR on line %d at position %d: %s", err.Line, err.Position, err.Value)
}

type lexemeType int

const (
	lexemeError lexemeType = iota

	lexemeSpace
	LexemeEOF

	lexemeComment
	lexemeIdentifier
	lexemeKeyword
	lexemeNumber
	lexemeString
	lexemeChar

	lexemeEquality
	lexemeInequality

	lexemeSectionDelim

	lexemePatternVar
	lexemePatternWildcard
)

func (l lexemeType) String() string {
	switch l {
	case lexemeError:
		return "error"
	case lexemeSpace:
		return "space"
	case LexemeEOF:
		return "EOF"
	case lexemeComment:
		return "comment"
	case lexemeIdentifier:
		return "identifier"
	case lexemeKeyword:
		return "keyword"
	case lexemeNumber:
		return "number"
	case lexemeString:
		return "string"
	case lexemeChar:
		return "char"
	case lexemeEquality:
		return "equality"
	case lexemeInequality:
		return "inequality"
	case lexemeSectionDelim:
		return "section delimiter"
	case lexemePatternVar:
		return "pattern var"
	case lexemePatternWildcard:
		return "pattern wildcard"
	}

	return "unknown"
}

type lexeme struct {
	typ   lexemeType
	value string
	line  int // line number this lexeme is on
	pos   int // position within the line
}

// sectionType is used to keep track of what section we are lexing
// We use this to limit the scope of keywords.
type sectionType int

const (
	sectionModel sectionType = iota
	sectionConfig
	sectionInit
	sectionProduction
)

// lexer_amod tracks our lexing and provides a channel to emit lexemes
type lexer_amod struct {
	name           string // used only for error reports
	input          string // the string being scanned.
	line           int    // the line number
	lastNewlinePos int
	start          int         // start position of this lexeme (offset from beginning of file)
	pos            int         // current position in the input  (offset from beginning of file)
	width          int         // width of last rune read from input
	lexemes        chan lexeme // channel of scanned lexemes

	inSectionHeader bool        // state: switch currentSection based on ~~ section headers
	currentSection  sectionType // which section are we lexing? used to switch out keywords
	inPattern       bool        // state: a pattern - delimited by [] is lexed specially
}

// stateFn is used to move through the lexing states
type stateFn func(*lexer_amod) stateFn

const (
	eof = -1

	commentDelim = "//"
)

// keywordsModel are only keywords for the model section
var keywordsModel []string = []string{
	"authors",
	"description",
	"examples",
	"name",
	"nil",
}

// keywordsModel are only keywords for the config section
var keywordsConfig []string = []string{
	"chunks",
	"gactar",
	"modules",
	"nil",
}

// keywordsModel are only keywords for the init section
var keywordsInit []string = []string{
	"nil",
	"similar",
}

// keywordsModel are only keywords for the productions section
var keywordsProductions []string = []string{
	"and",
	"any",
	"buffer_state",
	"clear",
	"description",
	"do",
	"match",
	"module_state",
	"nil",
	"print",
	"recall",
	"set",
	"stop",
	"to",
	"when",
	"with",
}

// Symbols provides a mapping from participle strings to our lexemes
func (lexer_def) Symbols() map[string]lexer.TokenType {
	return map[string]lexer.TokenType{
		"Comment":      lexer.TokenType(lexemeComment),
		"Whitespace":   lexer.TokenType(lexemeSpace),
		"Keyword":      lexer.TokenType(lexemeKeyword),
		"Ident":        lexer.TokenType(lexemeIdentifier),
		"Number":       lexer.TokenType(lexemeNumber),
		"String":       lexer.TokenType(lexemeString),
		"Char":         lexer.TokenType(lexemeChar),
		"Equality":     lexer.TokenType(lexemeEquality),
		"Inequality":   lexer.TokenType(lexemeInequality),
		"SectionDelim": lexer.TokenType(lexemeSectionDelim),
		"Var":          lexer.TokenType(lexemePatternVar),
		"Wildcard":     lexer.TokenType(lexemePatternWildcard),
	}
}

// Lex is called by the participle parser to lex a reader
func (lexer_def) Lex(filename string, r io.Reader) (lexer.Lexer, error) {
	s := &strings.Builder{}
	_, err := io.Copy(s, r)
	if err != nil {
		return nil, err
	}

	data := s.String()

	l := lex(filename, data)

	return l, nil
}

func lex(filename string, data string) *lexer_amod {
	cleanData(&data)

	l := &lexer_amod{
		name:            filename,
		input:           data,
		line:            1,
		lastNewlinePos:  1, // start @ 1 so first line gets 0 (see emit())
		lexemes:         make(chan lexeme),
		currentSection:  sectionModel,
		inSectionHeader: false,
		inPattern:       false,
	}

	go l.run()

	return l
}

// Next is used by participle to get the next token
func (l *lexer_amod) Next() (tok lexer.Token, err error) {
	next := <-l.lexemes

	pos := lexer.Position{
		Filename: l.name,
		Offset:   l.pos,
		Line:     next.line,
		Column:   next.pos,
	}

	if next.typ == LexemeEOF {
		return lexer.EOFToken(pos), nil
	}

	tok = lexer.Token{
		Type:  lexer.TokenType(next.typ),
		Value: next.value,
		Pos:   pos,
	}

	if next.typ == lexemeError {
		err = LexError{
			Line:     next.line,
			Position: next.pos,
			Value:    next.value,
		}
		return
	}

	if debugLex {
		fmt.Printf("TOK (%d, %d-%d):\t%+v (%s)\n", pos.Line, pos.Column, pos.Column+len(tok.Value), tok, next.typ.String())
	}
	return
}

func (l *lexer_amod) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width

	return r
}

// lookupKeyword looks up "id" to see if it is a keyword based on which section we are lexing
func (l *lexer_amod) lookupKeyword(id string) bool {
	switch l.currentSection {
	case sectionModel:
		return slices.Contains(keywordsModel, id)
	case sectionConfig:
		return slices.Contains(keywordsConfig, id)
	case sectionInit:
		return slices.Contains(keywordsInit, id)
	case sectionProduction:
		return slices.Contains(keywordsProductions, id)
	}

	return false
}

// skip over the pending input before this point
func (l *lexer_amod) ignore() {
	l.start = l.pos
}

// step back one rune
func (l *lexer_amod) backup() {
	l.pos -= l.width
}

// look at the next rune in the input, but don't eat it
func (l *lexer_amod) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// check if next rune is "r"
func (l *lexer_amod) nextIs(r rune) bool {
	return l.peek() == r
}

// accept any character in the string
func (l *lexer_amod) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}

	l.backup()

	return false
}

// accept a run of any characters in the string
func (l *lexer_amod) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}

	l.backup()
}

// pass an item back to the client via the channel
func (l *lexer_amod) emit(t lexemeType) {
	value := l.input[l.start:l.pos]
	l.lexemes <- lexeme{
		typ:   t,
		value: value,
		line:  l.line,
		pos:   l.start - l.lastNewlinePos + 1,
	}

	l.start = l.pos
}

// declare and error and let the client know where we are in the input
func (l *lexer_amod) errorf(format string, args ...interface{}) stateFn {
	l.lexemes <- lexeme{
		lexemeError,
		fmt.Sprintf(format, args...),
		l.line,
		l.pos - l.lastNewlinePos,
	}

	return nil
}

func (l *lexer_amod) run() {
	for state := lexStart; state != nil; {
		// name := runtime.FuncForPC(reflect.ValueOf(state).Pointer()).Name()
		// fmt.Printf("%s\n", name)

		state = state(l)
	}

	close(l.lexemes)
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

// newlines have been normalized, so just check the one
func isNewline(r rune) bool {
	return r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isDigit(r rune) bool {
	return ('0' <= r && r <= '9')
}

func lexStart(l *lexer_amod) stateFn {
	switch r := l.next(); {
	case isSpace(r):
		if isNewline(r) {
			l.lastNewlinePos = l.pos + 1
			l.line++
		}
		return lexSpace

	case isDigit(r) || r == '.':
		return lexNumber

	case (r == '+') || (r == '-'):
		p := l.peek()
		if isDigit(p) || p == '.' {
			// eat the +/- char and let the previous case handle a number
			break
		}
		l.emit(lexemeChar)

	case isAlphaNumeric(r):
		return lexIdentifier

	case r == '=':
		if l.nextIs('=') {
			l.next()
			l.emit(lexemeEquality)
		} else {
			l.emit(lexemeChar)
		}

	case r == '!':
		if l.nextIs('=') {
			l.next()
			l.emit(lexemeInequality)
		} else {
			l.emit(lexemeChar)
		}

	case r == '/':
		if l.nextIs('/') {
			l.backup()
			return lexComment
		}
		l.emit(lexemeChar)

	case r == '"' || r == '\'':
		l.backup()
		return lexQuotedString

	case r == '[':
		l.inPattern = true
		l.emit(lexemeChar)

	case r == ']':
		l.emit(lexemeChar)
		l.inPattern = false

	case r == '?':
		if isAlphaNumeric(l.peek()) {
			return lexIdentifier
		}

		l.emit(lexemeChar)

	case r == '*':
		if l.inPattern {
			l.emit(lexemePatternWildcard)
			break
		}

		l.emit(lexemeChar)

	case r == '~':
		if l.nextIs('~') {
			l.next()
			l.emit(lexemeSectionDelim)
			l.inSectionHeader = !l.inSectionHeader
		} else {
			l.emit(lexemeChar)
		}

	case r <= unicode.MaxASCII && unicode.IsPrint(r):
		l.emit(lexemeChar)

	case r == eof:
		l.emit(LexemeEOF)
		return nil
	}

	return lexStart
}

// consume 0 or more spaces
func eatSpace(l *lexer_amod) {
	for {
		r := l.next()

		if !isSpace(r) {
			l.backup()
			break
		}

		if isNewline(r) {
			l.lastNewlinePos = l.pos + 1
			l.line++
		}
	}
	l.ignore()
}

func lexSpace(l *lexer_amod) stateFn {
	eatSpace(l)
	return lexStart
}

func lexComment(l *lexer_amod) stateFn {
	l.pos += len(commentDelim)
	i := strings.Index(l.input[l.pos:], "\n")

	// If we are at the end of file there may not be a newline,
	// so take the rest of the input.
	if i == -1 {
		l.pos = len(l.input)
	} else {
		l.pos += i
	}

	l.emit(lexemeComment)

	eatSpace(l)
	return lexStart
}

func lexIdentifier(l *lexer_amod) stateFn {
	for {
		r := l.peek()

		if !isAlphaNumeric(r) {
			break
		}

		l.next()
	}

	id := l.input[l.start:l.pos]
	isKeyword := false

	// If we are in a section header, then change our current section
	if l.inSectionHeader {
		switch id {
		case "model":
			l.currentSection = sectionModel
		case "config":
			l.currentSection = sectionConfig
		case "init":
			l.currentSection = sectionInit
		case "productions":
			l.currentSection = sectionProduction
		default:
			return l.errorf("unrecognized section")
		}

		// these are keywords in this context
		isKeyword = true
	} else {
		isKeyword = l.lookupKeyword(id)
	}

	// Perhaps not the best way to do this.
	// I'm sure there's a char-by-char way we could implement which would be faster.
	switch {
	case isKeyword:
		l.emit(lexemeKeyword)

	case l.input[l.start] == '?':
		l.emit(lexemePatternVar)

	default:
		l.emit(lexemeIdentifier)
	}

	return lexStart
}

func lexNumber(l *lexer_amod) stateFn {
	current := l.input[l.pos-1]
	digits := "0123456789"

	// used to determine if we added anything to our buffer
	savePos := l.pos

	l.acceptRun(digits)

	if current != '.' && l.accept(".") {
		l.acceptRun(digits)
	}

	// If we only found '.' without any numbers before or after, return it as a char
	if (current == '.') && (l.pos == savePos) {
		l.emit(lexemeChar)
	} else {
		l.emit(lexemeNumber)
	}

	return lexStart
}

func lexQuotedString(l *lexer_amod) stateFn {
	quoteType := l.next()
	done := false

	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof:
			fallthrough
		case '\n':
			return l.errorf("unterminated quoted string")
		case quoteType:
			done = true
		}

		if done {
			break
		}
	}

	l.emit(lexemeString)

	return lexSpace
}

// cleanData normalizes line endings
func cleanData(data *string) {
	*data = strings.ReplaceAll(*data, "\r\n", "\n")
	*data = strings.ReplaceAll(*data, "\r", "\n")
}

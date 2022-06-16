package amod

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Uses participle to parse the lexemes.
// 	https://github.com/alecthomas/participle

// Railroad Diagrams
// ------
// First output the EBNF grammar to stdout with the command "gactar -ebnf".
//
// There are two ways to generate railroad diagrams:
// 	1. Use the "railroad" tool from participle like this:
//		./railroad -o amod-grammar.html -w
//		paste in the generated EBNF above & hit control-D
//	2. Use this page to convert the ebnf and generate a diagram:
//		https://bottlecaps.de/convert/
//		paste in the generated EBNF above, click "Convert" and then click "View Diagram"

type amodFile struct {
	ModelHeader string        `parser:"@SectionDelim 'model' @SectionDelim"`
	Model       *modelSection `parser:"@@"`

	ConfigHeader string         `parser:"@SectionDelim 'config' @SectionDelim"`
	Config       *configSection `parser:"(@@)?"`

	InitHeader string       `parser:"@SectionDelim 'init' @SectionDelim"`
	Init       *initSection `parser:"(@@)?"`

	ProductionsHeader string             `parser:"@SectionDelim 'productions' @SectionDelim"`
	Productions       *productionSection `parser:"(@@)?"`

	Tokens []lexer.Token
}

type modelSection struct {
	Name        string     `parser:"'name' ':' (@String|@Ident)"`
	Description string     `parser:"('description' ':' @String)?"`
	Authors     []string   `parser:"('authors' '{' @String* '}')?"`
	Examples    []*pattern `parser:"('examples' '{' @@* '}')?"`

	Tokens []lexer.Token
}

type arg struct {
	Nil    *bool   `parser:"( @('nil':Keyword)"`
	Var    *string `parser:"| @PatternVar"`
	ID     *string `parser:"| @Ident"`
	Str    *string `parser:"| @String"`
	Number *string `parser:"| @Number)"`

	Tokens []lexer.Token
}

type fieldValue struct {
	ID     *string  `parser:"( @Ident"`
	Str    *string  `parser:"| @String"`
	Number *float64 `parser:"| @Number)"`

	Tokens []lexer.Token
}

// Used for outputting errors
func (f fieldValue) String() string {
	switch {
	case f.ID != nil:
		return *f.ID

	case f.Str != nil:
		return *f.Str

	case f.Number != nil:
		return fmt.Sprintf("%f", *f.Number)
	}

	return "<error>"
}

type field struct {
	Key   string     `parser:"@Ident ':'"`
	Value fieldValue `parser:"@@"`

	Tokens []lexer.Token
}

type chunkDecl struct {
	StartBracket string   `parser:"'['"` // not used - must be set for parse
	Name         string   `parser:"@Ident ':'"`
	Slots        []string `parser:"@Ident+"`
	EndBracket   string   `parser:"']'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type module struct {
	Name       string   `parser:"@Ident"`
	InitFields []*field `parser:"'{' @@ (','? @@)* '}'"`

	Tokens []lexer.Token
}

type configSection struct {
	GACTAR     []*field     `parser:"('gactar' '{' @@ (','? @@)* '}')?"`
	Modules    []*module    `parser:"('modules' '{' @@* '}')?"`
	ChunkDecls []*chunkDecl `parser:"('chunks' '{' @@* '}')?"`

	Tokens []lexer.Token
}

type initialization struct {
	Name         string     `parser:"@Ident"`
	InitPatterns []*pattern `parser:"( '{' @@+ '}' | @@ )"`

	Tokens []lexer.Token
}

type initSection struct {
	Initializations []*initialization `parser:"@@*"`

	Tokens []lexer.Token
}

type patternSlot struct {
	Not      bool    `parser:"((@('!':Char)?"`
	Nil      *bool   `parser:"( @('nil':Keyword)"`
	ID       *string `parser:"| @Ident"`
	Num      *string `parser:"| @Number"` // we don't need to treat this as a number anywhere, so keep as a string
	Var      *string `parser:"| @PatternVar ))"`
	Wildcard *string `parser:"| @PatternWildcard)"`

	Tokens []lexer.Token
}

type pattern struct {
	StartBracket string         `parser:"'['"` // not used - must be set for parse
	ChunkName    string         `parser:"@Ident ':'"`
	Slots        []*patternSlot `parser:"@@+"`
	EndBracket   string         `parser:"']'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type comparisonOperator struct {
	Equal    *string `parser:"( @Equality"`
	NotEqual *string `parser:"| @Inequality )"`

	Tokens []lexer.Token
}

type whenExpression struct {
	OpenParen  string              `parser:"'('"` // not used - must be set for parse
	LHS        string              `parser:"@PatternVar"`
	Comparison *comparisonOperator `parser:"@@"`
	RHS        *arg                `parser:"@@"`
	CloseParen string              `parser:"')'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type whenClause struct {
	When        string             `parser:"'when':Keyword"`
	Expressions *[]*whenExpression `parser:"@@ ('and' @@)*"`

	Tokens []lexer.Token
}

type matchItem struct {
	Name    string      `parser:"@Ident"`
	Pattern *pattern    `parser:"@@"`
	When    *whenClause `parser:"@@?"`

	Tokens []lexer.Token
}

type match struct {
	Items []*matchItem `parser:"'match' '{' @@+ '}'"`

	Tokens []lexer.Token
}

type clearStatement struct {
	BufferNames []string `parser:"'clear' ( @Ident ','? )+"`

	Tokens []lexer.Token
}

type printStatement struct {
	Args []*arg `parser:"'print' ( @@ ','? )*"`

	Tokens []lexer.Token
}

type recallStatement struct {
	Pattern *pattern `parser:"'recall' @@"`

	Tokens []lexer.Token
}

type setStatement struct {
	Set        string  `parser:"'set'"` // not used, but must be visible for parse to work
	BufferName string  `parser:"@Ident"`
	Slot       *string `parser:"('.' @Ident)?"`

	To      string   `parser:"'to'"` // not used, but must be visible for parse to work
	Value   *arg     `parser:"( @@"`
	Pattern *pattern `parser:"| @@)"`

	Tokens []lexer.Token
}

type stopStatement struct {
	Stop string `parser:"'stop':Keyword"`

	Tokens []lexer.Token
}

type statement struct {
	Clear  *clearStatement  `parser:"  @@"`
	Print  *printStatement  `parser:"| @@"`
	Recall *recallStatement `parser:"| @@"`
	Set    *setStatement    `parser:"| @@"`
	Stop   *stopStatement   `parser:"| @@"`

	Tokens []lexer.Token
}

type do struct {
	Do         string        `parser:"'do'"` // not used, but must be visible for parse to work
	Statements *[]*statement `parser:"'{' @@+ '}'"`

	Tokens []lexer.Token
}

type production struct {
	Name        string  `parser:"@Ident '{'"`
	Description *string `parser:"('description' ':' @String)?"`
	Match       *match  `parser:"@@"`
	Do          *do     `parser:"@@"`
	End         string  `parser:"'}'"` // not used, but must be visible for parse to work

	Tokens []lexer.Token
}

type productionSection struct {
	Productions []*production `parser:"@@+"`

	Tokens []lexer.Token
}

var amodParser = participle.MustBuild(&amodFile{},
	participle.Lexer(LexerDefinition),
	participle.Elide("Comment", "Whitespace"),
	participle.Unquote(),
)

var patternParser = participle.MustBuild(&pattern{},
	participle.Lexer(LexerDefinition),
	participle.Elide("Comment", "Whitespace"),
	participle.Unquote(),
)

func parse(r io.Reader) (*amodFile, error) {
	var amod amodFile

	err := amodParser.Parse("", r, &amod)

	if err != nil {
		return nil, err
	}

	return &amod, nil
}

func parseFile(filename string) (*amodFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parse(file)
}

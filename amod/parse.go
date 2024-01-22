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
	ModelHeader string        `parser:"'~~':SectionDelim 'model':Keyword '~~':SectionDelim"`
	Model       *modelSection `parser:"@@"`

	ConfigHeader string         `parser:"'~~':SectionDelim 'config':Keyword '~~':SectionDelim"`
	Config       *configSection `parser:"(@@)?"`

	InitHeader string       `parser:"'~~':SectionDelim 'init':Keyword '~~':SectionDelim"`
	Init       *initSection `parser:"(@@)?"`

	ProductionsHeader string             `parser:"'~~':SectionDelim 'productions':Keyword '~~':SectionDelim"`
	Productions       *productionSection `parser:"(@@)?"`

	Tokens []lexer.Token
}

type modelSection struct {
	Name        string     `parser:"'name':Keyword ':' (@String|@Ident)"`
	Description string     `parser:"('description':Keyword ':' @String)?"`
	Authors     []string   `parser:"('authors':Keyword '{' @String* '}')?"`
	Examples    []*pattern `parser:"('examples':Keyword '{' @@* '}')?"`

	Tokens []lexer.Token
}

type arg struct {
	Var    *string `parser:"( @Var"`
	Str    *string `parser:"| @String"`
	Number *string `parser:"| @Number )"`

	Tokens []lexer.Token
}

func (a arg) hasVar() bool {
	return a.Var != nil
}

type whenArg struct {
	Arg *arg  `parser:"( @@"`
	Nil *bool `parser:"| @('nil':Keyword) )"`

	Tokens []lexer.Token
}

func (w whenArg) hasVar() bool {
	return w.Arg != nil && w.Arg.hasVar()
}

type withArg struct {
	Arg *arg    `parser:"( @@"`
	Nil *bool   `parser:"| @('nil':Keyword)"`
	ID  *string `parser:"| @Ident )"`

	Tokens []lexer.Token
}

func (w withArg) hasVar() bool {
	return w.Arg != nil && w.Arg.hasVar()
}

type printArg struct {
	Arg       *arg       `parser:"( @@"`
	BufferRef *bufferRef `parser:"| @@)"`

	Tokens []lexer.Token
}

func (p printArg) hasVar() bool {
	return p.Arg != nil && p.Arg.hasVar()
}

type setArg struct {
	Arg *arg    `parser:"( @@"`
	Nil *bool   `parser:"| @('nil':Keyword)"`
	ID  *string `parser:"| @Ident )"`

	Tokens []lexer.Token
}

func (s setArg) hasVar() bool {
	return s.Arg != nil && s.Arg.hasVar()
}

type fieldValue struct {
	Colon  string   `parser:"(':'"`
	ID     *string  `parser:"(  @Ident"`
	Str    *string  `parser:"| @String"`
	Number *float64 `parser:"| @Number )"`

	// We can't capture an empty field "{}", so capture the open brace to tell us it's a nested field
	OpenBrace  *string  `parser:"| @('{')"`
	Fields     []*field `parser:"@@*"`
	CloseBrace *string  `parser:"'}')"`

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

	case f.OpenBrace != nil:
		return "<nested field>"
	}

	return "<error>"
}

type bufferRef struct {
	BufferName string  `parser:"@Ident"`
	SlotName   *string `parser:"('.' @Ident)?"`

	Tokens []lexer.Token
}

func (br bufferRef) String() string {
	ref := br.BufferName

	if br.SlotName != nil {
		ref += fmt.Sprintf(".%s", *br.SlotName)
	}

	return ref
}

type field struct {
	Key   string     `parser:"@Ident"`
	Value fieldValue `parser:"@@"`

	Tokens []lexer.Token
}

type gactarConfig struct {
	GactarFields []*field `parser:"'gactar':Keyword '{' @@* '}'"`

	Tokens []lexer.Token
}

type module struct {
	ModuleName string   `parser:"@Ident"`
	Fields     []*field `parser:"'{' @@* '}'"`

	Tokens []lexer.Token
}

type moduleConfig struct {
	Modules []*module `parser:"'modules':Keyword '{' @@* '}'"`

	Tokens []lexer.Token
}

type chunkDecl struct {
	StartBracket string   `parser:"'['"` // not used - must be set for parse
	TypeName     string   `parser:"@Ident ':'"`
	Slots        []string `parser:"@Ident+"`
	EndBracket   string   `parser:"']'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type chunkConfig struct {
	ChunkDecls []*chunkDecl `parser:"'chunks':Keyword '{' @@* '}'"`

	Tokens []lexer.Token
}

type configSection struct {
	GactarConfig *gactarConfig `parser:"@@?"`
	ModuleConfig *moduleConfig `parser:"@@?"`
	ChunkConfig  *chunkConfig  `parser:"@@?"`

	Tokens []lexer.Token
}

type namedInitializer struct {
	ChunkName *string  `parser:"(@Ident)?"`
	Pattern   *pattern `parser:"@@"`

	Tokens []lexer.Token
}

type bufferInitializer struct {
	BufferName   string              `parser:"@Ident"`
	InitPatterns []*namedInitializer `parser:"( '{' @@+ '}' | @@ )"`

	Tokens []lexer.Token
}

type moduleInitializer struct {
	ModuleName         string               `parser:"@Ident"`
	InitPatterns       []*namedInitializer  `parser:"( '{' @@+ '}' | @@"`
	BufferInitPatterns []*bufferInitializer `parser:"| '{' @@+ '}' )"`

	Tokens []lexer.Token
}

type similar struct {
	OpenParen  string  `parser:"'('"`
	ChunkOne   string  `parser:"@Ident"`
	ChunkTwo   string  `parser:"@Ident"`
	Value      float64 `parser:"@Number"`
	CloseParen string  `parser:"')'"`

	Tokens []lexer.Token
}

type similarityInitializer struct {
	Similar     string    `parser:"'similar':Keyword"`
	OpenBrace   string    `parser:"'{'"`
	SimilarList []similar `parser:"@@+"`
	CloseBrace  string    `parser:"'}'"`

	Tokens []lexer.Token
}

type initialization struct {
	ModuleInitializer     *moduleInitializer     `parser:"( @@"`
	SimilarityInitializer *similarityInitializer `parser:"| @@ )"`

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
	Str      *string `parser:"| @String"`
	Num      *string `parser:"| @Number"` // we don't need to treat this as a number anywhere, so keep as a string
	Var      *string `parser:"| @Var ))"`
	Wildcard *string `parser:"| @Wildcard)"`

	Tokens []lexer.Token
}

type chunkPattern struct {
	Name  string         `parser:"@Ident ':'"`
	Slots []*patternSlot `parser:"@@+"`

	Tokens []lexer.Token
}

type pattern struct {
	StartBracket string `parser:"'['"` // not used - must be set for parse

	AnyChunk *string       `parser:"( @('any':Keyword)"`
	Chunk    *chunkPattern `parser:"| @@ )"`

	EndBracket string `parser:"']'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type comparisonOperator struct {
	Equal    *string `parser:"( @Equality"`
	NotEqual *string `parser:"| @Inequality )"`

	Tokens []lexer.Token
}

type whenExpression struct {
	OpenParen  string              `parser:"'('"` // not used - must be set for parse
	LHS        string              `parser:"@Var"`
	Comparison *comparisonOperator `parser:"@@"`
	RHS        *whenArg            `parser:"@@"`
	CloseParen string              `parser:"')'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type whenClause struct {
	When        string             `parser:"'when':Keyword"`
	Expressions *[]*whenExpression `parser:"@@ ('and' @@)*"`

	Tokens []lexer.Token
}

type matchBufferPatternItem struct {
	BufferName string      `parser:"@Ident"`
	Pattern    *pattern    `parser:"@@"`
	When       *whenClause `parser:"@@?"`

	Tokens []lexer.Token
}

type matchBufferStateItem struct {
	Keyword    string `parser:"'buffer_state':Keyword"`
	BufferName string `parser:"@Ident"`
	State      string `parser:"@Ident"`

	Tokens []lexer.Token
}

type matchModuleStateItem struct {
	Keyword    string `parser:"'module_state':Keyword"`
	ModuleName string `parser:"@Ident"`
	State      string `parser:"@Ident"`

	Tokens []lexer.Token
}

type matchItem struct {
	BufferPattern *matchBufferPatternItem `parser:"( @@"`
	BufferState   *matchBufferStateItem   `parser:"| @@"`
	ModuleState   *matchModuleStateItem   `parser:"| @@)"`

	Tokens []lexer.Token
}

type match struct {
	Items []*matchItem `parser:"'match':Keyword '{' @@+ '}'"`

	Tokens []lexer.Token
}

type clearStatement struct {
	BufferNames []string `parser:"'clear':Keyword ( @Ident ','? )+"`

	Tokens []lexer.Token
}

type printStatement struct {
	Args []*printArg `parser:"'print':Keyword ( @@ ','? )*"`

	Tokens []lexer.Token
}

type withExpression struct {
	OpenParen  string   `parser:"'('"` // not used - must be set for parse
	Param      string   `parser:"@Ident"`
	Value      *withArg `parser:"@@"`
	CloseParen string   `parser:"')'"` // not used - must be set for parse

	Tokens []lexer.Token
}

type withClause struct {
	With        string             `parser:"'with':Keyword"`
	Expressions *[]*withExpression `parser:"@@ ('and':Keyword @@)*"`

	Tokens []lexer.Token
}

type recallStatement struct {
	Pattern *pattern    `parser:"'recall':Keyword @@"`
	With    *withClause `parser:"@@?"`

	Tokens []lexer.Token
}

type setStatement struct {
	Set       string    `parser:"'set':Keyword"` // not used, but must be visible for parse to work
	BufferRef bufferRef `parser:"@@"`

	To      string   `parser:"'to':Keyword"` // not used, but must be visible for parse to work
	Value   *setArg  `parser:"( @@"`
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
	Do         string        `parser:"'do':Keyword"` // not used, but must be visible for parse to work
	Statements *[]*statement `parser:"'{' @@+ '}'"`

	Tokens []lexer.Token
}

type production struct {
	Name        string  `parser:"@Ident '{'"`
	Description *string `parser:"('description':Keyword ':' @String)?"`
	Match       *match  `parser:"@@"`
	Do          *do     `parser:"@@"`
	End         string  `parser:"'}'"` // not used, but must be visible for parse to work

	Tokens []lexer.Token
}

type productionSection struct {
	Productions []*production `parser:"@@+"`

	Tokens []lexer.Token
}

var amodParser = participle.MustBuild[amodFile](
	participle.Lexer(LexerDefinition),
	participle.Elide("Comment", "Whitespace"),
	participle.Unquote(),
)

func parseAMOD(r io.Reader) (amod *amodFile, err error) {
	amod, err = amodParser.Parse("", r)
	if err != nil {
		return nil, err
	}

	return
}

func parseAMODFile(filename string) (*amodFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseAMOD(file)
}

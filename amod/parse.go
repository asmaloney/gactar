package amod

import (
	"io"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Use participle to parse the lexemes.
// https://github.com/alecthomas/participle

type amodFile struct {
	Model       *modelSection      `parser:"'==model==' @@"`
	Config      *configSection     `parser:"'==config==' (@@)?"`
	Init        *initSection       `parser:"'==init==' (@@)?"`
	Productions *productionSection `parser:"'==productions==' (@@)?"`

	Pos lexer.Position
}

type modelSection struct {
	Name        string   `parser:"'name' ':' (@String|@Ident)"`
	Description string   `parser:"('description' ':' (@String|@Ident))?"`
	Examples    []string `parser:"('examples' '{' (@String)+ '}')?"`

	Pos lexer.Position
}

type identList struct {
	Identifiers []string `parser:"( @Ident ','? )+"`

	Pos lexer.Position
}

type stringList struct {
	Strings []string `parser:"( @String ','? )+"`

	Pos lexer.Position
}

type value struct {
	String *string  `parser:"  (@String|@Ident)"`
	Number *float64 `parser:"| @Number"`

	Pos lexer.Position
}

type field struct {
	Key   string `parser:"@Ident ':'"`
	Value value  `parser:"@@ (',')?"`

	Pos lexer.Position
}

type fieldList struct {
	Fields []*field `parser:"'{' @@+ '}'"`

	Pos lexer.Position
}

type memory struct {
	Name   string    `parser:"@Ident"`
	Fields fieldList `parser:"@@+"`

	Pos lexer.Position
}

type memoryList struct {
	Memories []*memory `parser:"'{' @@+ '}'"`

	Pos lexer.Position
}

type configSection struct {
	ACTR        *fieldList  `parser:"('actr' @@)?"`
	Buffers     *identList  `parser:"('buffers' '{' @@ '}')?"`
	Memories    *memoryList `parser:"('memories' @@)?"`
	TextOutputs *identList  `parser:"('text_outputs' '{' @@ '}')?"`

	Pos lexer.Position
}

type initializer struct {
	Name  string      `parser:"@Ident"`
	Items *stringList `parser:"'{' @@+ '}'"`

	Pos lexer.Position
}

type initSection struct {
	Initializers []*initializer `parser:"@@+"`

	Pos lexer.Position
}

type matchItem struct {
	Name string `parser:"@Ident ':'"`
	Text string `parser:"@String"`

	Pos lexer.Position
}

type match struct {
	Items []*matchItem `parser:"'match' '{' ( @@ )+ '}'"`

	Pos lexer.Position
}

type do struct {
	Texts []string `parser:"'do' '#<' (@DoCode)+ '>#'"`

	Pos lexer.Position
}

type production struct {
	Name  string `parser:"@Ident '{'"`
	Match *match `parser:"@@"`
	Do    *do    `parser:"@@"`
	End   string `parser:"'}'"` // not used, but must be visible for parse to work

	Pos lexer.Position
}

type productionSection struct {
	Productions []*production `parser:"( @@ )+"`

	Pos lexer.Position
}

var parser = participle.MustBuild(&amodFile{},
	participle.Lexer(LexerDefinition),
	participle.Unquote(),
	participle.Elide("Comment", "Whitespace"),
)

func parse(r io.Reader) (*amodFile, error) {
	var amod amodFile

	err := parser.Parse("", r, &amod)

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

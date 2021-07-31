package amod

import (
	"io"
	"os"

	"github.com/alecthomas/participle/v2"
)

// Use participle to parse the lexemes.
// https://github.com/alecthomas/participle

type amodFile struct {
	Model       *modelSection      `parser:"'==model==' @@"`
	Config      *configSection     `parser:"'==config==' @@"`
	Init        *initSection       `parser:"'==init==' @@"`
	Productions *productionSection `parser:"'==productions==' @@"`
}

type modelSection struct {
	Name        string   `parser:"'name' ':' (@String|@Ident)"`
	Description string   `parser:"('description' ':' (@String|@Ident))?"`
	Examples    []string `parser:"('examples' '{' (@String)+ '}')?"`
}

type identList struct {
	Identifiers []string `parser:"( @Ident ','? )+"`
}

type value struct {
	String *string  `parser:"  (@String|@Ident)"`
	Number *float64 `parser:"| @Number"`
}

type field struct {
	Key   string `parser:"@Ident ':'"`
	Value value  `parser:"@@ (',')?"`
}

type fieldList struct {
	Fields []*field `parser:"'{' @@+ '}'"`
}

type memory struct {
	Name   string    `parser:"@Ident"`
	Fields fieldList `parser:"@@+"`
}

type memoryList struct {
	Memories []*memory `parser:"'{' @@+ '}'"`
}

type configSection struct {
	ACTR        *fieldList  `parser:"('actr' @@)?"`
	Buffers     *identList  `parser:"('buffers' '{' @@ '}')?"`
	Memories    *memoryList `parser:"('memories' @@)?"`
	TextOutputs *identList  `parser:"('text_outputs' '{' @@ '}')?"`
}

type initItem struct {
	Item string `parser:"'{' @String '}'"`
}

type initializer struct {
	Name  string      `parser:"@Ident '{'"`
	Items []*initItem `parser:"@@+"`
	End   string      `parser:"'}'"`
}

type initSection struct {
	Initializers []*initializer `parser:"@@+"`
}

type matchItem struct {
	Name string `parser:"@Ident ':'"`
	Text string `parser:"@String"`
}

type match struct {
	Items []*matchItem `parser:"'match' '{' ( @@ )+ '}'"`
}

type do struct {
	Texts []string `parser:"'do' '#<' (@DoCode)+ '>#'"`
}

type production struct {
	Name  string `parser:"@Ident '{'"`
	Match *match `parser:"@@"`
	Do    *do    `parser:"@@"`
	End   string `parser:"'}'"`
}

type productionSection struct {
	Productions []*production `parser:"( @@ )+"`
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

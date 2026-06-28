package nogit

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	nogitLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Myself", Pattern: "watashi"},
		{Name: "RightDirection", Pattern: "migi"},
		{Name: "LeftDirection", Pattern: "hidari"},
		{Name: "EOL", Pattern: `[\n\r]+`},
		{Name: "Whitespace", Pattern: `[ \t]+`},
	})

	nogitParser = participle.MustBuild[NogitAST](
		participle.Lexer(nogitLexer),
		participle.Elide("Whitespace", "EOL"),
	)
)

func Parse(text string) (*NogitAST, error) {
	json, err := nogitParser.ParseString("from in game", text)
	if err != nil {
		return nil, err
	}

	return json, nil
}

type NogitAST struct {
	Noun *Noun `parser:"@@"`
}

type Noun struct {
	Myself         *string `parser:"@Myself"`
	RightDirection *string `parser:"@RightDirection"`
	LeftDirection  *string `parser:"@LeftDirection"`
}

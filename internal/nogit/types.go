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
		{Name: "Move", Pattern: "ugoku"},

		{Name: "Topic", Pattern: "ha"},
		{Name: "Direction", Pattern: "ni"},

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
	Words []*Word `parser:"@@*"`
}

type Word struct {
	Noun        string `parser:"@(Myself | RightDirection | LeftDirection | Move) |"`
	Grammatical string `parser:"@(Topic | Direction)"`
}

package nogit

import (
	"testing"
)

func testParseFromText(t *testing.T, text string, expected *NogitAST) {
	actual, err := Parse(text)
	if err != nil {
		t.Errorf("An error invoked from Parse func: %v", err)
	}

	if actual == nil {
		t.Error("The actual variable created by Parse func is nil")
	}

	if len(expected.Words) != len(actual.Words) {
		t.Errorf("The expected words' len is %d, but got %d", len(expected.Words), len(actual.Words))
	}

	for i := range expected.Words {
		expectedWord := expected.Words[i]
		actualWord := actual.Words[i]

		if expectedWord.Grammatical != "" {
			if expectedWord.Grammatical != actualWord.Grammatical {
				t.Errorf("Expected %#v, but got %#v", expectedWord.Grammatical, actualWord.Grammatical)
			}
		} else {
			if expectedWord.Noun != actualWord.Noun {
				t.Errorf("Expected %#v, but got %#v", expectedWord.Noun, actualWord.Noun)
			}
		}
	}
}

func TestParse(t *testing.T) {
	testParseFromText(t, "watashi migi hidari", &NogitAST{Words: []*Word{
		{Noun: "watashi"},
		{Noun: "migi"},
		{Noun: "hidari"},
	}})

	testParseFromText(t, "watashi ha migi ni", &NogitAST{Words: []*Word{
		{Noun: "watashi"},
		{Grammatical: "ha"},
		{Noun: "migi"},
		{Grammatical: "ni"},
	}})

	testParseFromText(t, "watashi ha hidari ni ugoku", &NogitAST{Words: []*Word{
		{Noun: "watashi"},
		{Grammatical: "ha"},
		{Noun: "hidari"},
		{Grammatical: "ni"},
		{Noun: "ugoku"},
	}})

	testParseFromText(t, "watashihahidariniugoku", &NogitAST{Words: []*Word{
		{Noun: "watashi"},
		{Grammatical: "ha"},
		{Noun: "hidari"},
		{Grammatical: "ni"},
		{Noun: "ugoku"},
	}})
}

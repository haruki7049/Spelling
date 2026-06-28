package nogit

import (
	"testing"
)

func TestParse(t *testing.T) {
	{
		text := "watashi migi hidari"
		result, err := Parse(text)

		if err != nil {
			t.Errorf("An error invoked from Parse func: %v", err)
		}

		if result == nil {
			t.Error("The result by Parse func is nil")
		}
	}

	{
		text := "watashi ha migi ni"
		result, err := Parse(text)

		if err != nil {
			t.Errorf("An error invoked from Parse func: %v", err)
		}

		if result == nil {
			t.Error("The result by Parse func is nil")
		}
	}

	{
		text := "watashi ha hidari ni ugoku"
		result, err := Parse(text)

		if err != nil {
			t.Errorf("An error invoked from Parse func: %v", err)
		}

		if result == nil {
			t.Error("The result by Parse func is nil")
		}
	}

	{
		text := "watashihahidariniugoku"
		actual, err := Parse(text)

		if err != nil {
			t.Errorf("An error invoked from Parse func: %v", err)
		}

		if actual == nil {
			t.Error("The result by Parse func is nil")
		}

		expected := &NogitAST{
			Words: []*Word{
				&Word{Noun: "watashi"},
				&Word{Grammatical: "ha"},
				&Word{Noun: "hidari"},
				&Word{Grammatical: "ni"},
				&Word{Noun: "ugoku"},
			},
		}

		if len(expected.Words) != len(actual.Words) {
			t.Errorf("The expected words' len is %d, but got %d", len(expected.Words), len(actual.Words))
		}

		for i, _ := range expected.Words {
			expectedWord := expected.Words[i]
			actualWord := actual.Words[i]

			if expectedWord != actualWord {
				t.Errorf("Expected %#v, but got %#v", expectedWord, actualWord)
			}
		}

		if expected != actual {
			t.Errorf("Expected %#v, but got %#v", expected, actual)
		}
	}
}

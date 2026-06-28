package nogit

import (
	"testing"
)

func TestParse(t *testing.T) {
	text := "watashi migi hidari"
	result, err := Parse(text)

	if err != nil {
		t.Errorf("An error invoked from Parse func: %v", err)
	}

	if result == nil {
		t.Error("The result by Parse func is nil")
	}
}

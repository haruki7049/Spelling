package asm

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// TestLLVMGolden checks that the assembler matches LLVM's RISC-V assembler
// on every case in testdata/llvm.golden: the same machine code where LLVM
// accepts the source, and an error where LLVM rejects it. Regenerate the
// file with testdata/gen-llvm-golden.sh after changing testdata/llvm-cases.txt.
func TestLLVMGolden(t *testing.T) {
	data, err := os.ReadFile("testdata/llvm.golden")
	if err != nil {
		t.Fatal(err)
	}
	_, body, _ := strings.Cut(string(data), "\n") // skip the header comment
	cases := strings.Split(strings.TrimSpace(body), "\n---\n")
	if len(cases) < 100 {
		t.Fatalf("only %d cases in llvm.golden", len(cases))
	}
	for _, c := range cases {
		i := strings.LastIndex(c, "\n=> ")
		if i < 0 {
			t.Fatalf("malformed case %q", c)
		}
		src, want := c[:i], c[i+len("\n=> "):]
		got, err := Assemble(src, 0x200)
		if want == "error" {
			if err == nil {
				t.Errorf("%q: LLVM rejects it, but got %#08x", src, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q: LLVM accepts it, but got error: %v", src, err)
			continue
		}
		var wantWords []uint32
		for w := range strings.FieldsSeq(want) {
			v, err := strconv.ParseUint(w, 16, 32)
			if err != nil {
				t.Fatalf("%q: malformed word %q", src, w)
			}
			wantWords = append(wantWords, uint32(v))
		}
		if !slices.Equal(got, wantWords) {
			t.Errorf("%q:\n got %s\nwant %s", src, fmt.Sprintf("%08x", got), want)
		}
	}
}

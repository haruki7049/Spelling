package asm

import (
	"slices"
	"strings"
	"testing"

	"github.com/haruki7049/spelling/internal/vm"
)

func TestAssembleEncoding(t *testing.T) {
	// Expected words match llvm-mc -triple=riscv32, except that li always
	// expands to lui+addi when the value does not fit in 12 bits.
	tests := []struct {
		src  string
		want []uint32
	}{
		{"addi a0, zero, 1", []uint32{0x00100513}},
		{"add a0, a1, a2", []uint32{0x00c58533}},
		{"sub x5, x6, x7", []uint32{0x407302b3}},
		{"srai a0, a0, 3", []uint32{0x40355513}},
		{"lui a0, 0x12345", []uint32{0x12345537}},
		{"lw a0, 8(sp)", []uint32{0x00812503}},
		{"lbu t0, (a1)", []uint32{0x0005c283}},
		{"sw a0, 8(sp)", []uint32{0x00a12423}},
		{"sh a0, -2(sp)", []uint32{0xfea11f23}},
		{"beq a0, a1, 8", []uint32{0x00b50463}},
		{"jal zero, 0", []uint32{0x0000006f}},
		{"jalr ra, 0(a0)", []uint32{0x000500e7}},
		{"ecall", []uint32{0x00000073}},
		{"ebreak", []uint32{0x00100073}},
		{"fence", []uint32{0x0ff0000f}},
		{"nop", []uint32{0x00000013}},
		{"ret", []uint32{0x00008067}},
		{"mv a0, a1", []uint32{0x00058513}},
		{"li a0, -1", []uint32{0xfff00513}},
		{"li a0, 0xffffffff", []uint32{0xfff00513}},
		{"li a0, 0x12345678", []uint32{0x12345537, 0x67850513}},
		{"li a0, 0x800", []uint32{0x00001537, 0x80050513}},
		{"li a0, 0x1000_2008", []uint32{0x10002537, 0x00850513}},
		{"ADDI A0, ZERO, 1", []uint32{0x00100513}},
		{"top: j top", []uint32{0x0000006f}},
		{"nop; nop # comment; not a statement", []uint32{0x13, 0x13}},
	}
	for _, tt := range tests {
		got, err := Assemble(tt.src, 0)
		if err != nil {
			t.Errorf("%q: %v", tt.src, err)
			continue
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("%q = %#08x, want %#08x", tt.src, got, tt.want)
		}
	}
}

// run assembles src at base, runs it on a CPU until it reaches the "done"
// label, and returns the CPU.
func run(t *testing.T, src string, base uint32) *vm.CPU {
	t.Helper()
	code, err := Assemble(src, base)
	if err != nil {
		t.Fatal(err)
	}
	ram := vm.NewRAM(4096)
	for i, w := range code {
		ram.Write(base+uint32(4*i), 4, w)
	}
	done := base + uint32(4*(len(code)-1)) // the last statement is "done: j done"
	c := vm.NewCPU(ram, base)
	for range 1000 {
		if c.PC == done {
			return c
		}
		if !c.Step() {
			t.Fatalf("illegal instruction at %#x", c.PC)
		}
	}
	t.Fatal("program did not finish")
	return nil
}

func TestAssembleRun(t *testing.T) {
	// Sum 1..10 in one line, as a player would type it.
	c := run(t, "li a0, 0; li t0, 10; loop: add a0, a0, t0; addi t0, t0, -1; bnez t0, loop; done: j done", 0x100)
	if c.Regs[10] != 55 {
		t.Errorf("sum = %d, want 55", c.Regs[10])
	}

	c = run(t, `
		li sp, 0x800
		la a1, value     # address of the data word
		call double
		j done
	double:
		lw a0, 0(a1)
		add a0, a0, a0
		sw a0, 0(a1)
		ret
	value:
		nop              # data word 0x13
	done: j done
	`, 0x200)
	if c.Regs[10] != 0x26 {
		t.Errorf("a0 = %#x, want 0x26", c.Regs[10])
	}
	if c.Bus.Read(c.Regs[11], 4) != 0x26 {
		t.Errorf("value not stored at la address %#x", c.Regs[11])
	}

	c = run(t, "li t0, 0x12345678; li t1, -2048; li t2, 0x7ff; beqz zero, done; li t0, 0; done: j done", 0)
	if c.Regs[5] != 0x12345678 || c.Regs[6] != 0xfffff800 || c.Regs[7] != 0x7ff {
		t.Errorf("li results = %#x %#x %#x", c.Regs[5], c.Regs[6], c.Regs[7])
	}
}

func TestAssembleErrors(t *testing.T) {
	tests := []struct {
		src  string
		line int
		want string
	}{
		{"fly a0", 1, "unknown instruction"},
		{"nop\naddi a0, a0", 2, "want 3 operands"},
		{"addi a0, a0, 2048", 1, "out of range"},
		{"slli a0, a0, 32", 1, "out of range"},
		{"lui a0, 0x100000", 1, "out of range"},
		{"add a0, a1, x32", 1, "invalid register"},
		{"lw a0, 8", 1, "memory operand"},
		{"j nowhere", 1, "undefined label"},
		{"j 3", 1, "multiple of 2"},
		{"beq a0, a1, 4096", 1, "out of range"},
		{"a: nop\na: nop", 2, "duplicate label"},
		{"1x: nop", 1, "invalid label"},
		{"li a0, 0x100000000", 1, "out of range"},
		{"la a0, 16", 1, "want a label"},
		{"ret a0", 1, "want 0 operands"},
	}
	for _, tt := range tests {
		_, err := Assemble(tt.src, 0)
		e, ok := err.(*Error)
		if !ok || e.Line != tt.line || !strings.Contains(e.Msg, tt.want) {
			t.Errorf("%q: err = %v, want line %d containing %q", tt.src, err, tt.line, tt.want)
		}
	}
}

func TestAssembleEmpty(t *testing.T) {
	code, err := Assemble("  # only a comment\n\n;;", 0)
	if err != nil || len(code) != 0 {
		t.Errorf("got %v, %v; want no code", code, err)
	}
}

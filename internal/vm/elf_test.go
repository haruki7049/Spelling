package vm

import (
	"debug/elf"
	"encoding/binary"
	"strings"
	"testing"
)

type testSegment struct {
	vaddr, memsz uint32
	flags        elf.ProgFlag
	data         []byte
}

// buildELF returns a minimal ELF32 little-endian RISC-V executable with one
// PT_LOAD program header per segment.
func buildELF(entry uint32, segs ...testSegment) []byte {
	const ehsize, phentsize = 52, 32
	le := binary.LittleEndian
	off := uint32(ehsize + phentsize*len(segs))

	b := make([]byte, ehsize)
	copy(b, elf.ELFMAG)
	b[elf.EI_CLASS] = byte(elf.ELFCLASS32)
	b[elf.EI_DATA] = byte(elf.ELFDATA2LSB)
	b[elf.EI_VERSION] = byte(elf.EV_CURRENT)
	le.PutUint16(b[16:], uint16(elf.ET_EXEC))
	le.PutUint16(b[18:], uint16(elf.EM_RISCV))
	le.PutUint32(b[20:], uint32(elf.EV_CURRENT))
	le.PutUint32(b[24:], entry)
	le.PutUint32(b[28:], ehsize) // e_phoff
	le.PutUint16(b[40:], ehsize)
	le.PutUint16(b[42:], phentsize)
	le.PutUint16(b[44:], uint16(len(segs)))

	var body []byte
	for _, s := range segs {
		ph := make([]byte, phentsize)
		le.PutUint32(ph[0:], uint32(elf.PT_LOAD))
		le.PutUint32(ph[4:], off)
		le.PutUint32(ph[8:], s.vaddr)
		le.PutUint32(ph[12:], s.vaddr)
		le.PutUint32(ph[16:], uint32(len(s.data)))
		le.PutUint32(ph[20:], s.memsz)
		le.PutUint32(ph[24:], uint32(s.flags))
		b = append(b, ph...)
		body = append(body, s.data...)
		off += uint32(len(s.data))
	}
	return append(b, body...)
}

func TestLoadELF(t *testing.T) {
	ram := NewRAM(4096)
	ram[0x204] = 0xaa                         // inside .bss; must be zero-filled
	ram[0x300] = 0xbb                         // outside every segment; must be kept
	program := []byte{0x13, 0x05, 0x10, 0x00} // addi a0, zero, 1
	data := buildELF(0x100,
		testSegment{vaddr: 0x100, memsz: 4, flags: elf.PF_R | elf.PF_X, data: program},
		testSegment{vaddr: 0x200, memsz: 8, flags: elf.PF_R | elf.PF_W, data: []byte{1, 2, 3, 4}},
	)

	entry, err := LoadELF(data, ram)
	if err != nil {
		t.Fatal(err)
	}
	if entry != 0x100 {
		t.Errorf("entry = %#x, want 0x100", entry)
	}
	if got := ram.Read(0x100, 4); got != 0x00100513 {
		t.Errorf("code = %#08x, want 0x00100513", got)
	}
	if got := ram.Read(0x200, 4); got != 0x04030201 {
		t.Errorf("data = %#08x, want 0x04030201", got)
	}
	if got := ram.Read(0x204, 4); got != 0 {
		t.Errorf(".bss = %#08x, want 0", got)
	}
	if ram[0x300] != 0xbb {
		t.Errorf("byte outside segments was modified")
	}

	c := NewCPU(ram, entry)
	step(t, c, 1)
	if c.Regs[10] != 1 {
		t.Errorf("a0 = %d, want 1", c.Regs[10])
	}
}

func TestLoadELFRejects(t *testing.T) {
	code := testSegment{vaddr: 0, memsz: 4, flags: elf.PF_R | elf.PF_X, data: make([]byte, 4)}
	valid := buildELF(0, code)

	patch := func(off int, v uint16) []byte {
		b := append([]byte(nil), valid...)
		binary.LittleEndian.PutUint16(b[off:], v)
		return b
	}
	class64 := append([]byte(nil), valid...)
	class64[elf.EI_CLASS] = byte(elf.ELFCLASS64)

	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"not ELF", []byte("hello"), "elf:"},
		{"64-bit", class64, "elf:"},
		{"wrong machine", patch(18, uint16(elf.EM_X86_64)), "EM_RISCV"},
		{"not executable", patch(16, uint16(elf.ET_DYN)), "ET_EXEC"},
		{"segment past end of RAM", buildELF(0, testSegment{vaddr: 4094, memsz: 4, flags: elf.PF_X, data: make([]byte, 4)}), "does not fit"},
		{"segment address overflow", buildELF(0, code, testSegment{vaddr: 0xffff_fff0, memsz: 0x20, flags: elf.PF_R}), "does not fit"},
		{"huge memory size", buildELF(0, testSegment{vaddr: 0, memsz: 0xffff_ffff, flags: elf.PF_X}), "does not fit"},
		{"file size exceeds memory size", buildELF(0, testSegment{vaddr: 0, memsz: 2, flags: elf.PF_X, data: make([]byte, 4)}), "exceeds memory size"},
		{"truncated segment data", buildELF(0, code)[:len(valid)-2], "segment 0"},
		{"entry outside segments", buildELF(0x800, code), "entry point"},
		{"entry in non-executable segment", buildELF(0, testSegment{vaddr: 0, memsz: 4, flags: elf.PF_R, data: make([]byte, 4)}), "entry point"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadELF(tt.data, NewRAM(4096))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("err = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func FuzzLoadELF(f *testing.F) {
	f.Add(buildELF(0, testSegment{vaddr: 0, memsz: 8, flags: elf.PF_X, data: make([]byte, 4)}))
	f.Add([]byte("\x7fELF"))
	f.Fuzz(func(t *testing.T, data []byte) {
		ram := NewRAM(4096)
		entry, err := LoadELF(data, ram)
		if err == nil && entry >= uint32(len(ram)) {
			t.Errorf("entry %#x outside RAM", entry)
		}
	})
}

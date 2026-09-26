package vm

import (
	"bytes"
	"debug/elf"
	"errors"
	"fmt"
	"io"
)

// LoadELF copies the loadable segments of a 32-bit little-endian RISC-V
// executable into ram and returns its entry point.
//
// The ELF file is untrusted input (see issue #15): every segment must fit
// inside ram, and the entry point must lie inside a loaded executable segment.
// Parts of ram not covered by a segment are left untouched; the part of a
// segment beyond its file data (for example .bss) is zero-filled. On error,
// ram may have been partially written.
func LoadELF(data []byte, ram RAM) (entry uint32, err error) {
	f, err := elf.NewFile(bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("elf: %w", err)
	}
	defer f.Close()

	if f.Class != elf.ELFCLASS32 {
		return 0, fmt.Errorf("elf: class %v, want ELFCLASS32", f.Class)
	}
	if f.Data != elf.ELFDATA2LSB {
		return 0, fmt.Errorf("elf: data encoding %v, want ELFDATA2LSB", f.Data)
	}
	if f.Machine != elf.EM_RISCV {
		return 0, fmt.Errorf("elf: machine %v, want EM_RISCV", f.Machine)
	}
	if f.Type != elf.ET_EXEC {
		return 0, fmt.Errorf("elf: type %v, want ET_EXEC", f.Type)
	}

	size := uint64(len(ram))
	entryOK := false
	for i, p := range f.Progs {
		if p.Type != elf.PT_LOAD || p.Memsz == 0 {
			continue
		}
		if p.Filesz > p.Memsz {
			return 0, fmt.Errorf("elf: segment %d: file size %d exceeds memory size %d", i, p.Filesz, p.Memsz)
		}
		if p.Vaddr > size || p.Memsz > size-p.Vaddr {
			return 0, fmt.Errorf("elf: segment %d: [%#x, +%#x) does not fit in %#x bytes of RAM", i, p.Vaddr, p.Memsz, size)
		}
		seg := ram[p.Vaddr : p.Vaddr+p.Memsz]
		if _, err := io.ReadFull(p.Open(), seg[:p.Filesz]); err != nil {
			return 0, fmt.Errorf("elf: segment %d: %w", i, err)
		}
		clear(seg[p.Filesz:])
		if p.Flags&elf.PF_X != 0 && f.Entry >= p.Vaddr && f.Entry < p.Vaddr+p.Memsz {
			entryOK = true
		}
	}
	if !entryOK {
		return 0, errors.New("elf: entry point is not inside a loaded executable segment")
	}
	return uint32(f.Entry), nil
}

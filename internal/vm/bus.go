package vm

// Bus is the memory seen by a CPU. Addresses are byte addresses and
// multi-byte values are little-endian, as in RISC-V.
//
// size is 1, 2, or 4 bytes. Read returns the value zero-extended to 32 bits.
// Accesses need not be aligned.
type Bus interface {
	Read(addr uint32, size int) uint32
	Write(addr uint32, size int, value uint32)
}

// RAM is a Bus backed by a byte slice starting at address 0.
// Reading outside the slice returns 0 and writing outside it is ignored,
// byte by byte.
type RAM []byte

// NewRAM returns a zero-filled RAM of the given size in bytes.
func NewRAM(size int) RAM {
	return make(RAM, size)
}

func (r RAM) Read(addr uint32, size int) uint32 {
	var v uint32
	for i := range size {
		a := uint64(addr) + uint64(i)
		if a < uint64(len(r)) {
			v |= uint32(r[a]) << (8 * i)
		}
	}
	return v
}

func (r RAM) Write(addr uint32, size int, value uint32) {
	for i := range size {
		a := uint64(addr) + uint64(i)
		if a < uint64(len(r)) {
			r[a] = byte(value >> (8 * i))
		}
	}
}

package h264

import (
	"encoding/binary"
	"math/bits"
)

func readExponentialGolombCoding(r *bitReader) (uint64, error) {
	x := uint8(0)
	for {
		b, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if b == BitOne {
			break
		}
		x++
	}
	g := uint64(1) << x
	for i := uint8(0); i < x; i++ {
		b, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if b == BitZero {
			continue
		}
		g |= 1 << (x - i - 1)
	}
	return g - 1, nil
}

func writeExponentialGolombCoding(w *bitWriter, code uint64) (writtenBit int, err error) {
	g := code + 1
	x := bits.Len64(g) - 1
	for i := 0; i < x; i++ {
		n, err := w.WriteBit(BitZero)
		if err != nil {
			return writtenBit, err
		}
		writtenBit += n
	}

	n, err := w.WriteBit(BitOne)
	if err != nil {
		return writtenBit, err
	}
	writtenBit += n

	if g == 1 {
		return writtenBit, nil
	}

	bb := make([]byte, 8)
	binary.BigEndian.PutUint64(bb, g)
	r := newBitReader(bb)
	if err := r.SeekBit(64 - x); err != nil {
		return writtenBit, nil
	}

	for i := 0; i < x; i++ {
		b, err := r.ReadBit()
		if err != nil {
			return writtenBit, err
		}
		n, err := w.WriteBit(b)
		if err != nil {
			return writtenBit, err
		}
		writtenBit += n
	}
	return writtenBit, nil
}

func Uint64ToGolombCodeNum(s uint64) uint64 {
	return s
}

func GolombCodeNumToUint64(codeNum uint64) uint64 {
	return codeNum
}

func Int64ToGolombCodeNum(s int64) uint64 {
	if s == 0 {
		return 0
	}
	if s > 0 {
		return uint64(2*s - 1)
	}
	return uint64(-2 * s)
}

func GolombCodeNumToInt64(codeNum uint64) int64 {
	if codeNum == 0 {
		return 0
	}
	if codeNum%2 == 1 {
		return int64((codeNum + 1) / 2)
	}
	return -int64(codeNum / 2)
}

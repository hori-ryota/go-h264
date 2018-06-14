package h264

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var GolombCodeNumTestData = []struct {
	Code   uint64
	Binary []byte
}{
	{
		Code:   0,
		Binary: []byte{0x80}, // 0b1 0000000
	},
	{
		Code:   0,
		Binary: []byte{0xff}, // 0b1 1111111
	},
	{
		Code:   1,
		Binary: []byte{0x40}, // 0b010 00000
	},
	{
		Code:   1,
		Binary: []byte{0x5f}, // 0b010 11111
	},
	{
		Code:   2,
		Binary: []byte{0x60}, // 0b011 00000
	},
	{
		Code:   2,
		Binary: []byte{0x7f}, // 0b011 11111
	},
	{
		Code:   3,
		Binary: []byte{0x20}, // 0b00100 000
	},
	{
		Code:   3,
		Binary: []byte{0x27}, // 0b00100 111
	},
}

func Test_readExponentialGolombCoding(t *testing.T) {
	for _, tt := range GolombCodeNumTestData {
		t.Run(fmt.Sprintf("%x to %d", tt.Binary, tt.Code), func(t *testing.T) {
			r := newBitReader(tt.Binary)
			got, err := readExponentialGolombCoding(r)
			require.NoError(t, err)
			assert.Equal(t, tt.Code, got)
		})
	}
}

func Test_writeExponentialGolombCoding(t *testing.T) {
	for _, tt := range GolombCodeNumTestData {
		t.Run(fmt.Sprintf("%d to %x", tt.Code, tt.Binary), func(t *testing.T) {
			w := newBitWriterSize(len(tt.Binary))
			writtenBit, err := writeExponentialGolombCoding(w, tt.Code)
			require.NoError(t, err)
			want := tt.Binary
			if writtenBit%8 != 0 {
				want[len(want)-1] &= 0xff ^ (1<<uint8(8-(writtenBit%8)) - 1)
			}
			assert.Equal(t, want, w.Bytes())
		})
	}
}

func TestUint64ToGolombCodeNum(t *testing.T) {
	for _, tt := range []struct {
		name string
		s    uint64
		d    uint64
	}{
		{
			name: "0 to 0",
			s:    0,
			d:    0,
		},
		{
			name: "1 to 1",
			s:    1,
			d:    1,
		},
		{
			name: "2 to 2",
			s:    2,
			d:    2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.d, Uint64ToGolombCodeNum(tt.s))
		})
	}
}

func TestGolombCodeNumToUint64(t *testing.T) {
	for _, tt := range []struct {
		name string
		s    uint64
		d    uint64
	}{
		{
			name: "0 to 0",
			s:    0,
			d:    0,
		},
		{
			name: "1 to 1",
			s:    1,
			d:    1,
		},
		{
			name: "2 to 2",
			s:    2,
			d:    2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.d, GolombCodeNumToUint64(tt.s))
		})
	}
}

func TestInt64ToGolombCodeNum(t *testing.T) {
	for _, tt := range []struct {
		name string
		s    int64
		d    uint64
	}{
		{
			name: "0 to 0",
			s:    0,
			d:    0,
		},
		{
			name: "1 to 1",
			s:    1,
			d:    1,
		},
		{
			name: "-1 to 2",
			s:    -1,
			d:    2,
		},
		{
			name: "2 to 3",
			s:    2,
			d:    3,
		},
		{
			name: "-2 to 4",
			s:    -2,
			d:    4,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.d, Int64ToGolombCodeNum(tt.s))
		})
	}
}

func TestGolombCodeNumToInt64(t *testing.T) {
	for _, tt := range []struct {
		name string
		s    uint64
		d    int64
	}{
		{
			name: "0 to 0",
			s:    0,
			d:    0,
		},
		{
			name: "1 to 1",
			s:    1,
			d:    1,
		},
		{
			name: "2 to -1",
			s:    2,
			d:    -1,
		},
		{
			name: "3 to 2",
			s:    3,
			d:    2,
		},
		{
			name: "4 to -2",
			s:    4,
			d:    -2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.d, GolombCodeNumToInt64(tt.s))
		})
	}
}

package h264

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var HypotheticalReferenceDecoderTestData = []struct {
	Name   string
	Struct HypotheticalReferenceDecoder
	Binary []byte
}{
	{
		Name: "minimal struct",
		Struct: HypotheticalReferenceDecoder{
			BitRateValueMinus1: make([]uint64, 1),
			CPBSizeValueMinus1: make([]uint64, 1),
			CBRFlag:            make([]bool, 1),
		},
		Binary: mustBitToBytes(
			l,          // CPBCntMinus1
			o, o, o, o, // BitRateScale
			o, o, o, o, // CPBSizeScale

			l, // BitRateValueMinus1[0]
			l, // CPBSizeValueMinus1[0]
			o, // CBRFlag[0]

			o, o, o, o, o, // InitialCPBRemovalDelayLengthMinus1
			o, o, o, o, o, // CPBRemovalDelayLengthMinus1
			o, o, o, o, o, // DPBOutputDelayLengthMinus1
			o, o, o, o, o, // TimeOffsetLength
		),
	},
	{
		Name: "BitRateScale, CPBSizeScale",
		Struct: HypotheticalReferenceDecoder{
			BitRateScale: 0xa,
			CPBSizeScale: 0x5,

			BitRateValueMinus1: make([]uint64, 1),
			CPBSizeValueMinus1: make([]uint64, 1),
			CBRFlag:            make([]bool, 1),
		},
		Binary: mustBitToBytes(
			l,          // CPBCntMinus1
			l, o, l, o, // BitRateScale
			o, l, o, l, // CPBSizeScale

			l, // BitRateValueMinus1[0]
			l, // CPBSizeValueMinus1[0]
			o, // CBRFlag[0]

			o, o, o, o, o, // InitialCPBRemovalDelayLengthMinus1
			o, o, o, o, o, // CPBRemovalDelayLengthMinus1
			o, o, o, o, o, // DPBOutputDelayLengthMinus1
			o, o, o, o, o, // TimeOffsetLength
		),
	},
	{
		Name: "CPBCntMinus1 is 1",
		Struct: HypotheticalReferenceDecoder{
			CPBCntMinus1: 1,
			BitRateValueMinus1: []uint64{
				1, 3,
			},
			CPBSizeValueMinus1: []uint64{
				2, 4,
			},
			CBRFlag: []bool{
				true, false,
			},
		},
		Binary: mustBitToBytes(
			o, l, o, // CPBCntMinus1
			o, o, o, o, // BitRateScale
			o, o, o, o, // CPBSizeScale

			o, l, o, // BitRateValueMinus1[0]
			o, l, l, // CPBSizeValueMinus1[0]
			l, // CBRFlag[0]

			o, o, l, o, o, // BitRateValueMinus1[0]
			o, o, l, o, l, // CPBSizeValueMinus1[0]
			o, // CBRFlag[0]

			o, o, o, o, o, // InitialCPBRemovalDelayLengthMinus1
			o, o, o, o, o, // CPBRemovalDelayLengthMinus1
			o, o, o, o, o, // DPBOutputDelayLengthMinus1
			o, o, o, o, o, // TimeOffsetLength
		),
	},
	{
		Name: "InitialCPBRemovalDelayLengthMinus1, CPBRemovalDelayLengthMinus1, DPBOutputDelayLengthMinus1, TimeOffsetLength",
		Struct: HypotheticalReferenceDecoder{
			BitRateValueMinus1: make([]uint64, 1),
			CPBSizeValueMinus1: make([]uint64, 1),
			CBRFlag:            make([]bool, 1),

			InitialCPBRemovalDelayLengthMinus1: 0x15,
			CPBRemovalDelayLengthMinus1:        0x0a,
			DPBOutputDelayLengthMinus1:         0x15,
			TimeOffsetLength:                   0x0a,
		},
		Binary: mustBitToBytes(
			l,          // CPBCntMinus1
			o, o, o, o, // BitRateScale
			o, o, o, o, // CPBSizeScale

			l, // BitRateValueMinus1[0]
			l, // CPBSizeValueMinus1[0]
			o, // CBRFlag[0]

			l, o, l, o, l, // InitialCPBRemovalDelayLengthMinus1
			o, l, o, l, o, // CPBRemovalDelayLengthMinus1
			l, o, l, o, l, // DPBOutputDelayLengthMinus1
			o, l, o, l, o, // TimeOffsetLength
		),
	},
}

func TestHypotheticalReferenceDecoder_MarshalBinary(t *testing.T) {
	for _, tt := range HypotheticalReferenceDecoderTestData {
		t.Run(tt.Name, func(t *testing.T) {
			w := newBitWriter()
			err := writeHypotheticalReferenceDecoder(w, tt.Struct)
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, w.Bytes())
		})
	}
}

func TestHypotheticalReferenceDecoder_UnmarshalBinary(t *testing.T) {
	for _, tt := range HypotheticalReferenceDecoderTestData {
		t.Run(tt.Name, func(t *testing.T) {
			r := newBitReader(tt.Binary)
			s, err := readHypotheticalReferenceDecoder(r)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

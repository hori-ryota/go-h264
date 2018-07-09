package h264

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var VideoUsabilityInformationTestData = []struct {
	Name   string
	Struct VideoUsabilityInformation
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: VideoUsabilityInformation{},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "AspectRatioInfoPresentFlag is true: AspectRatioIdc is not Extended_SAR",
		Struct: VideoUsabilityInformation{
			AspectRatioInfoPresentFlag: true,
			AspectRatioIdc:             254,
		},
		Binary: mustBitToBytes(
			l,                      // AspectRatioInfoPresentFlag
			l, l, l, l, l, l, l, o, // AspectRatioIdc

			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "AspectRatioInfoPresentFlag is true: AspectRatioIdc is Extended_SAR",
		Struct: VideoUsabilityInformation{
			AspectRatioInfoPresentFlag: true,
			AspectRatioIdc:             255,
			SarWidth:                   0xAAAA,
			SarHeight:                  0x5555,
		},
		Binary: mustBitToBytes(
			l,                      // AspectRatioInfoPresentFlag
			l, l, l, l, l, l, l, l, // AspectRatioIdc
			l, o, l, o, l, o, l, o, l, o, l, o, l, o, l, o, // SarWidth
			o, l, o, l, o, l, o, l, o, l, o, l, o, l, o, l, // SarHeight

			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "OverscanInfoPresentFlag is true: OverscanAppropriateFlag is false",
		Struct: VideoUsabilityInformation{
			OverscanInfoPresentFlag: true,
			OverscanAppropriateFlag: false,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			l, // OverscanInfoPresentFlag
			o, // OverscanAppropriateFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "OverscanInfoPresentFlag is true: OverscanAppropriateFlag is false",
		Struct: VideoUsabilityInformation{
			OverscanInfoPresentFlag: true,
			OverscanAppropriateFlag: true,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			l, // OverscanInfoPresentFlag
			l, // OverscanAppropriateFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "VideoSignalTypePresentFlag is true: VideoFormat is 0, VideoFullRangeFlag is false, ColourDescriptionPresentFlag is false",
		Struct: VideoUsabilityInformation{
			VideoSignalTypePresentFlag:   true,
			VideoFormat:                  0,
			VideoFullRangeFlag:           false,
			ColourDescriptionPresentFlag: false,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag

			l,       // VideoSignalTypePresentFlag
			o, o, o, // VideoFormat
			o, // VideoFullRangeFlag
			o, // ColourDescriptionPresentFlag

			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "VideoSignalTypePresentFlag is true: VideoFormat is 7, VideoFullRangeFlag is true, ColourDescriptionPresentFlag is true: ColourPrimaries, TransferCharacteristics, MatrixCoefficients",
		Struct: VideoUsabilityInformation{
			VideoSignalTypePresentFlag:   true,
			VideoFormat:                  7,
			VideoFullRangeFlag:           true,
			ColourDescriptionPresentFlag: true,
			ColourPrimaries:              0xaa,
			TransferCharacteristics:      0x55,
			MatrixCoefficients:           0xaa,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag

			l,       // VideoSignalTypePresentFlag
			l, l, l, // VideoFormat
			l,                      // VideoFullRangeFlag
			l,                      // ColourDescriptionPresentFlag
			l, o, l, o, l, o, l, o, // ColourPrimaries
			o, l, o, l, o, l, o, l, // TransferCharacteristics
			l, o, l, o, l, o, l, o, // MatrixCoefficients

			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "ChromaLocInfoPresentFlag is true",
		Struct: VideoUsabilityInformation{
			ChromaLocInfoPresentFlag:       true,
			ChromaSampleLocTypeTopField:    1,
			ChromaSampleLocTypeBottomField: 2,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag

			l,       // ChromaLocInfoPresentFlag
			o, l, o, // ChromaSampleLocTypeTopField
			o, l, l, // ChromaSampleLocTypeBottomField

			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "TimingInfoPresentFlag is true",
		Struct: VideoUsabilityInformation{
			TimingInfoPresentFlag: true,
			NumUnitsInTick:        0xaaaaaaaa,
			TimeScale:             0x55555555,
			FixedFrameRateFlag:    true,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag

			l, // TimingInfoPresentFlag

			l, o, l, o, l, o, l, o, l, o, l, o, l, o, l, o,
			l, o, l, o, l, o, l, o, l, o, l, o, l, o, l, o, // NumUnitsInTick

			o, l, o, l, o, l, o, l, o, l, o, l, o, l, o, l,
			o, l, o, l, o, l, o, l, o, l, o, l, o, l, o, l, // TimeScale

			l, // FixedFrameRateFlag

			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "NalHrdParametersPresentFlag is true",
		Struct: VideoUsabilityInformation{
			NalHrdParametersPresentFlag: true,
			HrdNal: &HypotheticalReferenceDecoder{ // minimal HypotheticalReferenceDecoder
				BitRateValueMinus1: make([]uint64, 1),
				CPBSizeValueMinus1: make([]uint64, 1),
				CBRFlag:            make([]bool, 1),
			},
			LowDelayHrdFlag: true,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag

			l, // NalHrdParametersPresentFlag

			/////////////////////////////////
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
			/////////////////////////////////

			o, // VclHrdParametersPresentFlag

			l, // LowDelayHrdFlag

			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "VclHrdParametersPresentFlag is true",
		Struct: VideoUsabilityInformation{
			VclHrdParametersPresentFlag: true,
			HrdVcl: &HypotheticalReferenceDecoder{ // minimal HypotheticalReferenceDecoder
				BitRateValueMinus1: make([]uint64, 1),
				CPBSizeValueMinus1: make([]uint64, 1),
				CBRFlag:            make([]bool, 1),
			},
			LowDelayHrdFlag: false,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag

			l, // VclHrdParametersPresentFlag

			/////////////////////////////////
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
			/////////////////////////////////

			o, // LowDelayHrdFlag

			o, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "PicStructPresentFlag is true",
		Struct: VideoUsabilityInformation{
			PicStructPresentFlag: true,
		},
		Binary: mustBitToBytes(
			o, // AspectRatioInfoPresentFlag
			o, // OverscanInfoPresentFlag
			o, // VideoSignalTypePresentFlag
			o, // ChromaLocInfoPresentFlag
			o, // TimingInfoPresentFlag
			o, // NalHrdParametersPresentFlag
			o, // VclHrdParametersPresentFlag
			l, // PicStructPresentFlag
			o, // BitstreamRestrictionFlag
		),
	},
	{
		Name: "BitstreamRestrictionFlag is true",
		Struct: VideoUsabilityInformation{
			BitstreamRestrictionFlag:           true,
			MotionVectorsOverPicBoundariesFlag: true,
			MaxBytesPerPicDenom:                1,
			MaxBitsPerMbDenom:                  2,
			Log2MaxMvLengthHorizontal:          3,
			Log2MaxMvLengthVertical:            4,
			MaxNumReorderFrames:                5,
			MaxDecFrameBuffering:               6,
		},
		Binary: mustBitToBytes(
			o,       // AspectRatioInfoPresentFlag
			o,       // OverscanInfoPresentFlag
			o,       // VideoSignalTypePresentFlag
			o,       // ChromaLocInfoPresentFlag
			o,       // TimingInfoPresentFlag
			o,       // NalHrdParametersPresentFlag
			o,       // VclHrdParametersPresentFlag
			o,       // PicStructPresentFlag
			l,       // BitstreamRestrictionFlag
			l,       // MotionVectorsOverPicBoundariesFlag
			o, l, o, // MaxBytesPerPicDenom
			o, l, l, // MaxBitsPerMbDenom
			o, o, l, o, o, // Log2MaxMvLengthHorizontal
			o, o, l, o, l, // Log2MaxMvLengthVertical
			o, o, l, l, o, // MaxNumReorderFrames
			o, o, l, l, l, // MaxDecFrameBuffering
		),
	},
}

func TestVideoUsabilityInformation_MarshalBinary(t *testing.T) {
	for _, tt := range VideoUsabilityInformationTestData {
		t.Run(tt.Name, func(t *testing.T) {
			w := newBitWriter()
			err := writeVideoUsabilityInformation(w, tt.Struct)
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, w.Bytes())
		})
	}
}

func TestVideoUsabilityInformation_UnmarshalBinary(t *testing.T) {
	for _, tt := range VideoUsabilityInformationTestData {
		t.Run(tt.Name, func(t *testing.T) {
			r := newBitReader(tt.Binary)
			s, err := readVideoUsabilityInformation(r)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

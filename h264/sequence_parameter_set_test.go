package h264

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustBitToBytes(b ...Bit) []byte {
	if len(b) == 0 {
		return []byte{}
	}
	w := newBitWriterSize((len(b)-1)/8 + 1)
	_, err := w.WriteBit(b...)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

const (
	o = BitZero
	l = BitOne
)

var SequenceParameterSetTestData = []struct {
	Name   string
	Struct SequenceParameterSet
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: SequenceParameterSet{},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag
			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1

			o,       // FrameMbsOnlyFlag
			o,       // MBAdaptiveFrameFieldFlag
			o,       // Direct8x8InterenceFlag
			o,       // FrameCroppingFlag
			o,       // VUIParametersPresentFlag
			l, o, o, // trailing bits
		),
	},
	{
		Name: "ProfileIDC, ConstraintFlags, LevelIDC",
		Struct: SequenceParameterSet{
			ProfileIDC:         170,
			ConstraintSet0Flag: true,
			ConstraintSet1Flag: false,
			ConstraintSet2Flag: true,
			ConstraintSet3Flag: false,
			ConstraintSet4Flag: true,
			ConstraintSet5Flag: true,
			LevelIDC:           170,
		},
		Binary: mustBitToBytes(
			l, o, l, o, l, o, l, o, // ProfileIDC
			l, o, l, o, l, l, o, o, // ConstraintFlags
			l, o, l, o, l, o, l, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag
			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1

			o,       // FrameMbsOnlyFlag
			o,       // MBAdaptiveFrameFieldFlag
			o,       // Direct8x8InterenceFlag
			o,       // FrameCroppingFlag
			o,       // VUIParametersPresentFlag
			l, o, o, // trailing bits
		),
	},
	{
		Name: "SequenceParameterSet is 1",
		Struct: SequenceParameterSet{
			SequenceParamterSetID: 1,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			o, l, o, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "ProfileIDC is 100 (use chroma)",
		Struct: SequenceParameterSet{
			ProfileIDC:      100,
			ChromaFormatIDC: 0,
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // ChromaFormatIDC
			l, // BitDepthLumaMinus8
			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			o, // SequenceScalingMatrixPresentFlag
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType

			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag
			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag

			o,                // FrameCroppingFlag
			o,                // VUIParametersPresentFlag
			l, o, o, o, o, o, // trailing bits
		),
	},
	{
		Name: "ChromaFormatIDC is 1",
		Struct: SequenceParameterSet{
			ProfileIDC:      100,
			ChromaFormatIDC: 1,
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			o, l, o, // ChromaFormatIDC
			l, // BitDepthLumaMinus8
			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			o, // SequenceScalingMatrixPresentFlag

			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag
			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag

			o,          // MBAdaptiveFrameFieldFlag
			o,          // Direct8x8InterenceFlag
			o,          // FrameCroppingFlag
			o,          // VUIParametersPresentFlag
			l, o, o, o, // trailing bits
		),
	},
	{
		Name: "ChromaFormatIDC is 3: SeparateColourPlaneFlag is false",
		Struct: SequenceParameterSet{
			ProfileIDC:              100,
			ChromaFormatIDC:         3,
			SeparateColourPlaneFlag: false,
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,             // SequenceParamterSetID
			o, o, l, o, o, // ChromaFormatIDC
			o, // SeparateColourPlaneFlag
			l, // BitDepthLumaMinus8

			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			o, // SequenceScalingMatrixPresentFlag
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "ChromaFormatIDC is 3: SeparateColourPlaneFlag is true",
		Struct: SequenceParameterSet{
			ProfileIDC:              100,
			ChromaFormatIDC:         3,
			SeparateColourPlaneFlag: true,
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,             // SequenceParamterSetID
			o, o, l, o, o, // ChromaFormatIDC
			l, // SeparateColourPlaneFlag
			l, // BitDepthLumaMinus8

			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			o, // SequenceScalingMatrixPresentFlag
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "BitDepthLumaMinus8, BitDepthLumaMinus8, QPPrimeYZeroTransformBypassFlag",
		Struct: SequenceParameterSet{
			ProfileIDC:                      100,
			BitDepthLumaMinus8:              1,
			BitDepthChromaMinus8:            2,
			QPPrimeYZeroTransformBypassFlag: true,
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			l,       // ChromaFormatIDC
			o, l, o, // BitDepthLumaMinus8
			o, l, l, // BitDepthChromaMinus8

			l, // QPPrimeYZeroTransformBypassFlag
			o, // SequenceScalingMatrixPresentFlag
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag
			l, // PicWidthInMbsMinus1

			l,    // PicHeightInMapUnitsMinus1
			o,    // FrameMbsOnlyFlag
			o,    // MBAdaptiveFrameFieldFlag
			o,    // Direct8x8InterenceFlag
			o,    // FrameCroppingFlag
			o,    // VUIParametersPresentFlag
			l, o, // trailing bits
		),
	},
	{
		Name: "SequenceScalingMatrix: ChromaFormatIDC is not 3: empty matrix",
		Struct: SequenceParameterSet{
			ProfileIDC:                       100,
			SequenceScalingMatrixPresentFlag: true,
			SequenceScalingListPresentFlag: []bool{
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
			},
			ScalingListDeltaScales: [][]int64{
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			},
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // ChromaFormatIDC
			l, // BitDepthLumaMinus8
			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			l, // SequenceScalingMatrixPresentFlag
			// ScalingMatrix
			// ScalingMatrix 0
			o,
			// ScalingMatrix 1
			o,
			// ScalingMatrix 2
			o,
			// ScalingMatrix 3
			o,
			// ScalingMatrix 4
			o,
			// ScalingMatrix 5
			o,
			// ScalingMatrix 6
			o,
			// ScalingMatrix 7
			o,

			l,                // Log2MaxFrameNumMinus4
			l,                // PicOrderCntType
			l,                // Log2MaxPicOrderCntLsbMinus4
			l,                // MaxNumRefFrames
			o,                // GapsInFrameNumValueAllowedFlag
			l,                // PicWidthInMbsMinus1
			l,                // PicHeightInMapUnitsMinus1
			o,                // FrameMbsOnlyFlag
			o,                // MBAdaptiveFrameFieldFlag
			o,                // Direct8x8InterenceFlag
			o,                // FrameCroppingFlag
			o,                // VUIParametersPresentFlag
			l, o, o, o, o, o, // trailing bits
		),
	},
	{
		Name: "SequenceScalingMatrix: ChromaFormatIDC is not 3",
		Struct: SequenceParameterSet{
			ProfileIDC:                       100,
			SequenceScalingMatrixPresentFlag: true,
			SequenceScalingListPresentFlag: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
			},
			ScalingListDeltaScales: [][]int64{
				{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8},
				nil,
				{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8},
				nil,
				{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8},
				nil,
				{
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
				},
				nil,
			},
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // ChromaFormatIDC
			l, // BitDepthLumaMinus8
			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			l, // SequenceScalingMatrixPresentFlag
			// ScalingMatrix
			// ScalingMatrix 0
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 1
			o,
			// ScalingMatrix 2
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 3
			o,
			// ScalingMatrix 4
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 5
			o,
			// ScalingMatrix 6
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 7
			o,

			l,    // Log2MaxFrameNumMinus4
			l,    // PicOrderCntType
			l,    // Log2MaxPicOrderCntLsbMinus4
			l,    // MaxNumRefFrames
			o,    // GapsInFrameNumValueAllowedFlag
			l,    // PicWidthInMbsMinus1
			l,    // PicHeightInMapUnitsMinus1
			o,    // FrameMbsOnlyFlag
			o,    // MBAdaptiveFrameFieldFlag
			o,    // Direct8x8InterenceFlag
			o,    // FrameCroppingFlag
			o,    // VUIParametersPresentFlag
			l, o, // trailing bits
		),
	},
	{
		Name: "SequenceScalingMatrix: ChromaFormatIDC is 3: empty matrix",
		Struct: SequenceParameterSet{
			ProfileIDC:                       100,
			ChromaFormatIDC:                  3,
			SequenceScalingMatrixPresentFlag: true,
			SequenceScalingListPresentFlag: []bool{
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
				false,
			},
			ScalingListDeltaScales: [][]int64{
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			},
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,             // SequenceParamterSetID
			o, o, l, o, o, // ChromaFormatIDC
			o, // SeparateColourPlaneFlag
			l, // BitDepthLumaMinus8
			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			l, // SequenceScalingMatrixPresentFlag
			// ScalingMatrix
			// ScalingMatrix 0
			o,
			// ScalingMatrix 1
			o,
			// ScalingMatrix 2
			o,
			// ScalingMatrix 3
			o,
			// ScalingMatrix 4
			o,
			// ScalingMatrix 5
			o,
			// ScalingMatrix 6
			o,
			// ScalingMatrix 7
			o,
			// ScalingMatrix 8
			o,
			// ScalingMatrix 9
			o,
			// ScalingMatrix 10
			o,
			// ScalingMatrix 11
			o,

			l,             // Log2MaxFrameNumMinus4
			l,             // PicOrderCntType
			l,             // Log2MaxPicOrderCntLsbMinus4
			l,             // MaxNumRefFrames
			o,             // GapsInFrameNumValueAllowedFlag
			l,             // PicWidthInMbsMinus1
			l,             // PicHeightInMapUnitsMinus1
			o,             // FrameMbsOnlyFlag
			o,             // MBAdaptiveFrameFieldFlag
			o,             // Direct8x8InterenceFlag
			o,             // FrameCroppingFlag
			o,             // VUIParametersPresentFlag
			l, o, o, o, o, // trailing bits
		),
	},
	{
		Name: "SequenceScalingMatrix: ChromaFormatIDC is not 3",
		Struct: SequenceParameterSet{
			ProfileIDC:                       100,
			ChromaFormatIDC:                  3,
			SequenceScalingMatrixPresentFlag: true,
			SequenceScalingListPresentFlag: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
			},
			ScalingListDeltaScales: [][]int64{
				{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8},
				nil,
				{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8},
				nil,
				{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8},
				nil,
				{
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
				},
				nil,
				{
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
				},
				nil,
				{
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
					0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5, 6, -6, 7, -7, 8,
				},
				nil,
			},
		},
		Binary: mustBitToBytes(
			o, l, l, o, o, l, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,             // SequenceParamterSetID
			o, o, l, o, o, // ChromaFormatIDC
			o, // SeparateColourPlaneFlag
			l, // BitDepthLumaMinus8
			l, // BitDepthChromaMinus8
			o, // QPPrimeYZeroTransformBypassFlag
			l, // SequenceScalingMatrixPresentFlag
			// ScalingMatrix
			// ScalingMatrix 0
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 1
			o,
			// ScalingMatrix 2
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 3
			o,
			// ScalingMatrix 4
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 5
			o,
			// ScalingMatrix 6
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 7
			o,
			// ScalingMatrix 8
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 9
			o,
			// ScalingMatrix 10
			l,
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			////
			/* 0 */ l /* 1 */, o, l, o /* -1 */, o, l, l /* 2 */, o, o, l, o, o /* -2 */, o, o, l, o, l /* 3 */, o, o, l, l, o /* -3 */, o, o, l, l, l,
			/* 4 */ o, o, o, l, o, o, o /* -4 */, o, o, o, l, o, o, l /* 5 */, o, o, o, l, o, l, o /* -5 */, o, o, o, l, o, l, l /* 6 */, o, o, o, l, l, o, o,
			/* -6 */ o, o, o, l, l, o, l /* 7 */, o, o, o, l, l, l, o /* -7 */, o, o, o, l, l, l, l /* 8 */, o, o, o, o, l, o, o, o, o,
			// ScalingMatrix 11
			o,

			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag
			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "Log2MaxFrameNumMinus4 is 1",
		Struct: SequenceParameterSet{
			Log2MaxFrameNumMinus4: 1,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			o, l, o, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "Log2MaxPicOrderCntLsbMinus4 is 1",
		Struct: SequenceParameterSet{
			PicOrderCntType:             0,
			Log2MaxPicOrderCntLsbMinus4: 1,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			l,       // Log2MaxFrameNumMinus4
			l,       // PicOrderCntType
			o, l, o, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "PicOrderCntType is 1: NumRefFramesInPicOrderCntCycle is empty",
		Struct: SequenceParameterSet{
			PicOrderCntType:                1,
			DeltaPicOrderAlwaysZeroFlag:    false,
			OffsetForNonRefPic:             0,
			OffsetForTopToBottomField:      0,
			NumRefFramesInPicOrderCntCycle: 0,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			l,       // Log2MaxFrameNumMinus4
			o, l, o, // PicOrderCntType
			o, // DeltaPicOrderAlwaysZeroFlag
			l, // OffsetForNonRefPic
			l, // OffsetForTopToBottomField
			l, // NumRefFramesInPicOrderCntCycle

			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,                // PicWidthInMbsMinus1
			l,                // PicHeightInMapUnitsMinus1
			o,                // FrameMbsOnlyFlag
			o,                // MBAdaptiveFrameFieldFlag
			o,                // Direct8x8InterenceFlag
			o,                // FrameCroppingFlag
			o,                // VUIParametersPresentFlag
			l, o, o, o, o, o, // trailing bits
		),
	},
	{
		Name: "PicOrderCntType is 1: DeltaPicOrderAlwaysZeroFlag, OffsetForNonRefPic, OffsetForTopToBottomField",
		Struct: SequenceParameterSet{
			PicOrderCntType:                1,
			DeltaPicOrderAlwaysZeroFlag:    true,
			OffsetForNonRefPic:             -1,
			OffsetForTopToBottomField:      2,
			NumRefFramesInPicOrderCntCycle: 0,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			l,       // Log2MaxFrameNumMinus4
			o, l, o, // PicOrderCntType
			l,       // DeltaPicOrderAlwaysZeroFlag
			o, l, l, // OffsetForNonRefPic
			o, o, l, o, o, // OffsetForTopToBottomField
			l, // NumRefFramesInPicOrderCntCycle

			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,                      // PicWidthInMbsMinus1
			l,                      // PicHeightInMapUnitsMinus1
			o,                      // FrameMbsOnlyFlag
			o,                      // MBAdaptiveFrameFieldFlag
			o,                      // Direct8x8InterenceFlag
			o,                      // FrameCroppingFlag
			o,                      // VUIParametersPresentFlag
			l, o, o, o, o, o, o, o, // trailing bits
		),
	},
	{
		Name: "PicOrderCntType is 1: NumRefFramesInPicOrderCntCycle is 5",
		Struct: SequenceParameterSet{
			PicOrderCntType:                1,
			DeltaPicOrderAlwaysZeroFlag:    false,
			OffsetForNonRefPic:             0,
			OffsetForTopToBottomField:      0,
			NumRefFramesInPicOrderCntCycle: 5,
			OffsetForRefFrame:              []int64{0, 1, -1, 2, -2},
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			l,       // Log2MaxFrameNumMinus4
			o, l, o, // PicOrderCntType
			o,             // DeltaPicOrderAlwaysZeroFlag
			l,             // OffsetForNonRefPic
			l,             // OffsetForTopToBottomField
			o, o, l, l, o, // NumRefFramesInPicOrderCntCycle
			l,       // OffsetForRefFrame[0]
			o, l, o, // OffsetForRefFrame[1]
			o, l, l, // OffsetForRefFrame[2]
			o, o, l, o, o, // OffsetForRefFrame[3]
			o, o, l, o, l, // OffsetForRefFrame[4]

			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l, // PicWidthInMbsMinus1
			l, // PicHeightInMapUnitsMinus1
			o, // FrameMbsOnlyFlag
			o, // MBAdaptiveFrameFieldFlag
			o, // Direct8x8InterenceFlag
			o, // FrameCroppingFlag
			o, // VUIParametersPresentFlag
			l, // trailing bits
		),
	},
	{
		Name: "MaxNumRefFrames, GapsInFrameNumValueAllowedFlag, PicWidthInMbsMinus1, PicHeightInMapUnitsMinus1",
		Struct: SequenceParameterSet{
			MaxNumRefFrames:                1,
			GapsInFrameNumValueAllowedFlag: true,
			PicWidthInMbsMinus1:            2,
			PicHeightInMapUnitsMinus1:      3,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l,       // SequenceParamterSetID
			l,       // Log2MaxFrameNumMinus4
			l,       // PicOrderCntType
			l,       // Log2MaxPicOrderCntLsbMinus4
			o, l, o, // MaxNumRefFrames
			l, // GapsInFrameNumValueAllowedFlag

			o, l, l, // PicWidthInMbsMinus1
			o, o, l, o, o, // PicHeightInMapUnitsMinus1
			o,       // FrameMbsOnlyFlag
			o,       // MBAdaptiveFrameFieldFlag
			o,       // Direct8x8InterenceFlag
			o,       // FrameCroppingFlag
			o,       // VUIParametersPresentFlag
			l, o, o, // trailing bits
		),
	},
	{
		Name: "FrameMbsOnlyFlag is false: MBAdaptiveFrameFieldFlag is true",
		Struct: SequenceParameterSet{
			FrameMbsOnlyFlag:         false,
			MBAdaptiveFrameFieldFlag: true,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,    // PicWidthInMbsMinus1
			l,    // PicHeightInMapUnitsMinus1
			o,    // FrameMbsOnlyFlag
			l,    // MBAdaptiveFrameFieldFlag
			o,    // Direct8x8InterenceFlag
			o,    // FrameCroppingFlag
			o,    // VUIParametersPresentFlag
			l, o, // trailing bits
		),
	},
	{
		Name: "FrameMbsOnlyFlag is true",
		Struct: SequenceParameterSet{
			FrameMbsOnlyFlag: true,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,       // PicWidthInMbsMinus1
			l,       // PicHeightInMapUnitsMinus1
			l,       // FrameMbsOnlyFlag
			o,       // Direct8x8InterenceFlag
			o,       // FrameCroppingFlag
			o,       // VUIParametersPresentFlag
			l, o, o, // trailing bits
		),
	},
	{
		Name: "Direct8x8InterenceFlag is true",
		Struct: SequenceParameterSet{
			Direct8x8InterenceFlag: true,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,    // PicWidthInMbsMinus1
			l,    // PicHeightInMapUnitsMinus1
			o,    // FrameMbsOnlyFlag
			o,    // MBAdaptiveFrameFieldFlag
			l,    // Direct8x8InterenceFlag
			o,    // FrameCroppingFlag
			o,    // VUIParametersPresentFlag
			l, o, // trailing bits
		),
	},
	{
		Name: "FrameCroppingFlag is true: flags are empty",
		Struct: SequenceParameterSet{
			FrameCroppingFlag: true,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,    // PicWidthInMbsMinus1
			l,    // PicHeightInMapUnitsMinus1
			o,    // FrameMbsOnlyFlag
			o,    // MBAdaptiveFrameFieldFlag
			o,    // Direct8x8InterenceFlag
			l,    // FrameCroppingFlag
			l,    // FrameCropLeftOffset
			l,    // FrameCropRightOffset
			l,    // FrameCropTopOffset
			l,    // FrameCropBottomOffset
			o,    // VUIParametersPresentFlag
			l, o, // trailing bits
		),
	},
	{
		Name: "FrameCroppingFlag is true: with flags",
		Struct: SequenceParameterSet{
			FrameCroppingFlag:     true,
			FrameCropLeftOffset:   1,
			FrameCropRightOffset:  2,
			FrameCropTopOffset:    3,
			FrameCropBottomOffset: 4,
		},
		Binary: mustBitToBytes(
			o, o, o, o, o, o, o, o, // ProfileIDC
			o, o, o, o, o, o, o, o, // ConstraintFlags
			o, o, o, o, o, o, o, o, // LevelIDC
			l, // SequenceParamterSetID
			l, // Log2MaxFrameNumMinus4
			l, // PicOrderCntType
			l, // Log2MaxPicOrderCntLsbMinus4
			l, // MaxNumRefFrames
			o, // GapsInFrameNumValueAllowedFlag

			l,       // PicWidthInMbsMinus1
			l,       // PicHeightInMapUnitsMinus1
			o,       // FrameMbsOnlyFlag
			o,       // MBAdaptiveFrameFieldFlag
			o,       // Direct8x8InterenceFlag
			l,       // FrameCroppingFlag
			o, l, o, // FrameCropLeftOffset
			o, l, l, // FrameCropRightOffset
			o, o, l, o, o, // FrameCropTopOffset
			o, o, l, o, l, // FrameCropBottomOffset
			o,    // VUIParametersPresentFlag
			l, o, // trailing bits
		),
	},
	// TODO VUI test
}

func TestSequenceParameterSet_MarshalBinary(t *testing.T) {
	for _, tt := range SequenceParameterSetTestData {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.Struct.MarshalBinary()
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, b)
		})
	}
}

func TestSequenceParameterSet_UnmarshalBinary(t *testing.T) {
	for _, tt := range SequenceParameterSetTestData {
		t.Run(tt.Name, func(t *testing.T) {
			s := SequenceParameterSet{}
			err := s.UnmarshalBinary(tt.Binary)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

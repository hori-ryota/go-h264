package h264

import "io"

type SequenceParameterSet struct {
	ProfileIDC                       uint8
	ConstraintSet0Flag               bool
	ConstraintSet1Flag               bool
	ConstraintSet2Flag               bool
	ConstraintSet3Flag               bool
	ConstraintSet4Flag               bool
	ConstraintSet5Flag               bool
	LevelIDC                         uint8
	SequenceParamterSetID            uint64
	ChromaFormatIDC                  uint64
	SeparateColourPlaneFlag          bool
	BitDepthLumaMinus8               uint64
	BitDepthChromaMinus8             uint64
	QPPrimeYZeroTransformBypassFlag  bool
	SequenceScalingMatrixPresentFlag bool
	SequenceScalingListPresentFlag   []bool
	ScalingListDeltaScales           [][]int64
	Log2MaxFrameNumMinus4            uint64
	PicOrderCntType                  uint64
	Log2MaxPicOrderCntLsbMinus4      uint64
	DeltaPicOrderAlwaysZeroFlag      bool
	OffsetForNonRefPic               int64
	OffsetForTopToBottomField        int64
	NumRefFramesInPicOrderCntCycle   uint64
	OffsetForRefFrame                []int64
	MaxNumRefFrames                  uint64
	GapsInFrameNumValueAllowedFlag   bool
	PicWidthInMbsMinus1              uint64
	PicHeightInMapUnitsMinus1        uint64
	FrameMbsOnlyFlag                 bool
	MBAdaptiveFrameFieldFlag         bool
	Direct8x8InterenceFlag           bool
	FrameCroppingFlag                bool
	FrameCropLeftOffset              uint64
	FrameCropRightOffset             uint64
	FrameCropTopOffset               uint64
	FrameCropBottomOffset            uint64
	VUIParametersPresentFlag         bool
	VUIs                             []VideoUsabilityInformation
}

func (m SequenceParameterSet) MarshalBinary() ([]byte, error) {

	w := newBitWriter()

	if err := w.WriteByte(m.ProfileIDC); err != nil {
		return nil, err
	}
	if _, err := w.WriteBit(
		m.ConstraintSet0Flag,
		m.ConstraintSet1Flag,
		m.ConstraintSet2Flag,
		m.ConstraintSet3Flag,
		m.ConstraintSet4Flag,
		m.ConstraintSet5Flag,
		BitZero,
		BitZero,
	); err != nil {
		return nil, err
	}
	if err := w.WriteByte(m.LevelIDC); err != nil {
		return nil, err
	}

	switch m.ProfileIDC {
	case 100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135:
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.ChromaFormatIDC)); err != nil {
			return nil, err
		}
		if m.ChromaFormatIDC == 3 {
			if _, err := w.WriteBit(m.SeparateColourPlaneFlag); err != nil {
				return nil, err
			}
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.BitDepthLumaMinus8)); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.BitDepthChromaMinus8)); err != nil {
			return nil, err
		}
		if _, err := w.WriteBit(
			m.QPPrimeYZeroTransformBypassFlag,
			m.SequenceScalingMatrixPresentFlag,
		); err != nil {
			return nil, err
		}
		if m.SequenceScalingMatrixPresentFlag {
			xx := 8
			if m.ChromaFormatIDC == 3 {
				xx = 12
			}
			for i := 0; i < xx; i++ {
				if _, err := w.WriteBit(m.SequenceScalingListPresentFlag[i]); err != nil {
					return nil, err
				}
				if m.SequenceScalingListPresentFlag[i] {
					for j := range m.ScalingListDeltaScales[i] {
						if _, err := writeExponentialGolombCoding(w, Int64ToGolombCodeNum(m.ScalingListDeltaScales[i][j])); err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.Log2MaxFrameNumMinus4)); err != nil {
		return nil, err
	}
	if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.PicOrderCntType)); err != nil {
		return nil, err
	}
	switch m.PicOrderCntType {
	case 0:
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.Log2MaxPicOrderCntLsbMinus4)); err != nil {
			return nil, err
		}
	case 1:
		if _, err := w.WriteBit(m.DeltaPicOrderAlwaysZeroFlag); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Int64ToGolombCodeNum(m.OffsetForNonRefPic)); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Int64ToGolombCodeNum(m.OffsetForTopToBottomField)); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.NumRefFramesInPicOrderCntCycle)); err != nil {
			return nil, err
		}
		for i := 0; i < int(m.NumRefFramesInPicOrderCntCycle); i++ {
			if _, err := writeExponentialGolombCoding(w, Int64ToGolombCodeNum(m.OffsetForRefFrame[i])); err != nil {
				return nil, err
			}
		}
	}
	if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.MaxNumRefFrames)); err != nil {
		return nil, err
	}
	if _, err := w.WriteBit(m.GapsInFrameNumValueAllowedFlag); err != nil {
		return nil, err
	}
	if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.PicWidthInMbsMinus1)); err != nil {
		return nil, err
	}
	if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.PicHeightInMapUnitsMinus1)); err != nil {
		return nil, err
	}
	if _, err := w.WriteBit(m.FrameMbsOnlyFlag); err != nil {
		return nil, err
	}
	if !m.FrameMbsOnlyFlag {
		if _, err := w.WriteBit(m.MBAdaptiveFrameFieldFlag); err != nil {
			return nil, err
		}
	}
	if _, err := w.WriteBit(m.Direct8x8InterenceFlag); err != nil {
		return nil, err
	}
	if _, err := w.WriteBit(m.FrameCroppingFlag); err != nil {
		return nil, err
	}
	if m.FrameCroppingFlag {
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.FrameCropLeftOffset)); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.FrameCropRightOffset)); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.FrameCropTopOffset)); err != nil {
			return nil, err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.FrameCropBottomOffset)); err != nil {
			return nil, err
		}
	}
	if _, err := w.WriteBit(m.VUIParametersPresentFlag); err != nil {
		return nil, err
	}
	if m.VUIParametersPresentFlag {
		for i := range m.VUIs {
			if err := writeVideoUsabilityInformation(w, m.VUIs[i]); err != nil {
				return nil, err
			}
		}
	}

	if w.BitLen()%8 != 0 {
		if _, err := w.WriteBit(BitOne); err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}

func (m *SequenceParameterSet) UnmarshalBinary(b []byte) error {
	var err error
	var g uint64
	r := newBitReader(b)

	m.ProfileIDC, err = r.ReadByte()
	if err != nil {
		return err
	}
	m.ConstraintSet0Flag, err = r.ReadBit()
	if err != nil {
		return err
	}
	m.ConstraintSet1Flag, err = r.ReadBit()
	if err != nil {
		return err
	}
	m.ConstraintSet2Flag, err = r.ReadBit()
	if err != nil {
		return err
	}
	m.ConstraintSet3Flag, err = r.ReadBit()
	if err != nil {
		return err
	}
	m.ConstraintSet4Flag, err = r.ReadBit()
	if err != nil {
		return err
	}
	m.ConstraintSet5Flag, err = r.ReadBit()
	if err != nil {
		return err
	}
	if err = r.SeekBit(2); err != nil {
		return err
	}
	m.LevelIDC, err = r.ReadByte()
	if err != nil {
		return err
	}
	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return err
	}
	m.SequenceParamterSetID = GolombCodeNumToUint64(g)
	switch m.ProfileIDC {
	case 100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135:
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.ChromaFormatIDC = GolombCodeNumToUint64(g)
		if m.ChromaFormatIDC == 3 {
			m.SeparateColourPlaneFlag, err = r.ReadBit()
			if err != nil {
				return err
			}
		}
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.BitDepthLumaMinus8 = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.BitDepthChromaMinus8 = GolombCodeNumToUint64(g)
		m.QPPrimeYZeroTransformBypassFlag, err = r.ReadBit()
		if err != nil {
			return err
		}
		m.SequenceScalingMatrixPresentFlag, err = r.ReadBit()
		if err != nil {
			return err
		}
		if m.SequenceScalingMatrixPresentFlag {
			xx := 8
			if m.ChromaFormatIDC == 3 {
				xx = 12
			}
			m.SequenceScalingListPresentFlag = make([]bool, xx)
			m.ScalingListDeltaScales = make([][]int64, xx)
			for i := 0; i < xx; i++ {
				m.SequenceScalingListPresentFlag[i], err = r.ReadBit()
				if err != nil {
					return err
				}
				if m.SequenceScalingListPresentFlag[i] {
					if i < 6 {
						m.ScalingListDeltaScales[i] = make([]int64, 16)
					} else {
						m.ScalingListDeltaScales[i] = make([]int64, 64)
					}
					for j := range m.ScalingListDeltaScales[i] {
						g, err = readExponentialGolombCoding(r)
						if err != nil {
							return err
						}
						m.ScalingListDeltaScales[i][j] = GolombCodeNumToInt64(g)
					}
				}
			}
		}
	}
	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return err
	}
	m.Log2MaxFrameNumMinus4 = GolombCodeNumToUint64(g)
	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return err
	}
	m.PicOrderCntType = GolombCodeNumToUint64(g)

	switch m.PicOrderCntType {
	case 0:
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.Log2MaxPicOrderCntLsbMinus4 = GolombCodeNumToUint64(g)
	case 1:
		m.DeltaPicOrderAlwaysZeroFlag, err = r.ReadBit()
		if err != nil {
			return err
		}
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.OffsetForNonRefPic = GolombCodeNumToInt64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.OffsetForTopToBottomField = GolombCodeNumToInt64(g)
		m.OffsetForNonRefPic = GolombCodeNumToInt64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.NumRefFramesInPicOrderCntCycle = GolombCodeNumToUint64(g)
		m.OffsetForRefFrame = make([]int64, m.NumRefFramesInPicOrderCntCycle)
		for i := range m.OffsetForRefFrame {
			g, err = readExponentialGolombCoding(r)
			if err != nil {
				return err
			}
			m.OffsetForRefFrame[i] = GolombCodeNumToInt64(g)
		}
	}
	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return err
	}
	m.MaxNumRefFrames = GolombCodeNumToUint64(g)
	m.GapsInFrameNumValueAllowedFlag, err = r.ReadBit()
	if err != nil {
		return err
	}
	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return err
	}
	m.PicWidthInMbsMinus1 = GolombCodeNumToUint64(g)
	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return err
	}
	m.PicHeightInMapUnitsMinus1 = GolombCodeNumToUint64(g)
	m.FrameMbsOnlyFlag, err = r.ReadBit()
	if err != nil {
		return err
	}
	if !m.FrameMbsOnlyFlag {
		m.MBAdaptiveFrameFieldFlag, err = r.ReadBit()
		if err != nil {
			return err
		}
	}
	m.Direct8x8InterenceFlag, err = r.ReadBit()
	if err != nil {
		return err
	}
	m.FrameCroppingFlag, err = r.ReadBit()
	if err != nil {
		return err
	}
	if m.FrameCroppingFlag {
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.FrameCropLeftOffset = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.FrameCropRightOffset = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.FrameCropTopOffset = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return err
		}
		m.FrameCropBottomOffset = GolombCodeNumToUint64(g)
	}
	m.VUIParametersPresentFlag, err = r.ReadBit()
	if err != nil {
		return err
	}
	if m.VUIParametersPresentFlag {
		for {
			vui, err := readVideoUsabilityInformation(r)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			m.VUIs = append(m.VUIs, vui)
		}
	}

	return nil
}

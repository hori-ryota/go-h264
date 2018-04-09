package h264

import "encoding/binary"

type VideoUsabilityInformation struct {
	AspectRatioInfoPresentFlag         bool
	AspectRatioIdc                     uint8
	SarWidth                           uint16
	SarHeight                          uint16
	OverscanInfoPresentFlag            bool
	OverscanAppropriateFlag            bool
	VideoSignalTypePresentFlag         bool
	VideoFormat                        uint8
	VideoFullRangeFlag                 bool
	ColourDescriptionPresentFlag       bool
	ColourPrimaries                    uint8
	TransferCharacteristics            uint8
	MatrixCoefficients                 uint8
	ChromaLocInfoPresentFlag           bool
	ChromaSampleLocTypeTopField        uint64
	ChromaSampleLocTypeBottomField     uint64
	TimingInfoPresentFlag              bool
	NumUnitsInTick                     uint32
	TimeScale                          uint32
	FixedFrameRateFlag                 bool
	NalHrdParametersPresentFlag        bool
	HrdNal                             *HypotheticalReferenceDecoder
	VclHrdParametersPresentFlag        bool
	HrdVcl                             *HypotheticalReferenceDecoder
	LowDelayHrdFlag                    bool
	PicStructPresentFlag               bool
	BitstreamRestrictionFlag           bool
	MotionVectorsOverPicBoundariesFlag bool
	MaxBytesPerPicDenom                uint64
	MaxBitsPerMbDenom                  uint64
	Log2MaxMvLengthHorizontal          uint64
	Log2MaxMvLengthVertical            uint64
	MaxNumReorderFrames                uint64
	MaxDecFrameBuffering               uint64
}

const (
	ExtendedSAR = 255
)

func writeVideoUsabilityInformation(w *bitWriter, m VideoUsabilityInformation) (err error) {
	if _, err := w.WriteBit(m.AspectRatioInfoPresentFlag); err != nil {
		return err
	}
	if m.AspectRatioInfoPresentFlag {
		if err := w.WriteByte(m.AspectRatioIdc); err != nil {
			return err
		}
		if m.AspectRatioIdc == ExtendedSAR {
			if _, err := w.Write([]byte{
				byte(m.SarWidth >> 8),
				byte(m.SarWidth),
				byte(m.SarHeight >> 8),
				byte(m.SarHeight),
			}); err != nil {
				return err
			}
		}
	}
	if _, err := w.WriteBit(m.OverscanInfoPresentFlag); err != nil {
		return err
	}
	if m.OverscanInfoPresentFlag {
		if _, err := w.WriteBit(m.OverscanAppropriateFlag); err != nil {
			return err
		}
	}
	if _, err := w.WriteBit(m.VideoSignalTypePresentFlag); err != nil {
		return err
	}
	if m.VideoSignalTypePresentFlag {
		if _, err := w.WriteBit(
			m.VideoFormat&(1<<2) > 0,
			m.VideoFormat&(1<<1) > 0,
			m.VideoFormat&1 > 0,
			m.VideoFullRangeFlag,
			m.ColourDescriptionPresentFlag,
		); err != nil {
			return err
		}
		if m.ColourDescriptionPresentFlag {
			if _, err := w.Write([]byte{
				m.ColourPrimaries,
				m.TransferCharacteristics,
				m.MatrixCoefficients,
			}); err != nil {
				return err
			}
		}
	}
	if _, err := w.WriteBit(m.ChromaLocInfoPresentFlag); err != nil {
		return err
	}
	if m.ChromaLocInfoPresentFlag {
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.ChromaSampleLocTypeTopField)); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.ChromaSampleLocTypeBottomField)); err != nil {
			return err
		}
	}
	if _, err := w.WriteBit(m.TimingInfoPresentFlag); err != nil {
		return err
	}
	if m.TimingInfoPresentFlag {
		if _, err := w.Write([]byte{
			byte(m.NumUnitsInTick >> 24),
			byte(m.NumUnitsInTick >> 16),
			byte(m.NumUnitsInTick >> 8),
			byte(m.NumUnitsInTick),
			byte(m.TimeScale >> 24),
			byte(m.TimeScale >> 16),
			byte(m.TimeScale >> 8),
			byte(m.TimeScale),
		}); err != nil {
			return err
		}
		if _, err := w.WriteBit(m.FixedFrameRateFlag); err != nil {
			return err
		}
	}
	if _, err := w.WriteBit(m.NalHrdParametersPresentFlag); err != nil {
		return err
	}
	if m.NalHrdParametersPresentFlag {
		if err := writeHypotheticalReferenceDecoder(w, *m.HrdNal); err != nil {
			return err
		}
	}
	if _, err := w.WriteBit(m.VclHrdParametersPresentFlag); err != nil {
		return err
	}
	if m.VclHrdParametersPresentFlag {
		if err := writeHypotheticalReferenceDecoder(w, *m.HrdVcl); err != nil {
			return err
		}
	}
	if m.NalHrdParametersPresentFlag || m.VclHrdParametersPresentFlag {
		if _, err := w.WriteBit(m.LowDelayHrdFlag); err != nil {
			return err
		}
	}
	if _, err := w.WriteBit(
		m.PicStructPresentFlag,
		m.BitstreamRestrictionFlag,
	); err != nil {
		return err
	}
	if m.BitstreamRestrictionFlag {
		if _, err := w.WriteBit(m.MotionVectorsOverPicBoundariesFlag); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.MaxBytesPerPicDenom)); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.MaxBitsPerMbDenom)); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.Log2MaxMvLengthHorizontal)); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.Log2MaxMvLengthVertical)); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.MaxNumReorderFrames)); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.MaxDecFrameBuffering)); err != nil {
			return err
		}
	}
	return err
}

func readVideoUsabilityInformation(r *bitReader) (m VideoUsabilityInformation, err error) {
	var g uint64
	var bb []byte

	m.AspectRatioInfoPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.AspectRatioInfoPresentFlag {
		m.AspectRatioIdc, err = r.ReadByte()
		if err != nil {
			return m, err
		}
		if m.AspectRatioIdc == ExtendedSAR {
			bb, err := r.ReadBytes(2)
			if err != nil {
				return m, err
			}
			m.SarWidth = binary.BigEndian.Uint16(bb)
			bb, err = r.ReadBytes(2)
			if err != nil {
				return m, err
			}
			m.SarHeight = binary.BigEndian.Uint16(bb)
		}
	}
	m.OverscanInfoPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.OverscanInfoPresentFlag {
		m.OverscanAppropriateFlag, err = r.ReadBit()
		if err != nil {
			return m, err
		}
	}
	m.VideoSignalTypePresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.VideoSignalTypePresentFlag {
		for i := 0; i < 3; i++ {
			b, err := r.ReadBit()
			if err != nil {
				return m, err
			}
			if b {
				m.VideoFormat |= 1 << (2 - uint8(i))
			}
		}
		m.VideoFullRangeFlag, err = r.ReadBit()
		if err != nil {
			return m, err
		}
		m.ColourDescriptionPresentFlag, err = r.ReadBit()
		if err != nil {
			return m, err
		}
		if m.ColourDescriptionPresentFlag {
			bb, err = r.ReadBytes(3)
			if err != nil {
				return m, err
			}
			m.ColourPrimaries = bb[0]
			m.TransferCharacteristics = bb[1]
			m.MatrixCoefficients = bb[2]
		}
	}
	m.ChromaLocInfoPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.ChromaLocInfoPresentFlag {
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.ChromaSampleLocTypeTopField = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.ChromaSampleLocTypeBottomField = GolombCodeNumToUint64(g)
	}
	m.TimingInfoPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.TimingInfoPresentFlag {
		bb, err := r.ReadBytes(4)
		if err != nil {
			return m, err
		}
		m.NumUnitsInTick = binary.BigEndian.Uint32(bb)
		bb, err = r.ReadBytes(4)
		if err != nil {
			return m, err
		}
		m.TimeScale = binary.BigEndian.Uint32(bb)
		m.FixedFrameRateFlag, err = r.ReadBit()
		if err != nil {
			return m, err
		}
	}
	m.NalHrdParametersPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.NalHrdParametersPresentFlag {
		hdr, err := readHypotheticalReferenceDecoder(r)
		if err != nil {
			return m, err
		}
		m.HrdNal = &hdr
	}
	m.VclHrdParametersPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.VclHrdParametersPresentFlag {
		hdr, err := readHypotheticalReferenceDecoder(r)
		if err != nil {
			return m, err
		}
		m.HrdVcl = &hdr
	}
	if m.NalHrdParametersPresentFlag || m.VclHrdParametersPresentFlag {
		m.LowDelayHrdFlag, err = r.ReadBit()
		if err != nil {
			return m, err
		}
	}
	m.PicStructPresentFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	m.BitstreamRestrictionFlag, err = r.ReadBit()
	if err != nil {
		return m, err
	}
	if m.BitstreamRestrictionFlag {
		m.MotionVectorsOverPicBoundariesFlag, err = r.ReadBit()
		if err != nil {
			return m, err
		}
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.MaxBytesPerPicDenom = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.MaxBitsPerMbDenom = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.Log2MaxMvLengthHorizontal = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.Log2MaxMvLengthVertical = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.MaxNumReorderFrames = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.MaxDecFrameBuffering = GolombCodeNumToUint64(g)
	}
	return m, err
}

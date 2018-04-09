package h264

type HypotheticalReferenceDecoder struct {
	CPBCntMinus1 uint64
	BitRateScale uint8
	CPBSizeScale uint8

	BitRateValueMinus1 []uint64
	CPBSizeValueMinus1 []uint64
	CBRFlag            []bool

	InitialCPBRemovalDelayLengthMinus1 uint8
	CPBRemovalDelayLengthMinus1        uint8
	DPBOutputDelayLengthMinus1         uint8
	TimeOffsetLength                   uint8
}

func writeHypotheticalReferenceDecoder(w *bitWriter, m HypotheticalReferenceDecoder) (err error) {
	if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.CPBCntMinus1)); err != nil {
		return err
	}
	if err := w.WriteByte(m.BitRateScale<<4 | m.CPBSizeScale); err != nil {
		return err
	}
	for i := 0; i <= int(m.CPBCntMinus1); i++ {
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.BitRateValueMinus1[i])); err != nil {
			return err
		}
		if _, err := writeExponentialGolombCoding(w, Uint64ToGolombCodeNum(m.CPBSizeValueMinus1[i])); err != nil {
			return err
		}
		if _, err := w.WriteBit(m.CBRFlag[i]); err != nil {
			return err
		}
	}
	if _, err := w.WriteBit(
		m.InitialCPBRemovalDelayLengthMinus1&(1<<4) > 0,
		m.InitialCPBRemovalDelayLengthMinus1&(1<<3) > 0,
		m.InitialCPBRemovalDelayLengthMinus1&(1<<2) > 0,
		m.InitialCPBRemovalDelayLengthMinus1&(1<<1) > 0,
		m.InitialCPBRemovalDelayLengthMinus1&1 > 0,
		m.CPBRemovalDelayLengthMinus1&(1<<4) > 0,
		m.CPBRemovalDelayLengthMinus1&(1<<3) > 0,
		m.CPBRemovalDelayLengthMinus1&(1<<2) > 0,
		m.CPBRemovalDelayLengthMinus1&(1<<1) > 0,
		m.CPBRemovalDelayLengthMinus1&1 > 0,
		m.DPBOutputDelayLengthMinus1&(1<<4) > 0,
		m.DPBOutputDelayLengthMinus1&(1<<3) > 0,
		m.DPBOutputDelayLengthMinus1&(1<<2) > 0,
		m.DPBOutputDelayLengthMinus1&(1<<1) > 0,
		m.DPBOutputDelayLengthMinus1&1 > 0,
		m.TimeOffsetLength&(1<<4) > 0,
		m.TimeOffsetLength&(1<<3) > 0,
		m.TimeOffsetLength&(1<<2) > 0,
		m.TimeOffsetLength&(1<<1) > 0,
		m.TimeOffsetLength&1 > 0,
	); err != nil {
		return err
	}
	return nil
}

func readHypotheticalReferenceDecoder(r *bitReader) (m HypotheticalReferenceDecoder, err error) {
	var g uint64

	g, err = readExponentialGolombCoding(r)
	if err != nil {
		return m, err
	}
	m.CPBCntMinus1 = GolombCodeNumToUint64(g)

	b, err := r.ReadByte()
	if err != nil {
		return m, err
	}
	m.BitRateScale = b >> 4
	m.CPBSizeScale = b & 0x0f
	m.BitRateValueMinus1 = make([]uint64, m.CPBCntMinus1)
	m.CPBSizeValueMinus1 = make([]uint64, m.CPBCntMinus1)
	m.CBRFlag = make([]bool, m.CPBCntMinus1)
	for i := 0; i <= int(m.CPBCntMinus1); i++ {
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.BitRateValueMinus1[i] = GolombCodeNumToUint64(g)
		g, err = readExponentialGolombCoding(r)
		if err != nil {
			return m, err
		}
		m.CPBSizeValueMinus1[i] = GolombCodeNumToUint64(g)
		m.CBRFlag[i], err = r.ReadBit()
		if err != nil {
			return m, err
		}
	}

	vv := make([]uint8, 4)
	for i := range vv {
		for j := uint8(0); j < 5; j++ {
			b, err := r.ReadBit()
			if err != nil {
				return m, err
			}
			if b {
				vv[i] |= 1 << (4 - j)
			}
		}
	}

	m.InitialCPBRemovalDelayLengthMinus1 = vv[0]
	m.CPBRemovalDelayLengthMinus1 = vv[1]
	m.DPBOutputDelayLengthMinus1 = vv[2]
	m.TimeOffsetLength = vv[3]

	return m, err
}

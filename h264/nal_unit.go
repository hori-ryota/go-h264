package h264

import (
	"bytes"

	"github.com/pkg/errors"
)

type NALUnit struct {
	NALRefIDC      uint8
	NALUnitType    uint8
	SVCExtension   *NALUnitHeaderSVCExtension
	AVC3dExtension *NALUnitHeaderAVC3dExtension
	MVCExtension   *NALUnitHeaderMVCExtension
	RBSPByte       []byte
}

func (m NALUnit) MarshalBinary() ([]byte, error) {

	w := newBitWriter()

	if err := w.WriteByte(
		0x80 | m.NALRefIDC<<5 | m.NALUnitType,
	); err != nil {
		return nil, err
	}

	switch m.NALUnitType {
	case 14, 20, 21:
		switch {
		case m.NALUnitType != 21 && m.SVCExtension != nil:
			bb, err := m.SVCExtension.MarshalBinary()
			if err != nil {
				return nil, err
			}
			if _, err := w.Write(bb); err != nil {
				return nil, err
			}
		case m.NALUnitType == 21 && m.AVC3dExtension != nil:
			bb, err := m.AVC3dExtension.MarshalBinary()
			if err != nil {
				return nil, err
			}
			if _, err := w.Write(bb); err != nil {
				return nil, err
			}
		case m.MVCExtension != nil:
			bb, err := m.MVCExtension.MarshalBinary()
			if err != nil {
				return nil, err
			}
			if _, err := w.Write(bb); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("invalid header")
		}
	}

	rbspb := bytes.Replace(m.RBSPByte, []byte{0, 0, 0x03}, []byte{0, 0, 0x03, 0x03}, -1)
	rbspb = bytes.Replace(rbspb, []byte{0, 0, 0x02}, []byte{0, 0, 0x03, 0x02}, -1)
	rbspb = bytes.Replace(rbspb, []byte{0, 0, 0x01}, []byte{0, 0, 0x03, 0x01}, -1)
	rbspb = bytes.Replace(rbspb, []byte{0, 0, 0x00}, []byte{0, 0, 0x03, 0x00}, -1)

	if _, err := w.Write(rbspb); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (m *NALUnit) UnmarshalBinary(b []byte) error {

	ind := 0

	m.NALRefIDC = (b[ind] >> 5) & 0x03
	m.NALUnitType = b[ind] & 0x1f
	ind++

	switch m.NALUnitType {
	case 14, 20, 21:
		var headerBytes []byte
		var svcExtensionFlag, avc3dExtensionFlag bool
		if m.NALUnitType != 21 {
			headerBytes = b[ind : ind+3]
			svcExtensionFlag = headerBytes[0]>>7 == 1
			ind += 3
		} else {
			avc3dExtensionFlag = b[ind]>>7 == 1
			if avc3dExtensionFlag {
				headerBytes = b[ind : ind+2]
				ind += 2
			} else {
				headerBytes = b[ind : ind+3]
				ind += 3
			}
		}

		switch {
		case svcExtensionFlag:
			m.SVCExtension = &NALUnitHeaderSVCExtension{}
			if err := m.SVCExtension.UnmarshalBinary(headerBytes); err != nil {
				return err
			}
		case avc3dExtensionFlag:
			m.AVC3dExtension = &NALUnitHeaderAVC3dExtension{}
			if err := m.AVC3dExtension.UnmarshalBinary(headerBytes); err != nil {
				return err
			}
		default:
			m.MVCExtension = &NALUnitHeaderMVCExtension{}
			if err := m.MVCExtension.UnmarshalBinary(headerBytes); err != nil {
				return err
			}
		}
	}

	m.RBSPByte = bytes.Join(bytes.Split(b[ind:], []byte{0, 0, 0x03}), []byte{0, 0})

	return nil
}

type NALUnitHeaderSVCExtension struct {
	IDRFlag              bool
	PriorityID           uint8
	NoInterLayerPredFlag bool
	DependencyID         uint8
	QualityID            uint8
	TemporalID           uint8
	UseRefBasePicFlag    bool
	DiscardableFlag      bool
	OutputFlag           bool
}

func (m NALUnitHeaderSVCExtension) MarshalBinary() ([]byte, error) {
	b := make([]byte, 3)
	b[0] |= 1 << 7
	if m.IDRFlag {
		b[0] |= 1 << 6
	}
	b[0] |= m.PriorityID
	if m.NoInterLayerPredFlag {
		b[1] |= 1 << 7
	}
	b[1] |= m.DependencyID << 4
	b[1] |= m.QualityID
	b[2] |= m.TemporalID << 5
	if m.UseRefBasePicFlag {
		b[2] |= 1 << 4
	}
	if m.DiscardableFlag {
		b[2] |= 1 << 3
	}
	if m.OutputFlag {
		b[2] |= 1 << 2
	}
	b[2] |= 3
	return b, nil
}

func (m *NALUnitHeaderSVCExtension) UnmarshalBinary(b []byte) error {
	if len(b) != 3 {
		return errors.Errorf("invalid binary length: len=%d", len(b))
	}

	m.IDRFlag = b[0]>>6&1 == 1
	m.PriorityID = b[0] & 0x3f
	m.NoInterLayerPredFlag = b[1]>>7&1 == 1
	m.DependencyID = b[1] >> 4 & 0x07
	m.QualityID = b[1] & 0x0f
	m.TemporalID = b[2] >> 5
	m.UseRefBasePicFlag = b[2]>>4&1 == 1
	m.DiscardableFlag = b[2]>>3&1 == 1
	m.OutputFlag = b[2]>>2&1 == 1

	return nil
}

type NALUnitHeaderAVC3dExtension struct {
	ViewIdx       uint8
	DepthFlag     bool
	NonIDRFlag    bool
	TemporalID    uint8
	AnchorPicFlag bool
	InterViewFlag bool
}

func (m NALUnitHeaderAVC3dExtension) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	b[0] |= 1 << 7
	b[0] |= m.ViewIdx >> 1
	b[1] |= (m.ViewIdx & 1) << 7
	if m.DepthFlag {
		b[1] |= 1 << 6
	}
	if m.NonIDRFlag {
		b[1] |= 1 << 5
	}
	b[1] |= m.TemporalID << 2
	if m.AnchorPicFlag {
		b[1] |= 1 << 1
	}
	if m.InterViewFlag {
		b[1] |= 1
	}
	return b, nil
}

func (m *NALUnitHeaderAVC3dExtension) UnmarshalBinary(b []byte) error {
	if len(b) != 2 {
		return errors.Errorf("invalid binary length: len=%d", len(b))
	}

	m.ViewIdx = (b[0] & 0x7f << 1) | (b[1] >> 7)
	m.DepthFlag = b[1]>>6&1 == 1
	m.NonIDRFlag = b[1]>>5&1 == 1
	m.TemporalID = b[1] >> 4 & 0x07
	m.AnchorPicFlag = b[1]>>3&1 == 1
	m.InterViewFlag = b[1]>>2&1 == 1

	return nil
}

type NALUnitHeaderMVCExtension struct {
	NonIDRFlag     bool
	PriorityID     uint8
	ViewID         uint16
	TemporalID     uint8
	AnchorPicFlag  bool
	InterViewFlag  bool
	ReservedOneBit bool
}

func (m NALUnitHeaderMVCExtension) MarshalBinary() ([]byte, error) {
	b := make([]byte, 3)
	if m.NonIDRFlag {
		b[0] |= 1 << 6
	}
	b[0] |= m.PriorityID
	b[1] = byte(m.ViewID >> 2)
	b[2] = byte(m.ViewID & 0x03 << 6)
	b[2] |= m.TemporalID << 3
	if m.AnchorPicFlag {
		b[2] |= 1 << 2
	}
	if m.InterViewFlag {
		b[2] |= 1 << 1
	}
	if m.ReservedOneBit {
		b[2] |= 1
	}
	return b, nil
}

func (m *NALUnitHeaderMVCExtension) UnmarshalBinary(b []byte) error {
	if len(b) != 3 {
		return errors.Errorf("invalid binary length: len=%d", len(b))
	}

	m.NonIDRFlag = b[0]>>6&1 == 1
	m.PriorityID = b[0] & 0x3f
	m.ViewID = uint16(b[1])<<2 | uint16(b[2]>>6)
	m.TemporalID = b[2] >> 3 & 0x7
	m.AnchorPicFlag = b[2]>>2&1 == 1
	m.InterViewFlag = b[2]>>1&1 == 1
	m.ReservedOneBit = b[2]&1 == 1

	return nil
}

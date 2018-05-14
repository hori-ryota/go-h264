package h264

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type AVCDecoderConfigurationRecord struct {
	ConfigurationVersion         uint8
	AVCProfileIndication         uint8
	ProfileCompatibility         uint8
	AVCLevelIndication           uint8
	LengthSizeMinusOne           uint8
	SequenceParameterSetNALUnits [][]byte
	PictureParameterSetNALUnits  [][]byte
}

func (m AVCDecoderConfigurationRecord) MarshalBinary() ([]byte, error) {
	l := 7
	for i := range m.SequenceParameterSetNALUnits {
		l += 2 + len(m.SequenceParameterSetNALUnits[i])
	}
	for i := range m.PictureParameterSetNALUnits {
		l += 2 + len(m.PictureParameterSetNALUnits[i])
	}

	b := make([]byte, l)
	b[0] = m.ConfigurationVersion
	b[1] = m.AVCProfileIndication
	b[2] = m.ProfileCompatibility
	b[3] = m.AVCLevelIndication
	b[4] = m.LengthSizeMinusOne | 0xfc

	b[5] = 0xe0 | byte(len(m.SequenceParameterSetNALUnits))
	ind := 6
	for i := range m.SequenceParameterSetNALUnits {
		binary.BigEndian.PutUint16(b[ind:ind+2], uint16(len(m.SequenceParameterSetNALUnits[i])))
		ind += 2
		copy(b[ind:ind+len(m.SequenceParameterSetNALUnits[i])], m.SequenceParameterSetNALUnits[i])
		ind += len(m.SequenceParameterSetNALUnits[i])
	}

	b[ind] = byte(len(m.PictureParameterSetNALUnits))
	ind += 1
	for i := range m.PictureParameterSetNALUnits {
		binary.BigEndian.PutUint16(b[ind:ind+2], uint16(len(m.PictureParameterSetNALUnits[i])))
		ind += 2
		copy(b[ind:ind+len(m.PictureParameterSetNALUnits[i])], m.PictureParameterSetNALUnits[i])
		ind += len(m.PictureParameterSetNALUnits[i])
	}

	return b, nil
}

func (m *AVCDecoderConfigurationRecord) UnmarshalBinary(b []byte) error {
	if len(b) < 7 {
		return errors.Errorf("invalid binary length: len=%d", len(b))
	}

	m.ConfigurationVersion = b[0]
	m.AVCProfileIndication = b[1]
	m.ProfileCompatibility = b[2]
	m.AVCLevelIndication = b[3]
	m.LengthSizeMinusOne = b[4] & 0x03

	numOfSequenceParameterSets := b[5] & 0x1f
	m.SequenceParameterSetNALUnits = make([][]byte, numOfSequenceParameterSets)
	ind := 6
	for i := uint8(0); i < numOfSequenceParameterSets; i++ {
		l := int(binary.BigEndian.Uint16(b[ind : ind+2]))
		ind += 2
		m.SequenceParameterSetNALUnits[i] = make([]byte, l)
		m.SequenceParameterSetNALUnits[i] = b[ind : ind+l]
		ind += l
	}

	numOfPictureParameterSets := b[ind]
	ind += 1
	m.PictureParameterSetNALUnits = make([][]byte, numOfPictureParameterSets)
	for i := uint8(0); i < numOfPictureParameterSets; i++ {
		l := int(binary.BigEndian.Uint16(b[ind : ind+2]))
		ind += 2
		m.PictureParameterSetNALUnits[i] = make([]byte, l)
		m.PictureParameterSetNALUnits[i] = b[ind : ind+l]
		ind += l
	}

	return nil
}

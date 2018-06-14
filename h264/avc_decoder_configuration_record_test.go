package h264

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var AVCDecoderConfigurationRecordTestData = []struct {
	Name   string
	Struct AVCDecoderConfigurationRecord
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: AVCDecoderConfigurationRecord{},
		Binary: []byte{0x00, 0x00, 0x00, 0x00, 0xfc, 0xe0, 0x00},
	},
	{
		Name: "sample",
		Struct: AVCDecoderConfigurationRecord{
			ConfigurationVersion: 1,
			AVCProfileIndication: 2,
			ProfileCompatibility: 3,
			AVCLevelIndication:   4,
			LengthSizeMinusOne:   1,
			SequenceParameterSetNALUnits: [][]byte{
				{0x06, 0x07},
			},
			PictureParameterSetNALUnits: [][]byte{
				{0x08, 0x09, 0x0a},
				{0x0b, 0x0c, 0x0d, 0x0e},
			},
		},
		Binary: []byte{
			0x01, 0x02, 0x03, 0x04, 0xfd, /* 0b11111100 | 1 */
			// SequenceParameterSetNALUnits
			0xe1, /* 0b11100000 | 1 */
			0x00, 0x02, 0x06, 0x07,
			// PictureParameterSetNALUnits
			0x02,
			0x00, 0x03, 0x08, 0x09, 0x0a,
			0x00, 0x04, 0x0b, 0x0c, 0x0d, 0x0e,
		},
	},
}

func TestAVCDecoderConfigurationRecord_MarshalBinary(t *testing.T) {
	for _, tt := range AVCDecoderConfigurationRecordTestData {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.Struct.MarshalBinary()
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, b)
		})
	}
}

func TestAVCDecoderConfigurationRecord_UnmarshalBinary(t *testing.T) {
	for _, tt := range AVCDecoderConfigurationRecordTestData {
		t.Run(tt.Name, func(t *testing.T) {
			s := AVCDecoderConfigurationRecord{}
			err := s.UnmarshalBinary(tt.Binary)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

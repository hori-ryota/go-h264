package h264

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var NALUnitTestData = []struct {
	Name   string
	Struct NALUnit
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: NALUnit{},
		Binary: []byte{0x80 /* 0b10000000 */},
	},
	{
		Name: "unit type = 14 && with SVCExtension",
		Struct: NALUnit{
			NALRefIDC:    1,
			NALUnitType:  14,
			SVCExtension: &NALUnitHeaderSVCExtension{},
			RBSPByte:     []byte{0x01, 0x02},
		},
		Binary: []byte{
			0xae,             /* 0b1 000000 | 0b0 01 00000 | 14 (0b00001110) */
			0x80, 0x00, 0x03, // empty SVCExtension
			0x01, 0x02, // RBSPByte
		},
	},
	{
		Name: "unit type = 14 && with MVCExtension",
		Struct: NALUnit{
			NALRefIDC:    1,
			NALUnitType:  14,
			MVCExtension: &NALUnitHeaderMVCExtension{},
			RBSPByte:     []byte{0x01, 0x02},
		},
		Binary: []byte{
			0xae,             /* 0b1 000000 | 0b0 01 00000 | 14 (0b00001110) */
			0x00, 0x00, 0x00, // empty MVCExtension
			0x01, 0x02, // RBSPByte
		},
	},
	{
		Name: "unit type = 20 && with SVCExtension",
		Struct: NALUnit{
			NALRefIDC:    1,
			NALUnitType:  20,
			SVCExtension: &NALUnitHeaderSVCExtension{},
			RBSPByte:     []byte{0x01, 0x02},
		},
		Binary: []byte{
			0xb4,             /* 0b1 000000 | 0b0 01 00000 | 20 (0b00010100) */
			0x80, 0x00, 0x03, // empty SVCExtension
			0x01, 0x02, // RBSPByte
		},
	},
	{
		Name: "unit type = 20 && with MVCExtension",
		Struct: NALUnit{
			NALRefIDC:    1,
			NALUnitType:  20,
			MVCExtension: &NALUnitHeaderMVCExtension{},
			RBSPByte:     []byte{0x01, 0x02},
		},
		Binary: []byte{
			0xb4,             /* 0b1 000000 | 0b0 01 00000 | 20 (0b00010100) */
			0x00, 0x00, 0x00, // empty MVCExtension
			0x01, 0x02, // RBSPByte
		},
	},
	{
		Name: "unit type = 21 && with AVC3dExtension",
		Struct: NALUnit{
			NALRefIDC:      1,
			NALUnitType:    21,
			AVC3dExtension: &NALUnitHeaderAVC3dExtension{},
			RBSPByte:       []byte{0x01, 0x02},
		},
		Binary: []byte{
			0xb5,       /* 0b1 000000 | 0b0 01 00000 | 21 (0b00010101) */
			0x80, 0x00, // empty AVC3dExtension
			0x01, 0x02, // RBSPByte
		},
	},
	{
		Name: "unit type = 21 && with MVCExtension",
		Struct: NALUnit{
			NALRefIDC:    1,
			NALUnitType:  21,
			MVCExtension: &NALUnitHeaderMVCExtension{},
			RBSPByte:     []byte{0x01, 0x02},
		},
		Binary: []byte{
			0xb5,             /* 0b1 000000 | 0b0 01 00000 | 21 (0b00010101) */
			0x00, 0x00, 0x00, // empty MVCExtension
			0x01, 0x02, // RBSPByte
		},
	},
	{
		Name: "has emulartion three byte: non contiguous",
		Struct: NALUnit{
			RBSPByte: []byte{
				0xff, // dummy
				0x00, 0x00, 0x03,
				0xff, // dummy
				0x00, 0x00, 0x02,
				0xff, // dummy
				0x00, 0x00, 0x01,
				0xff, // dummy
				0x00, 0x00, 0x00,
				0xff, // dummy
			},
		},
		Binary: []byte{
			0x80,
			0xff, // dummy
			0x00, 0x00, 0x03, 0x03,
			0xff, // dummy
			0x00, 0x00, 0x03, 0x02,
			0xff, // dummy
			0x00, 0x00, 0x03, 0x01,
			0xff, // dummy
			0x00, 0x00, 0x03, 0x00,
			0xff, // dummy
		},
	},
	{
		Name: "has emulartion three byte: contiguous",
		Struct: NALUnit{
			RBSPByte: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x03, 0x03, 0x03, 0x03,
				0x00, 0x00, 0x03,
			},
		},
		Binary: []byte{
			0x80,
			0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
			0x00, 0x00, 0x03, 0x00, 0x00, 0x03,
			0x03, 0x03, 0x03, 0x03,
			0x00, 0x00, 0x03, 0x03,
		},
	},
	{
		Name: "has emulartion three byte: example 1",
		Struct: NALUnit{
			RBSPByte: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		Binary: []byte{
			0x80,
			0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
		},
	},
	{
		Name: "has emulartion three byte: example 2",
		Struct: NALUnit{
			RBSPByte: []byte{
				0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00,
			},
		},
		Binary: []byte{
			0x80,
			0x00, 0x00, 0x03, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00,
		},
	},
}

func TestNALUnit_MarshalBinary(t *testing.T) {
	for _, tt := range NALUnitTestData {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.Struct.MarshalBinary()
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, b)
		})
	}
}

func TestNALUnit_UnmarshalBinary(t *testing.T) {
	for _, tt := range NALUnitTestData {
		t.Run(tt.Name, func(t *testing.T) {
			s := NALUnit{}
			err := s.UnmarshalBinary(tt.Binary)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

var NALUnitHeaderSVCExtensionTestData = []struct {
	Name   string
	Struct NALUnitHeaderSVCExtension
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: NALUnitHeaderSVCExtension{},
		Binary: []byte{0x80 /* 0b10000000 */, 0x00, 0x03},
	},
	{
		Name: "sample",
		Struct: NALUnitHeaderSVCExtension{
			IDRFlag:              false,
			PriorityID:           0x2a, // 0b101010
			NoInterLayerPredFlag: true,
			DependencyID:         0x02, // 0b010
			QualityID:            0x0a, // 0b1010
			TemporalID:           0x05, //0b101
			UseRefBasePicFlag:    false,
			DiscardableFlag:      true,
			OutputFlag:           false,
		},
		Binary: []byte{0xaa, 0xaa, 0xab},
	},
}

func TestNALUnitHeaderSVCExtension_MarshalBinary(t *testing.T) {
	for _, tt := range NALUnitHeaderSVCExtensionTestData {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.Struct.MarshalBinary()
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, b)
		})
	}
}

func TestNALUnitHeaderSVCExtension_UnmarshalBinary(t *testing.T) {
	for _, tt := range NALUnitHeaderSVCExtensionTestData {
		t.Run(tt.Name, func(t *testing.T) {
			s := NALUnitHeaderSVCExtension{}
			err := s.UnmarshalBinary(tt.Binary)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

var NALUnitHeaderAVC3dExtensionTestData = []struct {
	Name   string
	Struct NALUnitHeaderAVC3dExtension
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: NALUnitHeaderAVC3dExtension{},
		Binary: []byte{0x80, 0x00},
	},
	{
		Name: "sample",
		Struct: NALUnitHeaderAVC3dExtension{
			ViewIdx:       0x55, // 0b01010101
			DepthFlag:     false,
			NonIDRFlag:    true,
			TemporalID:    0x02, //0b010
			AnchorPicFlag: true,
			InterViewFlag: false,
		},
		Binary: []byte{0xaa, 0xaa},
	},
}

func TestNALUnitHeaderAVC3dExtension_MarshalBinary(t *testing.T) {
	for _, tt := range NALUnitHeaderAVC3dExtensionTestData {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.Struct.MarshalBinary()
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, b)
		})
	}
}

func TestNALUnitHeaderAVC3dExtension_UnmarshalBinary(t *testing.T) {
	for _, tt := range NALUnitHeaderAVC3dExtensionTestData {
		t.Run(tt.Name, func(t *testing.T) {
			s := NALUnitHeaderAVC3dExtension{}
			err := s.UnmarshalBinary(tt.Binary)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

var NALUnitHeaderMVCExtensionTestData = []struct {
	Name   string
	Struct NALUnitHeaderMVCExtension
	Binary []byte
}{
	{
		Name:   "empty struct",
		Struct: NALUnitHeaderMVCExtension{},
		Binary: []byte{0x00, 0x00, 0x00},
	},
	{
		Name: "sample",
		Struct: NALUnitHeaderMVCExtension{
			NonIDRFlag:     true,
			PriorityID:     0x15,   // 0b010101
			ViewID:         0x0155, // 0b0101010101
			TemporalID:     0x02,   //0b010
			AnchorPicFlag:  true,
			InterViewFlag:  false,
			ReservedOneBit: true,
		},
		Binary: []byte{0x55, 0x55, 0x55},
	},
}

func TestNALUnitHeaderMVCExtension_MarshalBinary(t *testing.T) {
	for _, tt := range NALUnitHeaderMVCExtensionTestData {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := tt.Struct.MarshalBinary()
			require.NoError(t, err)
			assert.Equal(t, tt.Binary, b)
		})
	}
}

func TestNALUnitHeaderMVCExtension_UnmarshalBinary(t *testing.T) {
	for _, tt := range NALUnitHeaderMVCExtensionTestData {
		t.Run(tt.Name, func(t *testing.T) {
			s := NALUnitHeaderMVCExtension{}
			err := s.UnmarshalBinary(tt.Binary)
			require.NoError(t, err)
			assert.Equal(t, tt.Struct, s)
		})
	}
}

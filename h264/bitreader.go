package h264

import "io"

type bitReader struct {
	buf []byte
	// bit index
	n int
}

func newBitReader(b []byte) *bitReader {
	return &bitReader{
		buf: b,
	}
}

func (r *bitReader) ReadByte() (byte, error) {
	if len(r.buf)*8-r.n < 8 {
		return 0, io.EOF
	}
	if r.n%8 == 0 {
		b := r.buf[r.n/8]
		r.n += 8
		return b, nil
	}
	b := (r.buf[r.n/8] & (1<<uint8(8-r.n%8) - 1) << uint8(r.n%8)) | (r.buf[r.n/8+1] >> uint8(8-r.n%8))
	r.n += 8
	return b, nil
}

func (r *bitReader) ReadBit() (Bit, error) {
	if len(r.buf)*8-r.n < 1 {
		return BitZero, io.EOF
	}
	b := Bit(r.buf[r.n/8]&(0xff)^(1<<uint8(7-r.n%8)) > 0)
	r.n++
	return b, nil
}

func (r *bitReader) SeekBit(ind int) error {
	if len(r.buf)*8-r.n < ind {
		return io.EOF
	}
	r.n += ind
	return nil
}

func (r *bitReader) ReadBytes(n int) ([]byte, error) {
	bb := make([]byte, n)
	for i := range bb {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		bb[i] = b

	}
	return bb, nil
}

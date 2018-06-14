package h264

type bitWriter struct {
	buf []byte
	// bit index
	n int
}

type Bit = bool

const (
	BitOne  Bit = true
	BitZero Bit = false
)

const defaultBufSize = 32

func newBitWriterSize(size int) *bitWriter {
	if size <= 0 {
		size = defaultBufSize
	}
	return &bitWriter{
		buf: make([]byte, size),
	}
}

func newBitWriter() *bitWriter {
	return newBitWriterSize(defaultBufSize)
}

func (w *bitWriter) Write(p []byte) (writtenByte int, err error) {
	w.mightAppendBuf(len(p) * 8)
	if w.n%8 == 0 {
		copy(w.buf[w.n/8:w.n/8+len(p)], p)
		w.n += len(p) * 8
		return len(p), nil
	}
	for i := range p {
		if err := w.WriteByte(p[i]); err != nil {
			return i * 8, err
		}
	}
	return len(p), nil
}

func (w *bitWriter) mightAppendBuf(bitLen int) {
	byteLen := bitLen/8 + 1
	if byteLen <= len(w.buf)-w.n/8 {
		return
	}
	var b []byte
	if byteLen <= len(w.buf) {
		b = make([]byte, len(w.buf)*2)
	} else {
		b = make([]byte, len(w.buf)+byteLen)
	}
	byteInd := (w.n+1)/8 - 1
	copy(b[:byteInd+1], w.buf[:byteInd+1])
	w.buf = b
}

func (w *bitWriter) WriteByte(p byte) error {
	w.mightAppendBuf(8)
	if w.n%8 == 0 {
		w.buf[w.n/8] = p
		w.n += 8
		return nil
	}
	w.buf[w.n/8] |= p >> uint8(w.n%8)
	w.buf[w.n/8+1] = p & (1<<uint8(w.n%8) - 1) << uint8(8-w.n%8)
	w.n += 8
	return nil
}

func (w *bitWriter) WriteBit(b ...Bit) (writtenBit int, err error) {
	w.mightAppendBuf(len(b))
	for i := range b {
		if b[i] {
			w.buf[w.n/8] |= 1 << uint8(7-w.n%8)
		} else {
			w.buf[w.n/8] &= (0xff) ^ (1 << uint8(7-w.n%8))
		}
		w.n++
	}
	return len(b), nil
}

func (w *bitWriter) Bytes() []byte {
	return w.buf[:(w.n+7)/8]
}

func (w *bitWriter) BitLen() int {
	return w.n
}

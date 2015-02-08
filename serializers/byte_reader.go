package serializers

type ByteReader struct {
	data   []byte
	pos    int
	length int
}

func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{
		data: data,
		pos:  0,
	}
}

func (b *ByteReader) UInt8(dst *uint8) {
	b.checkAvailLength(1)
	*dst = b.data[b.pos]
	b.pos++
}

func (b *ByteReader) UInt16(dst *uint16) {
	b.checkAvailLength(2)
	*dst = uint16(b.data[b.pos]) + uint16(b.data[b.pos+1])*256
	b.pos += 2
}

func (b *ByteReader) UInt32(dst *uint32) {
	b.checkAvailLength(4)
	*dst = uint32(b.data[b.pos]) + uint32(b.data[b.pos+1])*256 + uint32(b.data[b.pos+2])*256*256 + uint32(b.data[b.pos+3])*256*256*256
	b.pos += 4
}

func (b *ByteReader) CString(dst *string) {
	out := []byte{}
	for {
		var char uint8
		b.UInt8(&char)
		if char == 0 {
			break
		}
		out = append(out, char)
	}
	*dst = string(out)
	b.pos += len(out) + 1
}

func (b *ByteReader) Pos() int {
	return b.pos
}

func (b *ByteReader) IsReader() bool {
	return true
}

func (b *ByteReader) checkAvailLength(requiredLength int) {
	if len(b.data)-b.pos < requiredLength {
		panic("Need more bytes")
	}
}

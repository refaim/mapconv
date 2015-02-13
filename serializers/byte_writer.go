package serializers

import (
	"unsafe"
)

type ByteWriter struct {
	data []byte
}

func NewByteWriter() *ByteWriter {
	return &ByteWriter{data: nil}
}

func (b *ByteWriter) UInt8(dst *uint8) {
	b.append(*dst)
}

func (b *ByteWriter) UInt16(dst *uint16) {
	v := *dst
	b.append(uint8(v), uint8(v>>8))
}

func (b *ByteWriter) UInt32(dst *uint32) {
	v := *dst
	b.append(uint8(v), uint8(v>>8), uint8(v>>16), uint8(v>>24))
}

func (b *ByteWriter) Int8(dst *int8) {
	b.UInt8((*uint8)(unsafe.Pointer(dst)))
}

func (b *ByteWriter) Int16(dst *int16) {
	b.UInt16((*uint16)(unsafe.Pointer(dst)))
}

func (b *ByteWriter) Int32(dst *int32) {
	b.UInt32((*uint32)(unsafe.Pointer(dst)))
}

func (b *ByteWriter) CString(dst *string) {
	b.append([]uint8(*dst)...)
	b.append(0)
}

func (b *ByteWriter) Data() []byte {
	return b.data
}

func (b *ByteWriter) IsReader() bool {
	return false
}

func (b *ByteWriter) append(values ...uint8) {
	b.data = append(b.data, values...)
}

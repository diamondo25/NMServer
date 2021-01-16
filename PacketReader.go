package main

import (
	"bytes"
)

type PacketReader struct {
	opcode uint16
	reader *bytes.Reader
}

func NewPacketReader(opcode uint16, buffer []byte) *PacketReader {
	return &PacketReader{
		opcode,
		bytes.NewReader(buffer)}
}

func (m *PacketReader) ReadByte() byte {
	result, _ := m.reader.ReadByte()
	return result
}

func (m *PacketReader) ReadInt16() int16 {
	return int16(m.ReadByte() | m.ReadByte()<<8)
}

func (m *PacketReader) ReadInt32() int32 {
	return int32(m.ReadByte() | m.ReadByte()<<8 | m.ReadByte()<<16 | m.ReadByte()<<24)
}

func (m *PacketReader) ReadInt64() int64 {
	return int64(m.ReadUInt32() | m.ReadUInt32()<<32)
}

func (m *PacketReader) ReadUInt16() uint16 {
	return uint16(m.ReadByte() | m.ReadByte()<<8)
}

func (m *PacketReader) ReadUInt32() uint32 {
	return uint32(m.ReadByte() | m.ReadByte()<<8 | m.ReadByte()<<16 | m.ReadByte()<<24)
}

func (m *PacketReader) ReadUInt64() uint64 {
	return uint64(m.ReadUInt32() | m.ReadUInt32()<<32)
}

func (m *PacketReader) ReadBytes(length int) []byte {
	buffer := make([]byte, length)
	m.reader.Read(buffer)
	return buffer
}

func (m *PacketReader) ReadString(length int, ascii bool) string {
	if !ascii {
		length *= 2
	}

	return string(m.ReadBytes(length))
}

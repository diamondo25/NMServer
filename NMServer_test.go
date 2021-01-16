package main

import (
	"bytes"
	"testing"
)

func TestPacketReader(t *testing.T) {
	input := []byte{0x00, 0x24, 0x00, 0x11, 0x18, 0x00, 0x00, 0x20, 0x02, 0x00, 0x00, 0x14, 0x5F, 0xD6, 0x11, 0x91, 0x86, 0x4B, 0xB4, 0x35, 0xEE, 0x04, 0x2A, 0x1F, 0x89, 0x7D, 0xED, 0x4E, 0x88, 0xDB, 0x5F, 0xDC, 0x49, 0x7D, 0xCF, 0xBD, 0x18, 0x05, 0x46, 0x2B}
	dummyCon := NewDummyCon(input)

	packet, err := ReadPacket(dummyCon)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(packet)
}

func TestPacketWriter(t *testing.T) {
	dummyCon := NewDummyCon([]byte{})

	op := NewPacketWriter(0x00)
	op.WriteString("This is a unicode test", false)
	op.WriteString("This is an ascii test", true)
	for i := 1; i < 255; i++ {
		op.WriteByte(byte(i))
	}

	SendPacket(dummyCon, op)

	PrintBytes(dummyCon.writeBuffer.Bytes())

	input := dummyCon.writeBuffer.Bytes()
	dummyCon = NewDummyCon(input)

	packet, err := ReadPacket(dummyCon)
	if err != nil {
		t.Fatal(err)
	}
	PrintBytes(dummyCon.data)

	if packet.opcode != 0x0000 {
		t.Fatal("Not the same opcode read")
	}

	str := packet.ReadString(int(packet.ReadInt16()), false)
	if str != "This is a unicode test" {
		t.Fatal("Readstring 1 failed: ", str)
	}

	str = packet.ReadString(int(packet.ReadInt16()), false)
	if str != "This is an ascii test" {
		t.Fatal("Readstring 2 failed: ", str)
	}

	for i := 1; i < 255; i++ {
		b := packet.ReadByte()
		if b != byte(i) {
			t.Fatal("Not the same byte read; ", b, " != ", i)
		}
	}
	t.Log(packet)
}

func TestPacketBIGpacket(t *testing.T) {

	dummyCon := NewDummyCon([]byte{})
	hugeSize := 0x00123456

	writeSmallSize(hugeSize, dummyCon.writeBuffer)

	PrintBytes(dummyCon.writeBuffer.Bytes())

	reader := bytes.NewReader(dummyCon.writeBuffer.Bytes())

	readSize := readSmallSize(reader)
	if readSize != int32(hugeSize) {
		t.Fatalf("Sizes didnt match: %08X != %08X", readSize, hugeSize)
	}
}

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func readSmallSize(reader *bytes.Reader) int32 {
	a, _ := reader.ReadByte()
	b, _ := reader.ReadByte()
	c, _ := reader.ReadByte()
	return int32(int32(a)<<16 | int32(b)<<8 | int32(c)<<0)
}

func writeSmallSize(input int, writer *bytes.Buffer) {
	writer.WriteByte(byte(input >> 16))
	writer.WriteByte(byte(input >> 8))
	writer.WriteByte(byte(input >> 0))
}

func ReadPacket(conn net.Conn) (*PacketReader, error) {
	// fmt.Println("---------------------------------------------------------------------")
	var buffer []byte
	var reader *bytes.Reader

	var _uint32 uint32
	var _int16 int16

	// HEADER
	buffer, _ = ReadBytes(conn, 4)

	reader = bytes.NewReader(buffer)
	//fmt.Println("Header received")
	//PrintBytes(buffer)

	if err := binary.Read(reader, binary.BigEndian, &_int16); err != nil {
		fmt.Println("Error while reading 2 bytes", err)
		return nil, err
	}

	packetSize := int(_int16)

	var opcode int16

	if err := binary.Read(reader, binary.BigEndian, &opcode); err != nil {
		fmt.Println("Error while reading 2 bytes", err)
		return nil, err
	}

	// BODY
	buffer, _ = ReadBytes(conn, packetSize)
	reader = bytes.NewReader(buffer)
	//fmt.Println("Body received")
	//PrintBytes(buffer)

	checkCode, _ := reader.ReadByte() // 0 - 1

	if checkCode != 0x18 {
		fmt.Println("This thing is not 0x18", checkCode)
		return nil, errors.New("0x18 wasnt 0x18")
	}

	_length := readSmallSize(reader) // 1 - 4
	//fmt.Println("_length", _length)

	flags, _ := reader.ReadByte()         // 4 - 5
	originalSize := readSmallSize(reader) // 5 - 8
	//fmt.Println("flags", flags)
	//fmt.Println("originalSize", originalSize)

	if originalSize != _length-12 {
		return nil, errors.New("Original size should be 12 lower than size")
	}

	binary.Read(reader, binary.BigEndian, &_uint32) // 8 - 12
	xorValue := _uint32

	binary.Read(reader, binary.BigEndian, &_uint32) // 12 - 16
	/*
		fmt.Println("Data received")
		PrintBytes(buffer[:16])
	*/
	decryptedBuffer := buffer[16:]
	if (flags & 0x02) != 0 {
		decryptedBuffer = Decrypt(decryptedBuffer, xorValue)
	}

	fmt.Println("RAW")
	PrintBytes(decryptedBuffer)

	fmt.Println("STRING")
	fmt.Println(string(decryptedBuffer))

	packet := NewPacketReader(uint16(opcode), decryptedBuffer)

	return packet, nil
}

func SendPacket(conn net.Conn, packet *PacketWriter) {
	fmt.Println("---------------------------------------------------------------------")

	useXor := false
	useZlib := false // not implemented

	var flags byte = 0
	if useXor {
		flags += 0x02
	}
	if useZlib {
		flags += 0x04
	}

	size := 1 + 3 + 1 + 3
	size += 4
	size += 4 // I think its a checksum
	size += packet.writer.Len()

	buffer := new(bytes.Buffer)
	// header
	binary.Write(buffer, binary.BigEndian, int16(size))
	binary.Write(buffer, binary.BigEndian, packet.opcode)

	// body
	binary.Write(buffer, binary.BigEndian, byte(0x18))
	writeSmallSize(size, buffer)
	binary.Write(buffer, binary.BigEndian, flags)
	writeSmallSize(size-12, buffer)

	var xorKey uint32 = 0xDEADB00B
	binary.Write(buffer, binary.BigEndian, xorKey)

	binary.Write(buffer, binary.BigEndian, uint32(0xAABBCCDD))

	dataBuffer := packet.writer.Bytes()
	
	fmt.Println("buffer sending")
	PrintBytes(dataBuffer)
	fmt.Println(string(dataBuffer))

	if useXor {
		dataBuffer = Encrypt(dataBuffer, xorKey)
	}

	buffer.Write(dataBuffer)

	conn.Write(buffer.Bytes())
}

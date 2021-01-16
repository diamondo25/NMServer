package main

import ()

var xor_table []uint32 = []uint32{
	0x040FC1578, 0x0113B6C1F, 0x08389CA19, 0x0E2196CD8,
	0x074901489, 0x04AAB1566, 0x07B8C12A0, 0x00018FFCD,
	0x0CCAB704B, 0x07B5A8C0F, 0x0AA13B891, 0x0DE419807,
	0x012FFBCAE, 0x05F5FBA34, 0x010F5AC99, 0x0B1C1DD01}

func GetUInt32(buffer []byte, offset int) uint32 {
	a := int(buffer[offset+0] << 0)
	b := int(buffer[offset+1]) << 8
	c := int(buffer[offset+2]) << 16
	d := int(buffer[offset+3]) << 24
	return uint32(a | b | c | d)
}

func GetBytes(input uint32) (byte, byte, byte, byte) {
	return byte((input >> 0) & 0xFF),
		byte((input >> 8) & 0xFF),
		byte((input >> 16) & 0xFF),
		byte((input >> 24) & 0xFF)
}

func SetBytes(input uint32, buffer []byte, offset int) {
	a, b, c, d := GetBytes(input)
	buffer[offset+0] = a
	buffer[offset+1] = b
	buffer[offset+2] = c
	buffer[offset+3] = d
}

func Encrypt(buffer []byte, seed uint32) []byte {
	var temp uint32 = 0
	var temp2 uint32 = 0

	output := make([]byte, len(buffer))

	for i := 0; i < len(buffer)>>2; i++ {
		temp = temp2 ^ xor_table[i&15] ^ seed
		temp2 = GetUInt32(buffer, i*4)
		SetBytes(temp^temp2, output, i*4)
	}

	return output
}

func Decrypt(buffer []byte, seed uint32) []byte {
	var temp uint32 = 0
	var temp2 uint32 = 0

	output := make([]byte, len(buffer))

	for i := 0; i < len(buffer)/4; i++ {
		temp2 = GetUInt32(buffer, i*4)
		temp2 ^= (temp ^ xor_table[i&15] ^ seed)
		SetBytes(temp2, output, i*4)
		temp = temp2
	}

	return output
}

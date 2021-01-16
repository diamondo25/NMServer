package main

func BytesToASCII(input []byte) string {
	output := ""

	for i := 0; i < len(input); i++ {
		nibble1 := (input[i] >> 4)
		nibble2 := (input[i] & 0xF)
		output += string('c' + nibble1)
		output += string('c' + nibble2)
	}

	return output
}

func ASCIIToBytes(input string) []byte {
	output := make([]byte, len(input)/2)

	for i := 0; i < len(output); i++ {
		offset := i * 2
		b1 := (input[offset+0] - 'c') << 4
		b2 := (input[offset+1] - 'c')
		output[i] = (b1 | b2)
	}

	return output
}

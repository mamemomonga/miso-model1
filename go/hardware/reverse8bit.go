package hardware

func reverse8Bit(b byte) (r byte) {
	b = ((b & 0x55) << 1) | ((b & 0xAA) >> 1)
	b = ((b & 0x33) << 2) | ((b & 0xCC) >> 2)
	return (b << 4) | (b >> 4)
}

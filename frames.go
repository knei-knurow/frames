// Package frames provides useful functions to deal with data frames.
package frames

import "fmt"

// Frame represents a data frame that can be e.g sent by USART.
//
// Frame starts with a header that is always 2 bytes.
// Header can only contain uppercase ASCII letters.
// Directly afer a header comes length byte which describes how long is data.
// After the length byte comes a plus sign ("+").
// Then comes an arbitrary-length data.
// Data is terminated with a hash sign ("#").
// The last byte is a simple 8-bit CRC checksum.
//
// Some example frames (H = header byte, D = data byte, C = CRC byte):
//
// HH4+DDDD#C
//
// LD5+DDDDD#C
//
// XD7+DDDDDDD#C
type Frame []byte

// Header returns frame's header. It is always 2 bytes.
func (f Frame) Header() []byte {
	return f[:2]
}

// LenData returns the length of frame's data in bytes.
func (f Frame) LenData() int {
	return int(f[2])
}

// Data returns frame's data part from the first byte after a plus sign ("+") up
// to the antepenultimate (last but one - 1) byte.
func (f Frame) Data() []byte {
	headerLength := len(f.Header())
	begin := headerLength + 2 // example: LD4+DDDD : we want to start from D (so index 4)
	end := begin + f.LenData()

	return f[begin:end]
}

// Checksum returns frame's last byte - a simple CRC checksum.
func (f Frame) Checksum() byte {
	return f[len(f)-1]
}

// Create creates a new frame.
// The frame starts with header and contains data.
// Create also calculates the checksum using CalculateChecksum.
// Data length must not overflow byte.
func Create(header [2]byte, data []byte) (frame Frame) {
	frame = make(Frame, len(header)+1+1+len(data)+2)
	copy(frame[:2], header[:])
	frame[len(header)] = byte(len(data))
	frame[len(header)+1] = '+'
	copy(frame[len(header)+2:len(frame)-2], data)
	frame[len(frame)-2] = '#'
	frame[len(frame)-1] = CalculateChecksum(frame)

	return
}

// Recreate creates a new frame from already available byte buffer.
// It does not check whether buf represents a correct frame.
// To check if the newly created frame is correct, use Verify function.
// Data length must not overflow byte.
func Recreate(buf []byte) (frame Frame) {
	frame = make(Frame, len(buf))
	copy(frame[:], buf[:])
	return
}

// Assemble creates a frame from already available values.
//
// Deprecated: you should probably just use Recreate.
func Assemble(header [2]byte, length byte, data []byte, checksum byte) (frame Frame) {
	frame = make(Frame, len(header)+1+1+len(data)+2)

	copy(frame[:2], header[:])
	frame[len(header)] = length
	frame[len(header)+1] = '+'
	copy(frame[len(header)+2:len(frame)-2], data)
	frame[len(frame)-2] = '#'
	frame[len(frame)-1] = checksum

	return
}

// Verify checks whether the frame is valid (i.e of correct format).
//
// The frame must have:
//
// - at 0th and 1st index: a header consisting of uppercase ASCII header or
// numbers
//
// - at 2nd index: "length byte" that is equal to the length of data
//
// - at 3rd index: a plus sign ("+")
//
// - at penultimate position: a hash sign ("#")
//
// - at last position: a checksum must be correct
func Verify(frame Frame) bool {
	first := frame[0]
	valid1 := (first > 'A' && first < 'Z') || (first > '0' && first < '9')
	if !valid1 {
		return false
	}

	second := frame[1]
	valid2 := (second > 'A' && second < 'Z') || (second > '0' && second < '9')
	if !valid2 {
		return false
	}

	if frame[2] != byte(frame.LenData()) {
		return false
	}

	if frame[3] != '+' {
		return false
	}

	if frame[len(frame)-2] != '#' {
		return false
	}

	checksum := CalculateChecksum(frame)
	return checksum == frame.Checksum()
}

// CalculateChecksum calculates the simple CRC checksum of frame.
//
// It takes all frame's bytes into account, except the last byte, because
// the last byte is the checksum itself.
func CalculateChecksum(frame Frame) (crc byte) {
	crc = frame[0]
	for i := 1; i < len(frame)-1; i++ {
		crc ^= frame[i]
	}

	return
}

func (f Frame) String() string {
	return fmt.Sprintf("%s+%x#%x", f.Header(), f.Data(), f.Checksum())
}

// DescribeByte prints everything most common representations of a byte.
// It prints b's binary value, decimal, hexadecimal value and ASCII.
func DescribeByte(b byte) string {
	return fmt.Sprintf("byte(bin: %08b, dec: %3d, hex: %02x, ASCII: %+q)", b, b, b, b)
}

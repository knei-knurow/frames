package frames_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/knei-knurow/frames"
)

var testCases = []struct {
	inputHeader      [2]byte
	inputData        []byte
	expectedChecksum byte
	expectedLength   int
	frame            []byte
}{
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{},
		expectedChecksum: 0x00,
		frame:            []byte{'L', 'D', '0', '+', '#', 0x00},
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'A'},
		expectedChecksum: 0x40,
		frame:            []byte{'L', 'D', '1', '+', 'A', '#', 0x40},
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'t', 'e', 's', 't'},
		expectedChecksum: 0x12,
		frame:            []byte{'L', 'D', '4', '+', 't', 'e', 's', 't', '#', 0x12},
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'d', 'u', 'p', 'c', 'i', 'a'},
		expectedChecksum: 0x0c,
		frame:            []byte{'L', 'D', '4', '+', 'd', 'u', 'p', 'c', 'i', 'a', '#', 0x0c},
	},
	{
		inputHeader:      [2]byte{'L', 'D'},
		inputData:        []byte{'l', 'o', 'l', 'x', 'd'},
		expectedChecksum: 0x76,
		frame:            []byte{'L', 'D', '5', '+', 'l', 'o', 'l', 'x', 'd', '#', 0x76},
	},
	{
		inputHeader:      [2]byte{'M', 'T'},
		inputData:        []byte{'d', 'o', 'n', 'd', 'u'},
		expectedChecksum: 0x60,
		frame:            []byte{'M', 'T', '5', '+', 'd', 'o', 'n', 'd', 'u', '#', 0x60},
	},
}

func TestCreate(t *testing.T) {
	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Create(tc.inputHeader, tc.inputData)

			if !bytes.Equal(gotFrame.Header(), tc.inputHeader[:]) {
				t.Errorf("got header % x, want header % x", gotFrame.Header(), tc.inputHeader)
			}

			if !bytes.Equal(gotFrame.Data(), tc.inputData) {
				t.Errorf("got data % x, want data % x", gotFrame.Data(), tc.inputData)
			}

			if gotFrame.LenData() != len(tc.inputData) {
				t.Errorf("got data length %d, want data %d", gotFrame.LenData(), len(tc.inputData))
			}

			if gotFrame.Checksum() != tc.expectedChecksum {
				t.Errorf("got checksum % x, want checksum % x", gotFrame.Checksum(), tc.expectedChecksum)
			}

			if !bytes.Equal(gotFrame, tc.frame) {
				t.Errorf("frame create not correct")
			}
		})
	}
}

func TestAssemble(t *testing.T) {
	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Assemble(tc.inputHeader, byte(len(tc.inputData)), tc.inputData, tc.expectedChecksum)

			if !bytes.Equal(gotFrame.Header(), tc.inputHeader[:]) {
				t.Errorf("got header % x, want header % x", gotFrame.Header(), tc.inputHeader)
			}

			if !bytes.Equal(gotFrame.Data(), tc.inputData) {
				t.Errorf("got data % x, want data % x", gotFrame.Data(), tc.inputData)
			}

			if gotFrame.Checksum() != tc.expectedChecksum {
				t.Errorf("got checksum % x, want checksum % x", gotFrame.Checksum(), tc.expectedChecksum)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	verifyTestCases := []struct {
		frame []byte
		valid bool
	}{
		{
			frame: []byte{'L', 'D', '0', '+', '#', 0x00},
			valid: true,
		},
		{
			frame: []byte{'L', 'D', '1', '+', 'A', '#', 0x40},
			valid: true,
		},
		{
			frame: []byte{'L', 'D', '4', '+', 't', 'e', 's', 't', '#', 0x12},
			valid: true,
		},
		{
			frame: []byte{'L', 'D', '4', '+', 'd', 'u', 'p', 'c', 'i', 'a', '#', 0x0c},
			valid: true,
		},
		{
			frame: []byte{'L', 'D', '5', '+', 'l', 'o', 'l', 'x', 'd', '#', 0x76},
			valid: true,
		},
		{
			frame: []byte{'M', 'T', '5', '+', 'd', 'o', 'n', 'd', 'u', '#', 0x60},
			valid: true,
		},
	}

	for i, tc := range verifyTestCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			testFrame := frames.Create(testCases[i].inputHeader, testCases[i].inputData)
			log.Printf("now goes 'create' frame: %x, len: %d", testFrame, len(testFrame))
			log.Printf("now goes 'verify' frame: %x, len: %d", tc.frame, len(tc.frame))
			log.Print("tc: " + frames.DescribeByte(frames.CalculateChecksum(testFrame)))
			log.Print("vtc: " + frames.DescribeByte(frames.CalculateChecksum(tc.frame)))
			if frames.Verify(tc.frame) != tc.valid {
				t.Errorf("frame verification failed")
			}
		})
	}
}

func TestRecreate(t *testing.T) {
	for i, tc := range testCases {
		testName := fmt.Sprintf("test %d", i)
		t.Run(testName, func(t *testing.T) {
			gotFrame := frames.Create(tc.inputHeader, tc.inputData)
			recreatedFrame := frames.Recreate(gotFrame[:])

			if !bytes.Equal(gotFrame, recreatedFrame) {
				t.Error("frame recreation failed")
			}
		})
	}
}

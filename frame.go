package xbeeapi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const frameStartDelimiter = 0x7e
const minFrameSize = 4

// Frame is the structured data packet used in XBee API mode.
//
//    1 Byte          2 Bytes        Length bytes (Type is 1 Byte)    1 Byte
// +----------+--------------------+--------------------------------+---------+
// |  0x7e    |      Length        |        Frame Data              |Checksum |
// |          |                    |    Type  |       Data          |         |
// +----------+--------------------+--------------------------------+---------+
type Frame struct {
	Length    uint16
	FrameData *RawFrameData
	Checksum  byte
}

// FrameParseError describes error from parsing frames
type FrameParseError struct {
	msg string
}

func (e *FrameParseError) Error() string { return e.msg }

func checksumVerifyFrame(serializedFrame []byte) error {
	var cs byte

	if len(serializedFrame) < minFrameSize {
		return &FrameParseError{msg: "Frame too short for checksum"}
	}
	for _, b := range serializedFrame[3:] {
		cs += b
	}

	if cs != 0xff {
		return &FrameParseError{msg: "Checksum invalid"}
	}
	return nil
}

func startDelimiterValid(serializedFrame []byte) error {
	if len(serializedFrame) == 0 {
		return &FrameParseError{msg: "Start delimiter field not found: Frame length not long enough"}
	}
	if serializedFrame[0] != 0x7e {
		return &FrameParseError{msg: fmt.Sprintf("Invalid start delimiter: %x", 0x7e)}
	}
	return nil
}

func lengthField(serializedFrame []byte) (uint16, error) {
	if len(serializedFrame) < 3 {
		return 0, &FrameParseError{msg: "Length field not found: Frame length not long enough"}
	}
	return uint16(int(serializedFrame[1])<<8) + uint16(serializedFrame[2]), nil
}

func totalFrameLength(dataLen uint16) int {
	return 3 + int(dataLen) + 1
}

func frameDataField(serializedFrame []byte) []byte {
	return serializedFrame[3:(len(serializedFrame) - 1)]
}

// NewFrameFromData creates and initializes a new frame from
// frame type and data.
func NewFrame(frameData FrameData) *Frame {
	return &Frame{
		Length:    uint16(frameData.RawFrameData().Len()),
		FrameData: frameData.RawFrameData().Copy(),
		Checksum:  frameData.RawFrameData().Checksum(),
	}
}

// Deserialize creates a new frame from the byte slice.
// If the slice does not represent a valid frame, it returns an error.
func Deserialize(serializedFrame []byte) (*Frame, error) {
	if err := startDelimiterValid(serializedFrame); err != nil {
		return nil, err
	}
	if err := checksumVerifyFrame(serializedFrame); err != nil {
		return nil, err
	}

	data := frameDataField(serializedFrame)
	expectedLength, err := lengthField(serializedFrame)
	if err != nil {
		return nil, err
	}
	if int(expectedLength) != len(data) {
		return nil, &FrameParseError{
			msg: fmt.Sprintf("Expected length %u, received %u", expectedLength, len(data)),
		}
	}

	checksumIndex := len(serializedFrame) - 1
	frameData := NewRawFrameData(data...)

	if frameData.Checksum() != serializedFrame[checksumIndex] {
		return nil, &FrameParseError{msg: "Invalid checksum"}
	}

	f := &Frame{
		Length:    uint16(expectedLength),
		FrameData: NewRawFrameData(data...),
		Checksum:  serializedFrame[checksumIndex],
	}

	return f, nil
}

func (f *Frame) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(frameStartDelimiter)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, f.Length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, f.FrameData.buf)
	if err != nil {
		return nil, err
	}
	err = buf.WriteByte(f.FrameData.Checksum())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

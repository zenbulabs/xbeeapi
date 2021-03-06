package xbeeapi

import (
	"bytes"
	"testing"
)

func TestSimpleFrame(t *testing.T) {
	frameBytes := []byte{0x7e, 0x00, 0x07, 0x88, 0x01, 0x4d, 0x59, 0x00, 0x00, 0x00, 0xd0}
	frame, err := Deserialize(frameBytes)

	if err != nil {
		t.Error("Expected valid frame:", err)
		return
	}

	actualLen := uint16(frame.FrameData.Len())
	if frame.Length != actualLen {
		t.Error("Expected frame length", frame.Length, "but got", actualLen)
		return
	}

	serialized, err2 := frame.Serialize()
	if err2 != nil {
		t.Error("Error serializing frame")
		return
	}
	if !bytes.Equal(frameBytes, serialized) {
		t.Error("Serialization and deserialization mismatch")
		return
	}
}

func TestBadCheckSum(t *testing.T) {
	frameBytes := []byte{0x7e, 0x00, 0x07, 0x88, 0x01, 0x4d, 0x59, 0x00, 0x00, 0x00, 0xd1}
	_, err := Deserialize(frameBytes)

	if err == nil {
		t.Error("Expected bad checksum error")
		return
	}
}

func TestNewFrameFromData(t *testing.T) {
	frame := NewFrame(NewRawFrameData([]byte{FrameTypeATCommand, 0x7, 0x2}...))

	var expectedChecksum byte = 0xff - (FrameTypeATCommand + 0x7 + 0x2)
	if frame.Checksum != expectedChecksum {
		t.Error("Expected checksum", expectedChecksum, "but got", frame.Checksum)
	}
}

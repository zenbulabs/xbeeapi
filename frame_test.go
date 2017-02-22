package xbeeapi

import "testing"

func TestSimple(t *testing.T) {
	frameBytes := []byte{0x7e, 0x00, 0x07, 0x88, 0x01, 0x4d, 0x59, 0x00, 0x00, 0x00, 0xd0}
	frame, err := DeserializeFrame(frameBytes)

	if err != nil {
		t.Error("Expected valid frame:", err)
		return
	}

	actualLen := uint16(len(frame.Data) + 1)
	if frame.Length != actualLen {
		t.Error("Expected frame length", frame.Length, "but got", actualLen)
	}
}

func TestBadCheckSum(t *testing.T) {
	frameBytes := []byte{0x7e, 0x00, 0x07, 0x88, 0x01, 0x4d, 0x59, 0x00, 0x00, 0x00, 0xd1}
	_, err := DeserializeFrame(frameBytes)

	if err == nil {
		t.Error("Expected bad checksum error")
		return
	}
}

func TestNewFrameFromData(t *testing.T) {
	frame, err := NewFrameFromData(ATCommand, []byte{0x7, 0x2})
	if err != nil {
		t.Error("Expected new frame from data")
		return
	}

	var expectedChecksum byte = 0xff - (ATCommand + 0x7 + 0x2)
	if frame.Checksum != expectedChecksum {
		t.Error("Expected checksum", frame.Checksum, "but got", frame.Checksum)
	}
}

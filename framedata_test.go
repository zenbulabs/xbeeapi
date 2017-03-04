package xbeeapi

import (
	"bytes"
	"testing"
)

func TestATCommand(t *testing.T) {
	cmd := &ATCommand{FrameID: 1, Command: "AP", Params: nil}
	frame := NewFrame(cmd.RawFrameData())
	frameBytes, err := frame.Serialize()

	if err != nil {
		t.Error("Could not serialize AT command")
	}

	expectedFrameBytes := []byte{0x7e, 0x00, 0x04, 0x08, 0x01, 0x41, 0x50, 0x65}
	if !bytes.Equal(frameBytes, expectedFrameBytes) {
		t.Error("Expected:", expectedFrameBytes, "Got:", frameBytes)
	}
}

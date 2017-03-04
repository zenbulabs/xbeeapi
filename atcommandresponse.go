package xbeeapi

import (
	"bytes"
	"fmt"
)

const MinATCommandResponseSize = 4

const (
	ATCommandOK = iota
	ATCommandError
	ATCommandInvalidCommand
	ATCommandInvalidParam
	ATCommandRemoteTransFailed
	ATCommandStatusUnknown
)

type ATCommandResponse struct {
	FrameID byte
	Command string
	Status  byte
	Params  []byte
}

func ParseATCommandResponse(rfd *RawFrameData) (*ATCommandResponse, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeATCommandResponse {
		return nil, &FrameParseError{msg: "Expecting frame type ATCommandResponse"}
	}
	if len(rfd.Data()) < MinATCommandResponseSize {
		return nil, &FrameParseError{msg: "Frame data too small for ATCommandResponse"}
	}
	buf := bytes.NewBuffer(rfd.Data())
	at := &ATCommandResponse{
		FrameID: buf.Next(1)[0],
		Command: string(buf.Next(2)),
		Status:  buf.Next(1)[0],
		Params:  copySlice(buf.Bytes()),
	}
	if !at.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for ATCommandResponse"}
	}

	return at, nil
}

func (atr *ATCommandResponse) RawFrameData() *RawFrameData {
	rfd := concat([]byte{FrameTypeATCommandResponse, atr.FrameID}, []byte(atr.Command))
	rfd = append(rfd, byte(atr.Status))

	return NewRawFrameData(concat(rfd, atr.Params)...)
}

func (atr *ATCommandResponse) IsValid() bool {
	if atr.Status < ATCommandStatusUnknown && len(atr.Command) == 2 {
		return true
	}

	return false
}

func ATCommandStatusDescription(status byte) string {
	switch status {
	case ATCommandOK:
		return fmt.Sprintf("OK %d", status)
	case ATCommandError:
		return fmt.Sprintf("Error %d", status)
	case ATCommandInvalidCommand:
		return fmt.Sprintf("Invalid Command %d", status)
	case ATCommandInvalidParam:
		return fmt.Sprintf("Invalid Params%d", status)
	}

	return fmt.Sprintf("AT Command Status Unknown %d", status)
}

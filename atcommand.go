package xbeeapi

import "bytes"

const MinATCommandSize = 3

type ATCommand struct {
	FrameID byte
	Command string
	Params  []byte
}

func ParseATCommand(rfd *RawFrameData) (*ATCommand, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeATCommand {
		return nil, &FrameParseError{msg: "Expecting frame type ATCommand"}
	}
	if rfd.Len() < MinATCommandSize {
		return nil, &FrameParseError{msg: "Frame data too small for ATCommand"}
	}
	buf := bytes.NewBuffer(rfd.Data())
	at := &ATCommand{
		FrameID: buf.Next(1)[0],
		Command: string(buf.Next(2)),
		Params:  copySlice(buf.Bytes()),
	}
	if !at.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for ATCommand"}
	}

	return at, nil
}

func (at *ATCommand) RawFrameData() *RawFrameData {
	return NewRawFrameData(concat([]byte{FrameTypeATCommand, at.FrameID}, []byte(at.Command), at.Params)...)
}

func (at *ATCommand) IsValid() bool {
	if len(at.Command) == 2 {
		return true
	}

	return false
}

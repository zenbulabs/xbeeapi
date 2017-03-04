package xbeeapi

import "bytes"

type ATCommandQueue struct {
	FrameID byte
	Command string
	Params  []byte
}

func ParseATCommandQueue(rfd *RawFrameData) (*ATCommandQueue, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeATCommandQueueRegisterValue {
		return nil, &FrameParseError{msg: "Expecting frame type ATCommandQueue"}
	}
	if rfd.Len() < MinATCommandSize {
		return nil, &FrameParseError{msg: "Frame data too small for ATCommandQueue"}
	}
	buf := bytes.NewBuffer(rfd.Data())
	at := &ATCommandQueue{
		FrameID: buf.Next(1)[0],
		Command: string(buf.Next(2)),
		Params:  copySlice(buf.Bytes()),
	}
	if !at.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for ATCommandQueue"}
	}

	return at, nil
}

func (at *ATCommandQueue) RawFrameData() *RawFrameData {
	return NewRawFrameData(concat([]byte{FrameTypeATCommandQueueRegisterValue, at.FrameID}, []byte(at.Command), at.Params)...)
}

func (at *ATCommandQueue) IsValid() bool {
	if len(at.Command) == 2 {
		return true
	}

	return false
}

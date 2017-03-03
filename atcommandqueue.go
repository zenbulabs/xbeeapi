package xbeeapi

type ATCommandQueue struct {
	FrameID byte
	Command string
	Params  []byte
}

func ParseATCommandQueue(rfd *RawFrameData) (*ATCommandQueue, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeATCommandQueueRegisterValue {
		return nil, &FrameParseError{msg: "Expecting frame type AT command"}
	}
	if rfd.Len() < 3 {
		return nil, &FrameParseError{msg: "Frame data too small for AT command"}
	}
	at := &ATCommandQueue{FrameID: rfd.Data()[0], Command: string(rfd.Data()[1:3])}
	if rfd.Len() > 3 {
		at.Params = rfd.Data()[3:]
	}
	if !at.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for AT command"}
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

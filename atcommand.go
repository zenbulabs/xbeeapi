package xbeeapi

type ATCommand struct {
	frameData *RawFrameData
}

func NewATCommand(frameID byte, commandType string, params []byte) *ATCommand {
	if len([]byte(commandType)) != 2 {
		return nil
	}
	data := append([]byte{FrameTypeATCommand, frameID, commandType[0], commandType[1]}, params...)
	return &ATCommand{frameData: NewRawFrameData(data...)}
}

func ATCommandFrameData(rfd *RawFrameData) (*ATCommand, error) {
	at := &ATCommand{frameData: rfd}
	if !at.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for AT command"}
	}
	return at, nil
}

func (at *ATCommand) FrameID() byte {
	return at.frameData.buf[0]
}

func (at *ATCommand) CommandType() string {
	return string(at.frameData.buf[1:3])
}

func (at *ATCommand) Param() []byte {
	if len(at.frameData.buf) > 3 {
		return at.frameData.buf[3:]
	}
	return nil
}

func (at *ATCommand) RawFrameData() *RawFrameData {
	return at.frameData
}

func (at *ATCommand) IsValid() bool {
	return len(at.frameData.buf) >= 3
}

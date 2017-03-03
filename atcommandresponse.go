package xbeeapi

import "fmt"

type ATCommandStatus byte

const (
	ATCommandOK ATCommandStatus = iota
	ATCommandError
	ATCommandInvalidCommand
	ATCommandInvalidParam
	ATCommandRemoteTransFailed
)

type ATCommandResponse struct {
	frameData *RawFrameData
}

func NewATCommandResponse(frameID byte, commandType string, status ATCommandStatus, params []byte) *ATCommand {
	if len([]byte(commandType)) != 2 {
		return nil
	}
	data := append([]byte{FrameTypeATCommand, frameID, commandType[0], commandType[1]}, params...)
	return &ATCommand{frameData: NewRawFrameData(data...)}
}

func ATCommandResponseFrameData(rfd *RawFrameData) (*ATCommandResponse, error) {
	atr := &ATCommandResponse{frameData: rfd}
	if !atr.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for AT command response"}
	}
	return atr, nil
}

func (atr *ATCommandResponse) FrameID() byte {
	return atr.frameData.buf[0]
}

func (atr *ATCommandResponse) CommandType() string {
	return string(atr.frameData.buf[1:3])
}

func (atr *ATCommandResponse) Status() ATCommandStatus {
	return ATCommandStatus(atr.frameData.buf[3])
}

func (atr *ATCommandResponse) Param() []byte {
	if len(atr.frameData.buf) > 4 {
		return atr.frameData.buf[4:]
	}
	return nil
}

func (atr *ATCommandResponse) RawFrameData() *RawFrameData {
	return atr.frameData
}

func (atr *ATCommandResponse) IsValid() bool {
	return len(atr.frameData.buf) >= 4
}

func (ats ATCommandStatus) Description() string {
	switch ats {
	case ATCommandOK:
		return fmt.Sprintf("OK %d", ats)
	case ATCommandError:
		return fmt.Sprintf("Error %d", ats)
	case ATCommandInvalidCommand:
		return fmt.Sprintf("Invalid Command %d", ats)
	case ATCommandInvalidParam:
		return fmt.Sprintf("Invalid Params%d", ats)

	}
	return fmt.Sprintf("AT Command Status Unknown %d", ats)
}

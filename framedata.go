package xbeeapi

import "fmt"

const (
	FrameTypeTxRequest64                       = 0x00
	FrameTypeTxRequest16                       = 0x01
	FrameTypeATCommand                         = 0x08
	FrameTypeATCommandQueueRegisterValue       = 0x09
	FrameTypeTxRequest                         = 0x10
	FrameTypeExplicitAddressingCommandFrame    = 0x11
	FrameTypeTxSMS                             = 0x1f
	FrameTypeRemoteATCommand                   = 0x17
	FrameTypeTxIPv4                            = 0x20
	FrameTypeSendIPDataRequest                 = 0x28
	FrameTypeDeviceResponse                    = 0x2a
	FrameTypeRxPacket64                        = 0x80
	FrameTypeRxPacket16                        = 0x81
	FrameTypeRxPacketIO64                      = 0x82
	FrameTypeRxPacketIO16                      = 0x83
	FrameTypeWiFiRemoteATCommandResponse       = 0x87
	FrameTypeATCommandResponse                 = 0x88
	FrameTypeTxStatus                          = 0x89
	FrameTypeModemStatus                       = 0x8a
	FrameTypeXBTxStatus                        = 0x8b
	FrameTypeDigiMeshRouteInfoPacket           = 0x8d
	FrameTypeDigiMeshAggregateAddressingUpdate = 0x8e
	FrameTypeWifiIODataSampleRxIndicator       = 0x8f
	FrameTypeXBRxResponse                      = 0x90
	FrameTypeExplicitRxIndicator               = 0x91
	FrameTypeXBIODataSampleRxIndicator         = 0x92
	FrameTypeXBSensorReadIndicator             = 0x94
	FrameTypeXBNodeIdentificationIndicator     = 0x95
	FrameTypeRemoteATCommandResponse           = 0x97
	FrameTypeXBExtendedModemStatus             = 0x98
	FrameTypeRxSMS                             = 0x9f
	FrameTypeXBOTAFirmwareUpdateStatus         = 0xa0
	FrameTypeXBRouteRecordIndicator            = 0xa1
	FrameTypeXBManyToOneRouteRequestIndiator   = 0xa3
	FrameTypeXBJoinNotificationStatus          = 0xa5
	FrameTypeRxIPv4                            = 0xb0
	FrameTypeSendIPDataResponse                = 0xb8
	FrameTypeDeviceRequest                     = 0xb9
	FrameTypeDeviceResponseStatus              = 0xba
	FrameTypeFrameError                        = 0xfe
)

type FrameData interface {
	RawFrameData() *RawFrameData
	IsValid() bool
}

type RawFrameData struct {
	buf []byte
}

func NewRawFrameData(data ...byte) *RawFrameData {
	return &RawFrameData{buf: append([]byte(nil), data...)}
}

func (rfd *RawFrameData) FrameType() byte {
	return rfd.buf[0]
}

func (rfd *RawFrameData) Data() []byte {
	return rfd.buf[1:]
}

func (rfd *RawFrameData) Len() int {
	return len(rfd.buf)
}

func (rfd *RawFrameData) Copy() *RawFrameData {
	return &RawFrameData{buf: append([]byte(nil), rfd.buf...)}
}

func (rfd *RawFrameData) Checksum() byte {
	var cs byte = 0x0

	for _, b := range rfd.buf {
		cs += b
	}

	return 0xff - cs
}

func (rfd *RawFrameData) RawFrameData() *RawFrameData {
	return rfd
}

func (rfd *RawFrameData) IsValid() bool {
	return rfd.Len() > 0
}

func NewFrameData(fd FrameData) *RawFrameData {
	return NewRawFrameData(fd.RawFrameData().buf...)
}

func ParseFrameData(rfd *RawFrameData) (FrameData, error) {
	if rfd.Len() == 0 {
		return nil, &FrameParseError{msg: "Frame data not large enough"}
	}
	switch rfd.FrameType() {
	case FrameTypeATCommand:
		return ParseATCommand(rfd)
	case FrameTypeATCommandResponse:
		return ParseATCommandResponse(rfd)
	case FrameTypeModemStatus:
		return ParseModemStatus(rfd)
	case FrameTypeATCommandQueueRegisterValue:
		return ParseATCommandQueue(rfd)
	case FrameTypeExplicitAddressingCommandFrame:
		return ParseTxExplicitAddressing(rfd)
	}
	return nil, &FrameParseError{msg: fmt.Sprintf("Unsupported frame type: %02x", rfd.FrameType())}
}

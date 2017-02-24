package xbeeapi

import (
	"bytes"
	"encoding/binary"
)

const (
	TxRequest64                       = 0x00
	TxRequest16                       = 0x01
	ATCommand                         = 0x08
	ATCommandQueueRegisterValue       = 0x09
	TxRequest                         = 0x10
	ExplicitAddressingCommandFrame    = 0x11
	TxSMS                             = 0x1f
	RemoteATCommand                   = 0x17
	TxIPv4                            = 0x20
	SendIPDataRequest                 = 0x28
	DeviceResponse                    = 0x2a
	RxPacket64                        = 0x80
	RxPacket16                        = 0x81
	RxPacketIO64                      = 0x82
	RxPacketIO16                      = 0x83
	WiFiRemoteATCommandResponse       = 0x87
	ATCommandResponse                 = 0x88
	TxStatus                          = 0x89
	ModemStatus                       = 0x8a
	XBTxStatus                        = 0x8b
	DigiMeshRouteInfoPacket           = 0x8d
	DigiMeshAggregateAddressingUpdate = 0x8e
	WifiIODataSampleRxIndicator       = 0x8f
	XBRxResponse                      = 0x90
	XBExplicitRxIndicator             = 0x91
	XBIODataSampleRxIndicator         = 0x92
	XBSensorReadIndicator             = 0x94
	XBNodeIdentificationIndicator     = 0x95
	RemoteATCommandResponse           = 0x97
	XBExtendedModemStatus             = 0x98
	RxSMS                             = 0x9f
	XBOTAFirmwareUpdateStatus         = 0xa0
	XBRouteRecordIndicator            = 0xa1
	XBManyToOneRouteRequestIndiator   = 0xa3
	XBJoinNotificationStatus          = 0xa5
	RxIPv4                            = 0xb0
	SendIPDataResponse                = 0xb8
	DeviceRequest                     = 0xb9
	DeviceResponseStatus              = 0xba
	FrameError                        = 0xfe
)

const (
	FrameStartDelimiter = 0x7e
)

type Frameable interface {
	FrameType() byte
	FrameData() []byte
}

type Frame struct {
	Length    uint16
	FrameType byte
	Data      []byte
	Checksum  byte
}

type FrameParseError struct {
	msg string
}

func (e *FrameParseError) Error() string { return e.msg }

func checksumVerify(serializedFrame []byte) bool {
	var cs byte

	for _, b := range serializedFrame[3:] {
		cs += b
	}

	return cs == 0xff
}

func checksumFromData(frameType byte, data []byte) byte {
	cs := frameType

	for _, b := range data {
		cs += b
	}

	return 0xff - cs
}

func lengthFromHeader(serializedFrame []byte) uint16 {
	return uint16(int(serializedFrame[1])<<8) + uint16(serializedFrame[2])
}

func totalFrameSize(dataLen uint16) int {
	// start(1) + datalen(2) + data(variable) + checksum(1)
	return 3 + int(dataLen) + 1
}

func dataWithoutType(serializedFrame []byte) []byte {
	return serializedFrame[4:(len(serializedFrame) - 1)]
}

func DeserializeFrame(serializedFrame []byte) (*Frame, error) {
	if len(serializedFrame) < 6 {
		return nil, &FrameParseError{msg: "Frame too short"}
	} else if serializedFrame[0] != FrameStartDelimiter {
		return nil, &FrameParseError{msg: "Expecting 7e as start delimiter"}
	} else if !checksumVerify(serializedFrame) {
		return nil, &FrameParseError{msg: "Invalid check sum"}
	}

	data := dataWithoutType(serializedFrame)
	expectedLength := int(lengthFromHeader(serializedFrame))
	if expectedLength != len(data)+1 {
		return nil, &FrameParseError{msg: "Expected length does not match actual length"}
	}

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	checksumIndex := len(serializedFrame) - 1
	f := Frame{
		Length:    uint16(expectedLength),
		FrameType: serializedFrame[3],
		Data:      dataCopy,
		Checksum:  serializedFrame[checksumIndex],
	}

	return &f, nil
}

func NewFrameFromData(frameType byte, data []byte) (*Frame, error) {
	if len(data) == 0 {
		return nil, &FrameParseError{msg: "Data is empty"}
	} else if len(data) > 0xffff {
		return nil, &FrameParseError{msg: "Data length too long"}
	}

	dataCopy := make([]byte, len(data), len(data))
	copy(dataCopy, data)

	f := Frame{
		Length:    uint16(len(data) + 1),
		FrameType: frameType,
		Data:      dataCopy,
		Checksum:  checksumFromData(frameType, data),
	}

	return &f, nil
}

func (f *Frame) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(FrameStartDelimiter)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, f.Length)
	if err != nil {
		return nil, err
	}
	err = buf.WriteByte(f.FrameType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, f.Data)
	if err != nil {
		return nil, err
	}
	err = buf.WriteByte(checksumFromData(f.FrameType, f.Data))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

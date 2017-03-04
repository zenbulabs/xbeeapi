package xbeeapi

import (
	"bytes"
	"encoding/binary"
)

const MinRxExplicitIndicatorSize = 29

type RxExplicitIndicator struct {
	FrameID     byte
	Address64   string
	Address16   string
	SrcEndPoint byte
	DstEndPoint byte
	ClusterID   uint16
	ProfileID   uint16
	Options     byte
	Payload     []byte
}

func ParseRxExplicitIndicator(rfd *RawFrameData) (*RxExplicitIndicator, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeExplicitRxIndicator {
		return nil, &FrameParseError{msg: "Expecting frame type RxExplicitIndicator"}
	}
	if rfd.Len() < MinRxExplicitIndicatorSize {
		return nil, &FrameParseError{msg: "Frame data too small for RxExplicitIndicator"}
	}
	buf := bytes.NewBuffer(rfd.Data())

	tx := &RxExplicitIndicator{
		FrameID:     buf.Next(1)[0],
		Address64:   bytesToHex(buf.Next(16)),
		Address16:   bytesToHex(buf.Next(4)),
		SrcEndPoint: buf.Next(1)[0],
		DstEndPoint: buf.Next(1)[0],
		ClusterID:   binary.BigEndian.Uint16(buf.Next(2)),
		ProfileID:   binary.BigEndian.Uint16(buf.Next(2)),
		Options:     buf.Next(1)[0],
		Payload:     copySlice(buf.Bytes()),
	}

	if !tx.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for RxExplicitIndicator"}
	}

	return tx, nil
}

func (rx *RxExplicitIndicator) RawFrameData() *RawFrameData {
	b := []byte{FrameTypeExplicitRxIndicator, rx.FrameID}
	address64, _ := hexToBytes(rx.Address64)
	address16, _ := hexToBytes(rx.Address16)
	b = concat(b, address64, address16)
	b = append(b, rx.SrcEndPoint, rx.DstEndPoint, 0x00, 0x00, 0x00, 0x00)
	binary.BigEndian.PutUint16(b[(len(b)-4):], rx.ClusterID)
	binary.BigEndian.PutUint16(b[(len(b)-2):], rx.ProfileID)
	b = append(b, rx.Options)
	b = concat(b, rx.Payload)

	return NewRawFrameData(b...)
}

func (rx *RxExplicitIndicator) IsValid() bool {
	address64, _ := hexToBytes(rx.Address64)
	address16, _ := hexToBytes(rx.Address16)
	if len(address64) == 16 && len(address16) == 4 {
		return true
	}

	return false
}

func (rx *RxExplicitIndicator) FrameType() byte {
	return FrameTypeExplicitRxIndicator
}

func (rx *RxExplicitIndicator) SetOptionsFlags(rxOptionFlags ...RxOptionFlag) {
	rx.Options = setRxOptionsFlags(rx.Options, rxOptionFlags...)
}

func (rx *RxExplicitIndicator) IsOptionsFlagSet(rxOptionFlag RxOptionFlag) bool {
	return isRxOptionsFlagSet(rx.Options, rxOptionFlag)
}

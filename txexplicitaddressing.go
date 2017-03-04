package xbeeapi

import (
	"bytes"
	"encoding/binary"
)

const MinTxExplicitAddressingSize = 30

type TxExplicitAddressing struct {
	FrameID         byte
	Address64       string
	Address16       string
	SrcEndPoint     byte
	DstEndPoint     byte
	ClusterID       uint16
	ProfileID       uint16
	BroadcastRadius byte
	Options         byte
	Payload         []byte
}

func ParseTxExplicitAddressing(rfd *RawFrameData) (*TxExplicitAddressing, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeExplicitAddressingCommandFrame {
		return nil, &FrameParseError{msg: "Expecting frame type TxExplicitAddressing"}
	}
	if rfd.Len() < MinTxExplicitAddressingSize {
		return nil, &FrameParseError{msg: "Frame data too small for TxExplicitAddressing"}
	}
	buf := bytes.NewBuffer(rfd.Data())

	tx := &TxExplicitAddressing{
		FrameID:         buf.Next(1)[0],
		Address64:       bytesToHex(buf.Next(16)),
		Address16:       bytesToHex(buf.Next(4)),
		SrcEndPoint:     buf.Next(1)[0],
		DstEndPoint:     buf.Next(1)[0],
		ClusterID:       binary.BigEndian.Uint16(buf.Next(2)),
		ProfileID:       binary.BigEndian.Uint16(buf.Next(2)),
		BroadcastRadius: buf.Next(1)[0],
		Options:         buf.Next(1)[0],
		Payload:         copySlice(buf.Bytes()),
	}

	if !tx.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for TxExplicitAddressing"}
	}

	return tx, nil
}

func (tx *TxExplicitAddressing) RawFrameData() *RawFrameData {
	b := []byte{FrameTypeExplicitAddressingCommandFrame, tx.FrameID}
	address64, _ := hexToBytes(tx.Address64)
	address16, _ := hexToBytes(tx.Address16)
	b = concat(b, address64, address16)
	b = append(b, tx.SrcEndPoint, tx.DstEndPoint, 0x00, 0x00, 0x00, 0x00)
	binary.BigEndian.PutUint16(b[(len(b)-4):], tx.ClusterID)
	binary.BigEndian.PutUint16(b[(len(b)-2):], tx.ProfileID)
	b = append(b, tx.BroadcastRadius, tx.Options)
	b = concat(b, tx.Payload)

	return NewRawFrameData(b...)
}

func (tx *TxExplicitAddressing) IsValid() bool {
	address64, _ := hexToBytes(tx.Address64)
	address16, _ := hexToBytes(tx.Address16)
	if len(address64) == 16 && len(address16) == 4 {
		return true
	}

	return false
}

func (tx *TxExplicitAddressing) FrameType() byte {
	return FrameTypeExplicitAddressingCommandFrame
}

func (tx *TxExplicitAddressing) SetOptionsFlags(txOptionFlags ...TxOptionFlag) {
	tx.Options = setTxOptionsFlags(tx.Options, txOptionFlags...)
}

func (tx *TxExplicitAddressing) IsOptionsFlagSet(txOptionFlag TxOptionFlag) bool {
	return isTxOptionsFlagSet(tx.Options, txOptionFlag)
}

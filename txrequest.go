package xbeeapi

import "bytes"

const MinTxRequestSize = 23

type TxRequest struct {
	FrameID         byte
	Address64       string
	Address16       string
	BroadcastRadius byte
	Options         byte
	Payload         []byte
}

func ParseTxRequest(rfd *RawFrameData) (*TxRequest, error) {
	if !rfd.IsValid() || rfd.FrameType() != FrameTypeTxRequest {
		return nil, &FrameParseError{msg: "Expecting frame type TxRequest"}
	}
	if rfd.Len() < MinTxRequestSize {
		return nil, &FrameParseError{msg: "Frame data too small for TxRequest"}
	}
	buf := bytes.NewBuffer(rfd.Data())
	tx := &TxRequest{
		FrameID:         buf.Next(1)[0],
		Address64:       bytesToHex(buf.Next(16)),
		Address16:       bytesToHex(buf.Next(4)),
		BroadcastRadius: buf.Next(1)[0],
		Options:         buf.Next(1)[0],
		Payload:         copySlice(buf.Bytes()),
	}
	if !tx.IsValid() {
		return nil, &FrameParseError{msg: "Invalid frame data for AT command"}
	}

	return tx, nil
}

func (tx *TxRequest) RawFrameData() *RawFrameData {
	b := []byte{FrameTypeTxRequest, tx.FrameID}
	address64, _ := hexToBytes(tx.Address64)
	address16, _ := hexToBytes(tx.Address16)
	b = concat(b, address64, address16)
	b = append(b, tx.BroadcastRadius, tx.Options)
	b = concat(b, tx.Payload)

	return NewRawFrameData(b...)
}

func (tx *TxRequest) IsValid() bool {
	address64, _ := hexToBytes(tx.Address64)
	address16, _ := hexToBytes(tx.Address16)
	if len(address64) == 16 && len(address16) == 4 {
		return true
	}

	return false
}

func (tx *TxRequest) FrameType() byte {
	return FrameTypeTxRequest
}

func (tx *TxRequest) SetOptionsFlags(txOptionFlags ...TxOptionFlag) {
	tx.Options = setTxOptionsFlags(tx.Options, txOptionFlags...)
}

func (tx *TxRequest) IsOptionsFlagSet(txOptionFlag TxOptionFlag) bool {
	return isTxOptionsFlagSet(tx.Options, txOptionFlag)
}

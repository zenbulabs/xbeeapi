package xbeeapi

import (
	"errors"
	"io"
)

type XBeeStatusType byte
type ReadCallback func(frame *Frame, status *XBeeReadStatus)

const (
	XBeeOK XBeeStatusType = iota
	XBeeClose
	XBeeReadError
	XBeeUnknownStatus
)

type XBeeReadStatus struct {
	Status XBeeStatusType
	Error  error
}

type XBeeAPI struct {
	port   io.ReadWriter
	readCb ReadCallback
	done   chan bool
}

func NewXBeeAPI(serialPort io.ReadWriter, readCb ReadCallback) *XBeeAPI {
	return &XBeeAPI{
		port:   serialPort,
		readCb: readCb,
		done:   make(chan bool, 1),
	}
}

func (api *XBeeAPI) Connect() {
	go func() {
		defer close(api.done)
		for {
			select {
			case done := <-api.done:
				if done {
					return
				}
			default:
				api.readFrames()
			}
		}
	}()
}

func (api *XBeeAPI) SendRawFrame(f *Frame) (int, error) {
	frameBytes, err := f.Serialize()
	if err != nil {
		return 0, err
	}

	return api.port.Write(frameBytes)
}

func (api *XBeeAPI) SendFrame(f Frameable) (int, error) {
	frame, err := NewFrameFromData(f.FrameType(), f.FrameData())
	if err != nil {
		return 0, err
	}
	return api.SendRawFrame(frame)
}

func (api *XBeeAPI) readFrames() {
	frameBytes := make([]byte, 1024, 1024)
	n, err := api.port.Read(frameBytes[0:1])
	if n == 0 || err != nil || frameBytes[0] != FrameStartDelimiter {
		api.readCb(nil, &XBeeReadStatus{Status: XBeeReadError, Error: errors.New("Frame read: invalid start delimiter")})
		return
	}
	n, err = api.port.Read(frameBytes[1:3])

	if n != 2 || err != nil {
		api.readCb(nil, &XBeeReadStatus{Status: XBeeReadError, Error: errors.New("Frame read: invalid length field")})
		return
	}
	dataLen := lengthFromHeader(frameBytes[1:3])

	if dataLen == 0 && dataLen > 1024 {
		api.readCb(nil, &XBeeReadStatus{Status: XBeeReadError, Error: errors.New("Frame read: invalid data size")})
		return
	}
	totalSize := 3 + dataLen + 1
	n, err = api.port.Read(frameBytes[3:totalSize])

	if n == int(dataLen+1) && err == nil {
		frame, parseErr := DeserializeFrame(frameBytes[0:totalSize])
		if parseErr == nil {
			api.readCb(frame, &XBeeReadStatus{Status: XBeeOK})
		} else {
			api.readCb(nil, &XBeeReadStatus{Status: XBeeReadError, Error: parseErr})
		}
	} else {
		api.readCb(nil, &XBeeReadStatus{Status: XBeeReadError, Error: errors.New("Frame read: missing data or checksum")})
	}
}

func (api *XBeeAPI) Finish() {
	api.done <- true
}

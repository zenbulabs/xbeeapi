package xbeeapi

import (
	"io"
)

type XBeeStatusType byte
type ReadCallback func(frame *Frame, status *XBeeStatus)

const (
	XBeeOK XBeeStatusType = iota
	XBeeClose
	XBeeReadError
	XBeeUnknownStatus
)

type Frameable interface {
	Frame() *Frame
}

type XBeeReadStatus struct {
	Status XBeeStatusType
	Error  error
}

type XBeeAPI struct {
	Port         *io.ReadWrite
	readCallback ReadCallback
	done         chan bool
}

func NewXBeeAPI(serialPort *io.ReadWrite, readCb ReadCallback) {
	return &XBeeAPI{
		port:         serialPort,
		readCallback: readCallback,
		done: make(chan bool, 1)
	}
}

func (api *XBeeApi) Connect() {
	go func() {
		defer close(done)
		for {
			select {
			case done := <-api.frameChan:
				if done {
					return
				}
			default:
				api.readFrames()
			}
		}
	}()
}

func (api *XBeeApi) SendRawFrame(f *Frame) (int, error) {
	frameBytes, err := f.Serialize()
	if error != nil {
		return 0, err
	}

	return api.port.Write(frameBytes)
}

func (api *XBeeApi) SendFrame(f Frameable) (int, error) {
	return SendRawFrame(f.Frame())
}

func (api *XBeeApi) readFrames() {
	frameBytes := make([]byte, 1024, 1024)
	total := 0
	n, err := api.port.Read(frameBytes[0:1])
	if n == 0 || err != nil || frameBytes[0] != 0x7e {
		api.readCb(nil, &XBeeStatus{Status: XBeeReadError, Error: errors.New("Frame read: invalid start delimiter")})
		return
	}
	n, err = api.port.Read(frameBytes[1:3])

	if n != 2 || err != nil {
		api.readCb(nil, &XBeeStatus{Status: XBeeReadError, Error: errors.New("Frame read: invalid length field")})
		return
	}
	dataLen := lengthFromHeader(frameBytes[1:3])

	if dataLen == 0 && dataLen > 1024 {
		api.readCb(nil, &XBeeStatus{Status: XBeeReadError, Error: errors.New("Frame read: invalid data size")})
		return
	}
	totalSize := 3 + dataLen + 1
	n, err = api.port.Read(frameBytes[3:totalSize])

	if n == dataLen+1 && err == nil {
		frame, parseErr := DeserializeFrame(frameBytes[0:totalSize])
		if parseErr == nil {
			api.readCb(frame, &XBeeStatus{Status: XBeeOK})
		} else {
			api.readCb(nil, &XBeeStatus{Status: XBeeReadError, Error: parseErr})
		}
	} else {
		api.readCb(nil, &XBeeStatus{Status: XBeeReadError, Error: errors.New("Frame read: missing data or checksum")})
	}
}

func (api *XBeeApi) Finish() {
	api.done <- true
}

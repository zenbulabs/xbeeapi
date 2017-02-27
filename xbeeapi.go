package xbeeapi

import (
	"errors"
	"io"
	"log"
	"sync"
	"time"
)

type XBeeStatusType byte

type ReadCallback func(frame *Frame, status XBeeReadStatus)

const (
	XBeeOK XBeeStatusType = iota
	XBeeClose
	XBeeReadError
	XBeeUnknownStatus
)

type XBeeReadStatus struct {
	StatusCode XBeeStatusType
	Error      error
}

type XBeeAPI struct {
	fwr     *frameReadWriter
	readCb  ReadCallback
	mu      *sync.Mutex
	running bool
}

func NewXBeeAPI(port io.ReadWriter, readCb ReadCallback) *XBeeAPI {
	return &XBeeAPI{
		fwr:     newFrameReader(port),
		readCb:  readCb,
		mu:      &sync.Mutex{},
		running: false,
	}
}

func (api *XBeeAPI) Start() error {
	var err error

	api.mu.Lock()
	if api.running {
		err = errors.New("XBeeAPI already started")
	} else {
		api.running = true
	}
	api.mu.Unlock()

	if err != nil {
		return err
	}

	api.fwr.init()

	go func() {
		for {
			if !api.Running() {
				log.Println("Stopping...")
				return
			}
			err := api.readFrames()
			if err != nil {
				time.Sleep(200 * time.Millisecond)
				log.Println("Reader error:", err)
			}
		}
	}()

	return nil
}

func (api *XBeeAPI) SendRawFrames(frame ...*Frame) (int, error) {
	n, err := api.fwr.write(frame...)
	return n, err
}

func (api *XBeeAPI) SendFrames(frameData ...FrameData) (int, error) {
	frames := []*Frame(nil)

	for _, fd := range frameData {
		frames = append(frames, NewFrame(fd))
	}

	return api.SendRawFrames(frames...)
}

func (api *XBeeAPI) readFrames() error {
	frames, err := api.fwr.read()

	if err != nil {
		api.readCb(nil, XBeeReadStatus{StatusCode: XBeeReadError, Error: err})
		return err
	}

	for _, frame := range frames {
		api.readCb(frame, XBeeReadStatus{StatusCode: XBeeOK, Error: nil})
	}

	return nil
}

func (api *XBeeAPI) Running() (r bool) {
	api.mu.Lock()
	r = api.running
	api.mu.Unlock()
	return
}

func (api *XBeeAPI) Finish() {
	api.mu.Lock()
	api.running = false
	api.fwr.init()
	api.mu.Unlock()
}

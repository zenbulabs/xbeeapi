package xbeeapi

import (
	"bytes"
	"testing"
)

type TestPort struct {
	data *bytes.Buffer
}

type TestData struct {
	frame  chan *Frame
	status XBeeReadStatus
}

func (t *TestData) readCb(f *Frame, s XBeeReadStatus) {
	t.status = s
	t.frame <- f
}

func NewTestPort(b []byte) *TestPort {
	return &TestPort{data: bytes.NewBuffer(b)}
}

func (p *TestPort) Read(data []byte) (int, error) {
	return p.data.Read(data)
}

func (p *TestPort) Write(data []byte) (int, error) {
	return p.data.Write(data)
}

func TestStart(t *testing.T) {
	nf := NewFrame(NewRawFrameData([]byte{0x88, 0x01, 0x4d, 0x59, 0x00, 0x00, 0x00}...))
	testFrameBytes, _ := nf.Serialize()
	port := NewTestPort(testFrameBytes)
	td := &TestData{frame: make(chan *Frame, 1)}

	api := NewXBeeAPI(port, td.readCb)
	err := api.Start()
	if err != nil {
		t.Error("Could not start", err)
		return
	}
	err = api.Start()
	if err == nil {
		t.Error("Expected error from double start")
		return
	}
	f := <-td.frame
	if f == nil {
		t.Error("Invalid frame from read")
		t.Error(td.status)
		return
	}

	api.Finish()
}

func TestWrite(t *testing.T) {
	port := NewTestPort([]byte{})
	td := &TestData{frame: make(chan *Frame)}
	api := NewXBeeAPI(port, td.readCb)
	frameSend := NewFrame(NewRawFrameData([]byte{0x0f, 0x02, 0x04, 0x06}...))
	n, err := api.SendRawFrames(frameSend)
	if n == 0 || err != nil {
		t.Error("SendRawFrame error", n, err)
	}
	api.Start()

	frameRecv := <-td.frame
	f1, err1 := frameSend.Serialize()
	f2, err2 := frameRecv.Serialize()

	if err1 != nil || err2 != nil || !bytes.Equal(f1, f2) {
		t.Error("Error in sending frames.", "Sent:", f1, "Received:", f2)
	}

	api.Finish()
}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	b, _ := NewFrameFromData(0x88, []byte{0x01, 0x4d, 0x59, 0x00, 0x00, 0x00})
	//testFrameBytes := []byte{0x7e, 0x00, 0x07, 0x88, 0x01, 0x4d, 0x59, 0x00, 0x00, 0x00, 0xd0}
	testFrameBytes, _ := b.Serialize()
	port := NewTestPort(testFrameBytes)
	td := &TestData{frame: make(chan *Frame, 1)}
	api := NewXBeeAPI(port, td.readCb)
	api.Start()

	f := <-td.frame
	if f == nil {
		t.Error("Invalid frame from read")
		t.Error(td.status)
	}

	api.Finish()
}

func TestWrite(t *testing.T) {
	t.SkipNow()
	port := NewTestPort([]byte{})
	td := &TestData{frame: make(chan *Frame)}
	api := NewXBeeAPI(port, td.readCb)
	frame, _ := NewFrameFromData(0x0f, []byte{0x02, 0x04, 0x06})
	n, err := api.SendRawFrame(frame)
	if n == 0 || err != nil {
		t.Error("SendRawFrame error", n, err)
	}
	api.Start()

	f := <-td.frame
	if f != frame {
		t.Error("Invalid frame", td.status)
	}

	api.Finish()
}

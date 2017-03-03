package xbeeapi

import (
	"bytes"
	"io"
	"sync"
)

type frameReadWriter struct {
	rw  io.ReadWriter
	mu  *sync.Mutex
	buf *bytes.Buffer
}

func newFrameReader(rw io.ReadWriter) *frameReadWriter {
	return &frameReadWriter{
		rw:  rw,
		mu:  &sync.Mutex{},
		buf: bytes.NewBuffer([]byte{}),
	}
}

func (fr *frameReadWriter) read() ([]*Frame, error) {
	var b [16]byte
	n, err := fr.rw.Read(b[:])

	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}

	n, err = fr.buf.Write(b[:n])

	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}

	frames := make([]*Frame, 0)

	for {
		buf := fr.buf.Bytes()
		if len(buf) < minFrameSize {
			break
		}

		if buf[0] != frameStartDelimiter {
			fr.buf.ReadByte()
			continue
		}

		dataLen, err := lengthField(buf)
		if err != nil {
			break
		}

		frameLen := totalFrameLength(dataLen)
		if len(buf) < frameLen {
			break
		}

		buf = fr.buf.Next(frameLen)
		frame, err := Deserialize(buf)

		if err == nil {
			frames = append(frames, frame)
		}
	}

	if fr.buf.Cap() > 128 {
		// Shrink buffer if capacity grows too large
		newBuffer := bytes.NewBuffer([]byte{})
		newBuffer.Write(fr.buf.Bytes())
		fr.buf = newBuffer
	}

	return frames, nil
}

func (fr *frameReadWriter) write(frames ...*Frame) (int, error) {
	var err error
	totalWritten := 0

	//fr.mu.Lock()
	for _, f := range frames {
		frameBytes, err := f.Serialize()
		if err != nil {
			break
		}
		_, err = fr.rw.Write(frameBytes)
		if err != nil {
			break
		}
		totalWritten++
	}
	//fr.mu.Unlock()

	return totalWritten, err
}

func (fr *frameReadWriter) init() (int, error) {
	return fr.rw.Write([]byte{0x7e, 0x00, 0x04, 0x08, 0x01, 0x41, 0x50, 0x65})
}

package main

import (
	"bytes"
	"net"
	"time"
)

type DummyCon struct {
	data     []byte
	offset   int
	dataSize int

	writeBuffer *bytes.Buffer
}

func NewDummyCon(data []byte) *DummyCon {
	return &DummyCon{
		data,
		0,
		len(data),
		new(bytes.Buffer)}
}

func (m *DummyCon) Read(b []byte) (n int, err error) {
	toRead := len(b)
	i := 0
	for ; i < toRead && m.offset+i < m.dataSize; i++ {
		b[i] = m.data[m.offset+i]
	}
	m.offset += i

	return i, nil
}

func (m *DummyCon) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *DummyCon) Close() error {
	return nil
}

func (m *DummyCon) LocalAddr() net.Addr {
	return nil
}

func (m *DummyCon) RemoteAddr() net.Addr {
	return nil
}

func (m *DummyCon) SetDeadline(t time.Time) error {
	return nil
}

func (m *DummyCon) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *DummyCon) SetWriteDeadline(t time.Time) error {
	return nil
}

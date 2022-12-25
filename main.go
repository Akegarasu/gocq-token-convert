package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

type token struct {
	D2Key []byte
	D2    []byte
	TGT   []byte
}

type Reader struct {
	buf *bytes.Reader
}

func NewReader(data []byte) *Reader {
	buf := bytes.NewReader(data)
	return &Reader{
		buf: buf,
	}
}
func (r *Reader) ReadBytes(len int) []byte {
	b := make([]byte, len)
	if len > 0 {
		_, err := r.buf.Read(b)
		if err != nil {
			panic(err)
		}
	}
	return b
}

func (r *Reader) ReadBytesShort() []byte {
	return r.ReadBytes(int(r.ReadUInt16()))
}

func (r *Reader) ReadInt64() int64 {
	b := make([]byte, 8)
	_, _ = r.buf.Read(b)
	return int64(binary.BigEndian.Uint64(b))
}

func (r *Reader) ReadUInt16() uint16 {
	b := make([]byte, 2)
	_, _ = r.buf.Read(b)
	return binary.BigEndian.Uint16(b)
}

func readToken(data []byte) *token {
	rd := NewReader(data)
	t := &token{}
	rd.ReadInt64()
	t.D2 = rd.ReadBytesShort()
	t.D2Key = rd.ReadBytesShort()
	t.TGT = rd.ReadBytesShort()
	return t
}

func main() {
	data, err := os.ReadFile("session.token")
	if err != nil {
		log.Fatal(err)
	}
	t := readToken(data)
	buf := new(bytes.Buffer)
	buf.Write(t.D2Key)
	buf.Write(t.D2)
	buf.Write(t.TGT)
	err = os.WriteFile("token", buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

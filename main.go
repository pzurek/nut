package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Constant declarations
const (
	Null    = 0xC0 // 192
	Float64 = 0xC1 // 193
	False   = 0xC2 // 194
	True    = 0xC3 // 195

	Int8  = 0xC8 // 200
	Int16 = 0xC9 // 201
	Int32 = 0xCA // 202
	Int64 = 0xCB // 203

	Bytes8  = 0xCC // 204
	Bytes16 = 0xCD // 205
	Bytes32 = 0xCE // 206

	String8  = 0xD0 // 208
	String16 = 0xD1 // 209
	String32 = 0xD2 // 210

	List8      = 0xD4 // 212
	List16     = 0xD5 // 213
	List32     = 0xD6 // 214
	ListStream = 0xD7 // 215

	Map8      = 0xD8 // 216
	Map16     = 0xD9 // 217
	Map32     = 0xDA // 218
	MapStream = 0xDB // 219

	Struct8  = 0xDC // 220
	Struct16 = 0xDD // 221

	EndOfStream = 0xDF // 223

	TinyTextStart   = 0x80 // 128
	TinyTextEnd     = 0x8F // 143
	TinyListStart   = 0x90 // 144
	TinyListEnd     = 0x9F // 159
	TinyMapStart    = 0xA0 // 160
	TinyMapEnd      = 0xAF // 175
	TinyStructStart = 0xB0 // 176
	TinyStructEnd   = 0xBF // 191

	MinTinyInt = -16
	MaxTinyInt = 127
)

var ()

func init() {
	neo4jMessages := map[byte]string{
		0x01: "INIT",
		0x0E: "ACK_FAILURE",
		0x0F: "RESET",
		0x10: "RUN",
		0x2F: "DISCARD_ALL",
		0x3F: "PULL_ALL",
		0x70: "SUCCESS",
		0x71: "RECORD",
		0x7E: "IGNORED",
		0x7F: "FAILURE",
	}
}

func main() {

}

// Decoder decodes
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a new Decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Encoder encodes
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new Encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Decode returns
func (decoder Decoder) Decode() (interface{}, error) {
	r := bufio.NewReader(decoder.r)
	marker, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	markerHigh := marker & 0xF0

	switch marker {
	case Null:
		return nil, nil
	case False:
		return false, nil
	case True:
		return true, nil
	case Float64:
		return decodeFloat(r)

		// Ints
	case Int8:
		return decodeInt(r, 1)
	case Int16:
		return decodeInt(r, 2)
	case Int32:
		return decodeInt(r, 4)
	case Int64:
		return decodeInt(r, 8)

		// Bytes
	case Bytes8:
		return decodeByte(r, 1)
	case Bytes16:
		return decodeByte(r, 2)
	case Bytes32:
		return decodeByte(r, 2)

		// Strings
	case String8:
		decodeString(r, 1)
	case String16:
		decodeString(r, 2)
	case String32:
		decodeString(r, 4)

	case List8:

	case List16:

	case List32:

	case ListStream:

	case Map8:

	case Map16:

	case Map32:

	case MapStream:

	case Struct8:

	case Struct16:

	case EndOfStream:
	}

	return nil, fmt.Errorf("decoding error: unsupported type: %x", marker)
}

func decodeByte(r io.Reader, size int) ([]byte, error) {
	b := make([]byte, size)
	n, err := r.Read(b[:])
	if err != nil {
		return b, err
	}
	if n != cap(b) {
		return b, fmt.Errorf("failed to read all bytes")
	}

	return b, nil
}

func decodeString(r io.Reader, size int) (string, error) {
	s := ""
	b, err := decodeByte(r, size)
	if err != nil {
		return s, err
	}
	return string(b), nil
}

func decodeInt(r io.Reader, size int) (int, error) {
	i := 0
	b, err := decodeByte(r, size)
	if err != nil {
		return i, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

func decodeFloat(r io.Reader) (float64, error) {
	f := 0.0
	b, err := decodeByte(r, 8)
	if err != nil {
		return f, err
	}

	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.BigEndian, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}

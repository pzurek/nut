package packstream

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Encoder encodes
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new Encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) write(bytes ...byte) error {
	n, err := e.w.Write(bytes)
	if n < len(bytes) {
		return fmt.Errorf("failed to write all bytes")
	}
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) encodeInt(i int) error {
	// special casefor TinyInt
	if i >= MinTinyInt && i <= MaxTinyInt {
		return e.write(uint8(i)) // casting the byte to uint to let it overflow
	}

	b := []byte{}
	var err error

	switch {
	case i < MinTinyInt && i >= math.MinInt8, i > MaxTinyInt && i <= math.MaxInt8:
		b, err = numberToBytes(int8(i))
	case i < math.MinInt8 && i >= math.MinInt16, i > math.MaxInt8 && i <= math.MaxInt16:
		b, err = numberToBytes(int16(i))
	case i < math.MinInt16 && i >= math.MinInt32, i > math.MaxInt16 && i <= math.MaxInt32:
		b, err = numberToBytes(int32(i))
	default:
		b, err = numberToBytes(int64(i))
	}
	if err != nil {
		return err
	}
	err = e.write(b...)
	return err
}

func (e *Encoder) encodeNull() error {
	return e.write(Null)
}

func (e *Encoder) encodeBool(v bool) error {
	if v {
		return e.write(True)
	}
	return e.write(False)
}

func (e *Encoder) encodeFloat64(f float64) error {
	b, err := numberToBytes(f)
	if err != nil {
		return err
	}
	return e.write(b...)
}

func (e *Encoder) encodeString(s string) error {
	return e.write([]byte(s)...)
}

// numberToBytes is a little wrapper around binary.Write
// it takes any fixed size number (int8, int64, float64 etc, but not int)
// and returns BigEndianbytes
func numberToBytes(n interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	err := binary.Write(b, binary.BigEndian, n)
	return b.Bytes(), err
}

func (e *Encoder) encodeStringHeader(s string) error {
	l := len([]byte(s))
	return e.write(byte(l))
}

func (e *Encoder) encodeTinyIntHeader(i int) error {
	b := uint8(i)
	return e.write(b)
}

func (e *Encoder) encodeInt8Header() error {
	return e.write(Int8)
}

func (e *Encoder) encodeInt16Header() error {
	return e.write(Int16)
}

func (e *Encoder) encodeInt32Header() error {
	return e.write(Int32)
}

func (e *Encoder) encodeInt64Header() error {
	return e.write(Int64)
}

func (e *Encoder) encodeBoolHeader(b bool) error {
	if b {
		return e.write(True)
	}
	return e.write(False)
}

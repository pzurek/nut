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

func (encoder *Encoder) Write(bytes ...byte) error {
	n, err := encoder.w.Write(bytes)
	if n < len(bytes) {
		return fmt.Errorf("failed to write all bytes")
	}
	if err != nil {
		return err
	}
	return nil
}

func bytesFromInt(i interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	err := binary.Write(b, binary.BigEndian, i)
	return b.Bytes(), err
}

func encodeInt(i int) (interface{}, error) {
	switch {
	case i >= MinTinyInt && i <= MaxTinyInt:
		return byte(i), nil
	case i < MinTinyInt && i >= math.MinInt8, i > MaxTinyInt && i <= math.MaxInt8:
		return bytesFromInt(int8(i))
	case i < math.MinInt8 && i >= math.MinInt16, i > math.MaxInt8 && i <= math.MaxInt16:
		return bytesFromInt(int16(i))
	case i < math.MinInt16 && i >= math.MinInt32, i > math.MaxInt16 && i <= math.MaxInt32:
		return bytesFromInt(int32(i))
	}
	return bytesFromInt(int64(i))
}

package packstream

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Decoder decodes
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a new Decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode returns
func (decoder *Decoder) Decode() (interface{}, error) {
	r := bufio.NewReader(decoder.r)
	marker, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	// markerHighNibble := marker & 0xF0
	// markerLowNibble := marker & 0x0F

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
		return decodeBytes(r, 1)
	case Bytes16:
		return decodeBytes(r, 2)
	case Bytes32:
		return decodeBytes(r, 2)

		// Strings
	case String8:
		return decodeString(r, 1)
	case String16:
		return decodeString(r, 2)
	case String32:
		return decodeString(r, 4)

	case List8:
		return nil, nil
	case List16:
		return nil, nil
	case List32:
		return nil, nil
	case ListStream:
		return nil, nil
	case Map8:
		return nil, nil
	case Map16:
		return nil, nil
	case Map32:
		return nil, nil
	case MapStream:
		return nil, nil
	case Struct8:
		return nil, nil
	case Struct16:
		return nil, nil
	case EndOfStream:
		return nil, nil
	}
	return nil, fmt.Errorf("decoding error: unsupported type: %x", marker)
}

func decodeBytes(r io.Reader, size int) ([]byte, error) {
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
	b, err := decodeBytes(r, size)
	if err != nil {
		return s, err
	}
	return string(b), nil
}

func decodeInt(r io.Reader, size int) (int, error) {
	i := 0
	b, err := decodeBytes(r, size)
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
	b, err := decodeBytes(r, 8)
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

func (decoder *Decoder) peekNextType() Type {
	reader := bufio.NewReader(decoder.r)
	markerbytes, err := reader.Peek(1)
	if err != nil {
		return PSNull
	}
	marker := markerbytes[0]
	markerHighNibble := marker & 0xF0

	switch markerHighNibble {
	case TinyString:
		return PSString
	case TinyList:
		return PSList
	case TinyMap:
		return PSMap
	case TinyStruct:
		return PSStruct
	}

	switch marker {
	case Null:
		return PSNull
	case True:
	case False:
		return PSBool
	case Float64:
		return PSFloat
	case Bytes8:
	case Bytes16:
	case Bytes32:
		return PSBytes
	case String8:
	case String16:
	case String32:
		return PSString
	case List8:
	case List16:
	case List32:
		return PSList
	case Struct8:
	case Struct16:
		return PSStruct
	default:
		return PSInt
	}
	return PSNull
}

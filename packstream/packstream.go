package packstream

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// PackStream marker declarations
const (
	TinyString = 0x80 // 128
	TinyList   = 0x90 // 144
	TinyMap    = 0xA0 // 160
	TinyStruct = 0xB0 // 176

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

	MinTinyInt = -16
	MaxTinyInt = 127
)

type Type int

// PackStream types
const (
	PSNull PackStreamType = iota
	PSBool
	PSInt
	PSFloat
	PSBytes
	PSString
	PSList
	PSMap
	PSStruct
)

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

func encodeInt(i int) (interface{}, error) {
	switch {
	case i >= MinTinyInt && i <= MaxTinyInt:
		return byte(i), nil
	case i < MinTinyInt && i >= math.MinInt8:
	case i > MaxTinyInt && i <= math.MaxInt8:
		return bytesFromInt(int8(i))
	case i < math.MinInt8 && i >= math.MinInt16:
	case i > math.MaxInt8 && i <= math.MaxInt16:
		return bytesFromInt(int16(i))
	case i < math.MinInt16 && i >= math.MinInt32:
	case i > math.MaxInt16 && i <= math.MaxInt32:
		return bytesFromInt(int32(i))
	}
	return bytesFromInt(int64(i))
}

func bytesFromInt(i interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	err := binary.Write(b, binary.BigEndian, i)
	return b.Bytes(), err
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

func (decoder *Decoder) peekNextType() PackStreamType {
	reader := bufio.NewReader(decoder.r)
	marker, err := reader.Peek(1)
	if err != nil {
		return PSNull
	}
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
}

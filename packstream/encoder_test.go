package packstream

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		i int
		b []byte
	}{
		// tiny int
		{-16, []byte{0xf0}},
		{-1, []byte{0xff}},
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{21, []byte{0x15}},
		{42, []byte{0x2a}},
		{127, []byte{0x7f}},

		// int 8
		{-128, []byte{0xff}},
		{-42, []byte{0x7f, 0xd6}},
		{-17, []byte{0x7f}},

		// int 16
		{-32768, []byte{0x80, 0x00}},
		{-129, []byte{0xff, 0x7f}},
		// {128, []byte{0x7f}},
		// {32767, []byte{0x7f}},

		// int 32
		{-2147483648, []byte{0x80, 0x00, 0x00, 0x00}},
		// {-32769, []byte{0x7f}},
		// {32768, []byte{0x7f}},
		// {2147483647, []byte{0x7f}},

		// int 64
		// {-9223372036854775808, []byte{0x7f}},
		// {-2147483649, []byte{0x7f}},
		// {2147483648, []byte{0x7f}},
		// {9223372036854775807, []byte{0x7f}},

	}

	buf := &bytes.Buffer{}
	e := NewEncoder(buf)

	for _, test := range tests {
		err := e.encodeInt(test.i)
		if err != nil {
			t.Error(err)
		}
		bytes, err := ioutil.ReadAll(buf)
		if err != nil {
			t.Error(err)
		} else {
			for i := range bytes {
				if bytes[i] != test.b[i] {
					t.Errorf("expected: %#v, got: %#v\n", test.b[i], bytes[i])
				}
			}
		}
	}
}

//
// func testencodeNull(t testing.T) {
// }
//
// func testencodeBool(t testing.T) {
//
// }

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
		{-16, []byte{0xf0}},
		{-1, []byte{0xff}},
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{21, []byte{0x15}},
		{127, []byte{0x7f}},
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

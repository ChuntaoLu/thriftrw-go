// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package protocol

import (
	"bytes"
	"math"
	"testing"

	"github.com/uber/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

var encodeDecodeTests = []struct {
	value    wire.Value
	expected []byte
}{
	// bool
	{vbool(false), []byte{0x00}},
	{vbool(true), []byte{0x01}},

	// byte
	{vbyte(0), []byte{0x00}},
	{vbyte(1), []byte{0x01}},
	{vbyte(-1), []byte{0xff}},
	{vbyte(127), []byte{0x7f}},
	{vbyte(-128), []byte{0x80}},

	// i16
	{vi16(1), []byte{0x00, 0x01}},
	{vi16(255), []byte{0x00, 0xff}},
	{vi16(256), []byte{0x01, 0x00}},
	{vi16(257), []byte{0x01, 0x01}},
	{vi16(32767), []byte{0x7f, 0xff}},
	{vi16(-1), []byte{0xff, 0xff}},
	{vi16(-2), []byte{0xff, 0xfe}},
	{vi16(-256), []byte{0xff, 0x00}},
	{vi16(-255), []byte{0xff, 0x01}},
	{vi16(-32768), []byte{0x80, 0x00}},

	// i32
	{vi32(1), []byte{0x00, 0x00, 0x00, 0x01}},
	{vi32(255), []byte{0x00, 0x00, 0x00, 0xff}},
	{vi32(65535), []byte{0x00, 0x00, 0xff, 0xff}},
	{vi32(16777215), []byte{0x00, 0xff, 0xff, 0xff}},
	{vi32(2147483647), []byte{0x7f, 0xff, 0xff, 0xff}},
	{vi32(-1), []byte{0xff, 0xff, 0xff, 0xff}},
	{vi32(-256), []byte{0xff, 0xff, 0xff, 0x00}},
	{vi32(-65536), []byte{0xff, 0xff, 0x00, 0x00}},
	{vi32(-16777216), []byte{0xff, 0x00, 0x00, 0x00}},
	{vi32(-2147483648), []byte{0x80, 0x00, 0x00, 0x00}},

	// i64
	{vi64(1), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
	{vi64(4294967295), []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}},
	{vi64(1099511627775), []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff}},
	{vi64(281474976710655), []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	{vi64(72057594037927935), []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	{vi64(9223372036854775807), []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	{vi64(-1), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	{vi64(-4294967296), []byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00}},
	{vi64(-1099511627776), []byte{0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{vi64(-281474976710656), []byte{0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{vi64(-72057594037927936), []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{vi64(-9223372036854775808), []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},

	// double
	{vdouble(0.0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{vdouble(1.0), []byte{0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{vdouble(1.0000000001), []byte{0x3f, 0xf0, 0x0, 0x0, 0x0, 0x6, 0xdf, 0x38}},
	{vdouble(1.1), []byte{0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
	{vdouble(-1.1), []byte{0xbf, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
	{vdouble(3.141592653589793), []byte{0x40, 0x9, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}},
	{vdouble(-1.0000000001), []byte{0xbf, 0xf0, 0x0, 0x0, 0x0, 0x6, 0xdf, 0x38}},
	{vdouble(math.NaN()), []byte{0x7f, 0xf8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}},
	{vdouble(math.Inf(0)), []byte{0x7f, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
	{vdouble(math.Inf(-1)), []byte{0xff, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},

	// binary~4
	{vbinary(""), []byte{0x00, 0x00, 0x00, 0x00}},
	{vbinary("hello"), []byte{
		0x00, 0x00, 0x00, 0x05, // len:4 = 5
		0x68, 0x65, 0x6c, 0x6c, 0x6f, // 'h', 'e', 'l', 'l', 'o'
	}},

	// struct = (type:1 id:2 value)* stop
	// stop = 0
	{vstruct(), []byte{0x00}},
	{vstruct(vfield(1, vbool(true))), []byte{
		0x02,       // type:1 = bool
		0x00, 0x01, // id:2 = 1
		0x01, // value = true
		0x00, // stop
	}},
	{
		vstruct(
			vfield(1, vi16(42)),
			vfield(2, vlist(wire.TBinary, vbinary("foo"), vbinary("bar"))),
		), []byte{
			0x06,       // type:1 = i16
			0x00, 0x01, // id:2 = 1
			0x00, 0x2a, // value = 42

			0x0F,       // type:1 = list
			0x00, 0x02, // id:2 = 2

			// <list>
			0x0B,                   // type:1 = binary
			0x00, 0x00, 0x00, 0x02, // size:4 = 2
			// <binary>
			0x00, 0x00, 0x00, 0x03, // len:4 = 3
			0x66, 0x6f, 0x6f, // 'f', 'o', 'o'
			// </binary>
			// <binary>
			0x00, 0x00, 0x00, 0x03, // len:4 = 3
			0x62, 0x61, 0x72, // 'b', 'a', 'r'
			// </binary>
			// </list>

			0x00, // stop
		},
	},

	// map = ktype:1 vtype:1 count:4 (key value){count}
	{vmap(wire.TI64, wire.TBinary), []byte{0x0A, 0x0B, 0x00, 0x00, 0x00, 0x00}},
	{
		vmap(
			wire.TBinary, wire.TList,
			vitem(vbinary("a"), vlist(wire.TI16, vi16(1))),
			vitem(vbinary("b"), vlist(wire.TI16, vi16(2), vi16(3))),
		), []byte{
			0x0B,                   // ktype = binary
			0x0F,                   // vtype = list
			0x00, 0x00, 0x00, 0x02, // count:4 = 2

			// <item>
			// <key>
			0x00, 0x00, 0x00, 0x01, // len:4 = 1
			0x61, // 'a'
			// </key>
			// <value>
			0x06,                   // type:1 = i16
			0x00, 0x00, 0x00, 0x01, // count:4 = 1
			0x00, 0x01, // 1
			// </value>
			// </item>

			// <item>
			// <key>
			0x00, 0x00, 0x00, 0x01, // len:4 = 1
			0x62, // 'b'
			// </key>
			// <value>
			0x06,                   // type:1 = i16
			0x00, 0x00, 0x00, 0x02, // count:4 = 2
			0x00, 0x02, // 2
			0x00, 0x03, // 3
			// </value>
			// </item>
		},
	},

	// set = vtype:1 count:4 (value){count)
	{vset(wire.TBool), []byte{0x02, 0x00, 0x00, 0x00, 0x00}},
	{
		vset(wire.TBool, vbool(true), vbool(false), vbool(true)),
		[]byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x01},
	},

	// list = vtype:1 count:4 (value){count}
	{vlist(wire.TStruct), []byte{0x0C, 0x00, 0x00, 0x00, 0x00}},
	{
		vlist(
			wire.TStruct,
			vstruct(
				vfield(1, vi16(1)),
				vfield(2, vi32(2)),
			),
			vstruct(
				vfield(1, vi16(3)),
				vfield(2, vi32(4)),
			),
		),
		[]byte{
			0x0C,                   // vtype:1 = struct
			0x00, 0x00, 0x00, 0x02, // count:4 = 2

			// <struct>
			0x06,       // type:1 = i16
			0x00, 0x01, // id:2 = 1
			0x00, 0x01, // value = 1

			0x08,       // type:1 = i32
			0x00, 0x02, // id:2 = 2
			0x00, 0x00, 0x00, 0x02, // value = 2

			0x00, // stop
			// </struct>

			// <struct>
			0x06,       // type:1 = i16
			0x00, 0x01, // id:2 = 1
			0x00, 0x03, // value = 3

			0x08,       // type:1 = i32
			0x00, 0x02, // id:2 = 2
			0x00, 0x00, 0x00, 0x04, // value = 4

			0x00, // stop
			// </struct>
		},
	},
}

func TestEncode(t *testing.T) {
	for _, tt := range encodeDecodeTests {
		buffer := bytes.Buffer{}
		err := Binary.Encode(tt.value, &buffer)
		if assert.NoError(t, err, "Encode failed:\n%s", tt.value) {
			assert.Equal(t, tt.expected, buffer.Bytes())
		}
	}
}

// TODO test input too short errors

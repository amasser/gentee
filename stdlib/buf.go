// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/gentee/gentee/core"
)

// InitBuffer appends stdlib buffer functions to the virtual machine
func InitBuffer(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{core.Link{bufºStr, 44<<16 | core.EMBED}, `str`, `buf`}, // buf( str ) buf
		{core.Link{strºBuf, 45<<16 | core.EMBED}, `buf`, `str`}, // str( buf ) str
		{core.Link{LenºBuf, core.Bcode(core.TYPEBUF<<16) | core.LEN},
			`buf`, `int`}, // the length of the buffer
		{core.Link{AddºBufBuf, 50<<16 | core.EMBED}, `buf,buf`, `buf`},        // buf + buf
		{core.Link{AssignºBufBuf, core.ASSIGN}, `buf,buf`, `buf`},             // buf = buf
		{core.Link{AssignAddºBufInt, core.ASSIGN + 4}, `buf,int`, `buf`},      // buf += int
		{core.Link{AssignAddºBufStr, core.ASSIGN + 6}, `buf,str`, `buf`},      // buf += str
		{core.Link{AssignAddºBufChar, core.ASSIGN + 3}, `buf,char`, `buf`},    // buf += char
		{core.Link{AssignAddºBufBuf, core.ASSIGN + 5}, `buf,buf`, `buf`},      // buf += buf
		{core.Link{AssignBitAndºBufBuf, core.ASSIGNPTR}, `buf,buf`, `buf`},    // buf &= buf
		{core.Link{Base64ºBuf, 51<<16 | core.EMBED}, `buf`, `str`},            // Base64( buf ) str
		{core.Link{DelºBufIntInt, 52<<16 | core.EMBED}, `buf,int,int`, `buf`}, // Del( buf, int, int ) buf
		{core.Link{HexºBuf, 46<<16 | core.EMBED}, `buf`, `str`},               // Hex( buf ) str
		{core.Link{InsertºBufIntBuf, 47<<16 | core.EMBED},
			`buf,int,buf`, `buf`}, // Insert( buf, int, buf ) buf
		{core.Link{UnBase64ºStr, 48<<16 | core.EMBED}, `str`, `buf`}, // UnBase64( str ) buf
		{core.Link{UnHexºStr, 49<<16 | core.EMBED}, `str`, `buf`},    // UnHex( str ) buf
		{core.Link{sysBufNil, 53<<16 | core.EMBED}, ``, `buf`},       // sysBufNil() buf
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// AddºBufBuf adds two buffers
func AddºBufBuf(left *core.Buffer, right *core.Buffer) (out *core.Buffer) {
	out = core.NewBuffer()
	out.Data = left.Data
	out.Data = append(out.Data, right.Data...)
	return out
}

// AssignºBufBuf copies one buf to another one
func AssignºBufBuf(ptr *interface{}, value *core.Buffer) *core.Buffer {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Buffer)
}

// AssignAddºBufChar appends rune to buffer
func AssignAddºBufChar(ptr *interface{}, value rune) *core.Buffer {
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, []byte(string([]rune{value}))...)
	return (*ptr).(*core.Buffer)
}

// AssignAddºBufInt appends one byte to buffer
func AssignAddºBufInt(ptr *interface{}, value int64) (*core.Buffer, error) {
	if uint64(value) > 255 {
		return nil, fmt.Errorf(core.ErrorText(core.ErrByteOut))
	}
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, byte(value))
	return (*ptr).(*core.Buffer), nil
}

// AssignAddºBufBuf appends buffer to buffer
func AssignAddºBufBuf(ptr *interface{}, value *core.Buffer) *core.Buffer {
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, value.Data...)
	return (*ptr).(*core.Buffer)
}

// AssignAddºBufStr appends string to buffer
func AssignAddºBufStr(ptr *interface{}, value string) *core.Buffer {
	(*ptr).(*core.Buffer).Data = append((*ptr).(*core.Buffer).Data, []byte(value)...)
	return (*ptr).(*core.Buffer)
}

// AssignBitAndºBufBuf assigns a pointer to buffer
func AssignBitAndºBufBuf(ptr *interface{}, value *core.Buffer) *core.Buffer {
	*ptr = value
	return (*ptr).(*core.Buffer)
}

// bufºStr converts string to buffer
func bufºStr(value string) *core.Buffer {
	b := core.NewBuffer()
	b.Data = []byte(value)
	return b
}

// strºBuf converts buffer to string
func strºBuf(buf *core.Buffer) string {
	return string(buf.Data)
}

// Base64ºBuf encodes buf to base64 string
func Base64ºBuf(buf *core.Buffer) string {
	return base64.StdEncoding.EncodeToString(buf.Data)
}

// HexºBuf encodes buf to hex string
func HexºBuf(buf *core.Buffer) string {
	return hex.EncodeToString(buf.Data)
}

// DelºBufIntInt deletes part of the buffer
func DelºBufIntInt(buf *core.Buffer, off, length int64) (*core.Buffer, error) {
	size := int64(len(buf.Data))
	if off < 0 || off > size {
		return buf, fmt.Errorf(core.ErrorText(core.ErrInvalidParam))
	}
	if length < 0 {
		off += length
		length = -length
	}
	if off < 0 {
		off = 0
	}
	if off+length > size {
		length = size - off
	}
	buf.Data = append(buf.Data[:off], buf.Data[off+length:]...)
	return buf, nil
}

// InsertºBufIntBuf inserts one buf object into another one
func InsertºBufIntBuf(buf *core.Buffer, off int64, b *core.Buffer) (*core.Buffer, error) {
	size := int64(len(buf.Data))
	if off < 0 || off > size {
		return buf, fmt.Errorf(core.ErrorText(core.ErrInvalidParam))
	}
	buf.Data = append(buf.Data[:off], append(b.Data, buf.Data[off:]...)...)
	return buf, nil
}

// LenºBuf returns the length of the buffer
func LenºBuf(buf *core.Buffer) int64 {
	return int64(len(buf.Data))
}

// UnBase64ºStr decodes base64 string to buf
func UnBase64ºStr(value string) (buf *core.Buffer, err error) {
	buf = core.NewBuffer()
	buf.Data, err = base64.StdEncoding.DecodeString(value)
	return
}

// UnHexºStr decodes hex string to the buffer
func UnHexºStr(value string) (*core.Buffer, error) {
	var err error
	buf := core.NewBuffer()
	buf.Data, err = hex.DecodeString(value)
	return buf, err
}

// sysBufNil return nil buffer
func sysBufNil() *core.Buffer {
	b := core.NewBuffer()
	b.Data = nil
	return b
}

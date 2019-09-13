// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"crypto/md5"
	"crypto/sha256"

	"github.com/gentee/gentee/core"
)

// InitCrypto appends stdlib crypto functions to the virtual machine
func InitCrypto(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{core.Link{Md5ºBuf, 58<<16 | core.EMBED}, `buf`, `buf`},    // Md5( buf ) buf
		{core.Link{Md5ºStr, 59<<16 | core.EMBED}, `str`, `buf`},    // Md5( str ) buf
		{core.Link{Sha256ºBuf, 60<<16 | core.EMBED}, `buf`, `buf`}, // Sha256( buf ) buf
		{core.Link{Sha256ºStr, 61<<16 | core.EMBED}, `str`, `buf`}, // Sha256( str ) buf
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

func md5Hash(in []byte) (out *core.Buffer) {
	out = core.NewBuffer()
	md5Hash := md5.Sum(in)
	out.Data = md5Hash[:]
	return out
}

// Md5ºBuf returns md5 hash of the buffer
func Md5ºBuf(in *core.Buffer) (out *core.Buffer) {
	return md5Hash(in.Data)
}

// Md5ºStr returns md5 hash of the string as hex string
func Md5ºStr(in string) (out *core.Buffer) {
	return md5Hash([]byte(in))
}

func sha256Hash(in []byte) (out *core.Buffer) {
	out = core.NewBuffer()
	shaHash := sha256.Sum256(in)
	out.Data = shaHash[:]
	return out
}

// Sha256ºBuf returns md5 hash of the buffer
func Sha256ºBuf(in *core.Buffer) (out *core.Buffer) {
	return sha256Hash(in.Data)
}

// Sha256ºStr returns md5 hash of the string as hex string
func Sha256ºStr(in string) (out *core.Buffer) {
	return sha256Hash([]byte(in))
}

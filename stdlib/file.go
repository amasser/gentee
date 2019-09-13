// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/gentee/gentee/core"
)

// InitFile appends stdlib int functions to the virtual machine
func InitFile(ws *core.Workspace) {
	NewStructType(ws, `finfo`, []string{
		`Name:str`, `Size:int`, `Mode:int`,
		`Time:time`, `IsDir:bool`,
	})

	for _, item := range []interface{}{
		core.Link{AppendFileºStrBuf, 69<<16 | core.EMBED},  // AppendFile( str, buf )
		core.Link{AppendFileºStrStr, 70<<16 | core.EMBED},  // AppendFile( str, str )
		core.Link{ChDirºStr, 71<<16 | core.EMBED},          // ChDir( str )
		core.Link{CopyFileºStrStr, 72<<16 | core.EMBED},    // CopyFile( str, str )
		core.Link{CreateDirºStr, 73<<16 | core.EMBED},      // CreateDir( str )
		core.Link{GetCurDir, 74<<16 | core.EMBED},          // GetCurDir( ) str
		core.Link{Md5FileºStr, 75<<16 | core.EMBED},        // Md5File( str ) str
		core.Link{ReadFileºStr, 40<<16 | core.EMBED},       // ReadFile( str ) str
		core.Link{ReadFileºStrBuf, 76<<16 | core.EMBED},    // ReadFile( str, buf ) buf
		core.Link{ReadFileºStrIntInt, 77<<16 | core.EMBED}, // ReadFile( str, int, int ) buf
		core.Link{RemoveºStr, 78<<16 | core.EMBED},         // Remove( str )
		core.Link{RemoveDirºStr, 79<<16 | core.EMBED},      // RemoveDir( str )
		core.Link{RenameºStrStr, 80<<16 | core.EMBED},      // Rename( str, str )
		core.Link{Sha256FileºStr, 81<<16 | core.EMBED},     // Sha256File( str ) str
		core.Link{TempDir, 82<<16 | core.EMBED},            // TempDir()
		core.Link{TempDirºStrStr, 83<<16 | core.EMBED},     // TempDir(str, str)
		core.Link{WriteFileºStrBuf, 84<<16 | core.EMBED},   // WriteFile( str, buf )
		core.Link{WriteFileºStrStr, 85<<16 | core.EMBED},   // WriteFile( str, str )
	} {
		ws.StdLib().NewEmbed(item)
	}
	for _, item := range []embedInfo{
		{core.Link{FileInfoºStr, 1016<<16 | core.EMBED}, `str`, `finfo`},        // FileInfo( str ) finfo
		{core.Link{ReadDirºStr, 1017<<16 | core.EMBED}, `str`, `arr.finfo`},     // ReadDir( str ) arr.finfo
		{core.Link{SetFileTimeºStrTime, 1018<<16 | core.EMBED}, `str,time`, ``}, // SetFileTime( str, time )
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

func appendFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func fromFileInfo(fileInfo os.FileInfo, finfo *core.Struct) *core.Struct {
	finfo.Values[0] = fileInfo.Name()
	finfo.Values[1] = fileInfo.Size()
	finfo.Values[2] = fileInfo.Mode()
	fromTime(finfo.Values[3].(*core.Struct), fileInfo.ModTime())
	finfo.Values[4] = fileInfo.IsDir()
	return finfo
}

// AppendFileºStrBuf appends a buffer to a file
func AppendFileºStrBuf(filename string, buf *core.Buffer) error {
	return appendFile(filename, buf.Data)
}

// AppendFileºStrStr appends a string to a file
func AppendFileºStrStr(filename, s string) error {
	return appendFile(filename, []byte(s))
}

// ChDirºStr change the current directory
func ChDirºStr(dirname string) error {
	return os.Chdir(dirname)
}

// CopyFileºStrStr copies a file
func CopyFileºStrStr(src, dest string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()
	return io.Copy(destFile, srcFile)
}

// GetCurDir returns the current directory
func GetCurDir() (string, error) {
	return os.Getwd()
}

// CreateDirºStr creates the directory(s)
func CreateDirºStr(dirname string) error {
	return os.MkdirAll(dirname, os.ModePerm)
}

// FileInfoºStr returns the finfo describing the named file.
func FileInfoºStr(rt *core.RunTime, name string) (*core.Struct, error) {
	finfo := core.NewStructObj(rt, `finfo`)
	handle, err := os.Open(name)
	if err != nil {
		return finfo, err
	}
	defer handle.Close()
	fileInfo, err := handle.Stat()
	if err != nil {
		return finfo, err
	}
	return fromFileInfo(fileInfo, finfo), nil
}

// Md5FileºStr returns md5 hash of the file as a hex string
func Md5FileºStr(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ``, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ``, err
	}
	return hex.EncodeToString(hash.Sum(nil)[:]), nil
}

// ReadDirºStr reads a directory
func ReadDirºStr(rt *core.RunTime, dirname string) (*core.Array, error) {
	ret := core.NewArray()
	fileList, err := ioutil.ReadDir(dirname)
	if err != nil {
		return ret, err
	}
	for _, fileInfo := range fileList {
		ret.Data = append(ret.Data, fromFileInfo(fileInfo, core.NewStructObj(rt, `finfo`)))
	}
	return ret, nil
}

// ReadFileºStr reads a file
func ReadFileºStr(filename string) (string, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

// ReadFileºStrBuf reads a file to buffer
func ReadFileºStrBuf(filename string, buf *core.Buffer) (*core.Buffer, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return buf, err
	}
	buf.Data = out
	return buf, nil
}

// ReadFileºStrIntInt reads a part of the file to the buffer
func ReadFileºStrIntInt(filename string, off int64, length int64) (buf *core.Buffer, err error) {
	var (
		fhandle *os.File
		n       int
	)
	buf = core.NewBuffer()
	if fhandle, err = os.Open(filename); err != nil {
		return
	}
	defer fhandle.Close()
	fi, err := fhandle.Stat()
	fsize := fi.Size()
	if off < 0 {
		off = fsize + off
	}
	if off < 0 {
		off = 0
	} else if off > fsize-1 {
		return
	}
	if off+length > fsize {
		length = fsize - off
	}
	buf.Data = make([]byte, length)
	n, err = fhandle.ReadAt(buf.Data, off)
	if err != nil && err == io.EOF {
		err = nil
	}
	buf.Data = buf.Data[:n]
	return
}

// RenameºStrStr renames a file or a directory
func RenameºStrStr(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// RemoveºStr removes a file or an empty directory
func RemoveºStr(filename string) error {
	return os.Remove(filename)
}

// RemoveDirºStr removes a directory
func RemoveDirºStr(dirname string) error {
	return os.RemoveAll(dirname)
}

// SetFileTimeºStrTime changes the modification time of the named file
func SetFileTimeºStrTime(name string, ftime *core.Struct) error {
	mtime := toTime(ftime)
	return os.Chtimes(name, mtime, mtime)
}

// Sha256FileºStr returns sha256 hash of the file as a hex string
func Sha256FileºStr(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ``, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ``, err
	}
	return hex.EncodeToString(hash.Sum(nil)[:]), nil
}

// TempDir returns the temporary directory
func TempDir() string {
	return os.TempDir()
}

// TempDirºStrStr creates a directory in the temporary directory
func TempDirºStrStr(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}

// WriteFileºStrBuf writes a buffer to a file
func WriteFileºStrBuf(filename string, buf *core.Buffer) error {
	return ioutil.WriteFile(filename, buf.Data, os.ModePerm)
}

// WriteFileºStrStr writes a string to a file
func WriteFileºStrStr(filename, in string) error {
	return ioutil.WriteFile(filename, []byte(in), os.ModePerm)
}

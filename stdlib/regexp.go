// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"regexp"

	"github.com/gentee/gentee/core"
)

// InitRegExp appends stdlib regexp functions to the virtual machine
func InitRegExp(ws *core.Workspace) {
	for _, item := range []interface{}{
		core.Link{MatchºStrStr, 64<<16 | core.EMBED},         // Match( str, str ) bool
		core.Link{ReplaceRegExpºStrStr, 65<<16 | core.EMBED}, // ReplaceRegExp( str, str ) str
	} {
		ws.StdLib().NewEmbed(item)
	}
	for _, item := range []embedInfo{
		{core.Link{FindRegExpºStrStr, 63<<16 | core.EMBED}, `str,str`, `arr.arr.str`}, // FindRegExp( str, str ) arr.arr.str
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// FindRegExpºStrStr returns an array of all successive matches of the expression
func FindRegExpºStrStr(src, rePattern string) (*core.Array, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return nil, err
	}
	list := re.FindAllStringSubmatch(src, -1)
	out := core.NewArray()
	for _, ilist := range list {
		sub := core.NewArray()
		for _, sublist := range ilist {
			sub.Data = append(sub.Data, sublist)
		}
		out.Data = append(out.Data, sub)
	}
	return out, nil
}

// MatchºStrStr reports whether the string s contains any match of the regular expression
func MatchºStrStr(s string, rePattern string) (bool, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

// ReplaceRegExpºStrStr returns a copy of src, replacing matches of the Regexp with the replacement string
func ReplaceRegExpºStrStr(src, rePattern, repl string) (string, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return ``, err
	}
	return re.ReplaceAllString(src, repl), nil
}

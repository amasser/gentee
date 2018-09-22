// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"reflect"

	"github.com/gentee/gentee/core"
)

type initType struct {
	name     string
	original reflect.Type
	index    string // support index of
}

// InitTypes appends stdlib types to the virtual machine
func InitTypes(vm *core.VirtualMachine) {
	for _, item := range []initType{
		{`int`, reflect.TypeOf(int64(0)), ``},
		{`bool`, reflect.TypeOf(true), ``},
		{`char`, reflect.TypeOf('a'), ``},
		{`str`, reflect.TypeOf(``), `char`},
		{`range`, reflect.TypeOf(core.Range{}), `int`},
	} {
		vm.StdLib().NewType(item.name, item.original, item.index)
	}
}

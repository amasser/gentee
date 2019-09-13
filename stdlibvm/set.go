// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

/*
// AssignºSetSet copies one set to another one
func AssignºSetSet(ptr *interface{}, value *core.Set) *core.Set {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Set)
}

// AssignBitAndºSetSet assigns a pointer to set
func AssignBitAndºSetSet(ptr *interface{}, value *core.Set) *core.Set {
	*ptr = value
	return (*ptr).(*core.Set)
}

// AssignAddºSetSet appends set to set
func AssignAddºSetSet(ptr *interface{}, value *core.Set) *core.Set {
	for i, v := range value.Data {
		for pos := uint64(0); pos < 64; pos++ {
			if v&(1<<pos) != 0 {
				(*ptr).(*core.Set).Set(int64(i<<6)+int64(pos), true)
			}
		}
	}
	return (*ptr).(*core.Set)
}

// LenºSet returns the length of the set size
func LenºSet(set *core.Set) int64 {
	return int64(len(set.Data) << 6)
}
*/
func checkIndex(set *core.Set, index int64) error {
	if index < 0 || index >= core.MaxSet {
		return fmt.Errorf(core.ErrorText(core.ErrIndexOut))
	}
	return nil
}

/*
func bitSet(left *core.Set, right *core.Set, and bool) *core.Set {
	ret := core.NewSet()
	if len(left.Data) < len(right.Data) {
		left, right = right, left
	}
	ret.Data = make([]uint64, len(left.Data))
	for i, v := range left.Data {
		if i < len(right.Data) {
			if and {
				v &= right.Data[i]
			} else {
				v |= right.Data[i]
			}
		}
		ret.Data[i] = v
	}
	return ret
}

// BitAndºSetSet equals set & set
func BitAndºSetSet(left *core.Set, right *core.Set) *core.Set {
	return bitSet(left, right, true)
}

// BitOrºSetSet equals set & set
func BitOrºSetSet(left *core.Set, right *core.Set) *core.Set {
	return bitSet(left, right, false)
}

// BitNotºSet changes boolean value of set
func BitNotºSet(set *core.Set) *core.Set {
	ret := core.NewSet()
	ret.Data = make([]uint64, len(set.Data))
	for i, v := range set.Data {
		ret.Data[i] = ^v
	}
	return ret
}

// SetºSet sets the item in the set
func SetºSet(set *core.Set, index int64) (*core.Set, error) {
	var err error
	if err = checkIndex(set, index); err == nil {
		set.Set(index, true)
	}
	return set, err
}

// ToggleºSetInt changes the value of the set
func ToggleºSetInt(set *core.Set, index int64) (prev bool, err error) {
	if err = checkIndex(set, index); err == nil {
		prev = set.IsSet(index)
		set.Set(index, !prev)
	}
	return
}
*/
// UnSetºSet sets the item in the set
func UnSetºSet(set *core.Set, index int64) (*core.Set, error) {
	var err error
	if err = checkIndex(set, index); err == nil {
		set.Set(index, false)
	}
	return set, err
}

// setºStr converts string to set
func setºStr(value string) (*core.Set, error) {
	s := core.NewSet()
	for i, ch := range value {
		switch ch {
		case '0':
		case '1':
			s.Set(int64(i), true)
		default:
			return nil, fmt.Errorf(core.ErrorText(core.ErrInvalidParam))
		}
	}
	return s, nil
}

/*
// arrºSet converts set to array of integers
func arrºSet(set *core.Set) *core.Array {
	ret := core.NewArray()
	for i, v := range set.Data {
		for pos := uint64(0); pos < 64; pos++ {
			if v&(1<<pos) != 0 {
				ret.Data = append(ret.Data, int64(i<<6)+int64(pos))
			}
		}
	}
	return ret
}
*/
// setºArr converts array of integers to set
func setºArr(arr *core.Array) (set *core.Set, err error) {
	var ind int64
	set = core.NewSet()
	for _, v := range arr.Data {
		ind = v.(int64)
		if err = checkIndex(set, ind); err == nil {
			set.Set(ind, true)
		}
	}
	return
}

/*
// strºSet converts set to string
func strºSet(set *core.Set) string {
	return set.String()
}
*/

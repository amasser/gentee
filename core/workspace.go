// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Workspace contains information of compiled source code
type Workspace struct {
	Units     []*Unit
	UnitNames map[string]int
	Objects   []IObject
	Linked    map[string]int // compiled files
	IotaID    int32
	Embedded  []Embed
}

const (
	// DefName is the key name for stdlib
	DefName = `stdlib`

	// PubOne means the only next object is public
	PubOne = 1
	// PubAll means all objects are public
	PubAll = 2
)

// Unit is a common structure for source code
type Unit struct {
	VM        *Workspace
	Index     uint32            // Index of the Unit
	NameSpace map[string]uint32 // name space of the unit
	Included  map[uint32]bool   // false - included or true - imported units
	Lexeme    *Lex              // The array of source code
	RunID     int               // The index of run function. Undefined (-1) - run has not yet been defined
	Name      string            // The name of the unit
	Pub       int               // Public mode
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandName returns random latin name
func RandName() string {
	alpha := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := len(alpha)
	b := make([]rune, 16)
	for i := range b {
		b[i] = alpha[rand.Intn(length)]
	}
	return string(b)
}

// NewVM returns a new virtual machine
func NewVM(Embedded []Embed) *Workspace {
	ws := Workspace{
		UnitNames: make(map[string]int),
		Units:     make([]*Unit, 0, 32),
		Objects:   make([]IObject, 0, 500),
		Linked:    make(map[string]int),
		Embedded:  Embedded,
	}
	return &ws
}

// InitUnit initialize a unit structure
func (ws *Workspace) InitUnit() *Unit {
	return &Unit{
		VM:        ws,
		RunID:     Undefined,
		NameSpace: make(map[string]uint32),
		Included:  make(map[uint32]bool),
	}
}

// TypeByGoType returns the type by the go type name
func (unit *Unit) TypeByGoType(goType reflect.Type) *TypeObject {
	var name string
	switch goType.String() {
	case `int64`:
		name = `int`
	case `float64`:
		name = `float`
	case `bool`:
		name = `bool`
	case `string`:
		name = `str`
	case `int32`:
		name = `char`
	case `core.KeyValue`:
		name = `keyval`
	case `core.Range`:
		name = `range`
	case `*core.Buffer`:
		name = `buf`
	case `*core.Set`:
		name = `set`
	case `*core.Array`:
		name = `arr`
	case `*core.Map`:
		name = `map`
	case `*core.Obj`:
		name = `obj`
	default:
		return nil
	}
	if obj := unit.FindType(name); obj != nil {
		return obj.(*TypeObject)
	}
	return nil
}

// StdLib returns the pointer to Standard Library Unit
func (ws *Workspace) StdLib() *Unit {
	return ws.Unit(DefName)
}

// Unit returns the pointer to Unit by its name
func (ws *Workspace) Unit(name string) *Unit {
	return ws.Units[ws.UnitNames[name]]
}

func (unit *Unit) GetHeader(name string) string {
	for _, line := range strings.Split(unit.Lexeme.Header, "\n") {
		ret := regexp.MustCompile(`\s*` + strings.ReplaceAll(name, `.`, `\.`) +
			`\s*=\s*(.*)$`).FindStringSubmatch(strings.TrimSpace(line))
		if len(ret) == 2 {
			return ret[1]
		}
	}
	return ``
}

// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	ws "github.com/gentee/gentee"
	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/vm"
)

// Source contains source code and result value
type Source struct {
	Src  string
	Want string
	Line int
}

func loadTest(filename string) (src []Source, err error) {
	var input []byte
	src = make([]Source, 0, 64)
	input, err = ioutil.ReadFile(filepath.Join(`../tests`, filename))
	if err != nil {
		return
	}
	list := strings.Split(string(input), "\n")
	source := make([]string, 0, 32)
	on := true
	for i, line := range list {
		if on && strings.HasPrefix(line, `OFF`) {
			on = false
			continue
		}
		if !on {
			if strings.HasPrefix(line, `ON`) {
				on = true
			}
			continue
		}

		if !strings.HasPrefix(line, `=====`) {
			source = append(source, line)
			continue
		}
		src = append(src, Source{
			Src:  strings.Join(source, "\n"),
			Want: strings.TrimSpace(strings.TrimLeft(line, `=`)),
			Line: i,
		})
		source = source[:0]
	}
	return
}

func TestVM(t *testing.T) {
	workspace := ws.New()

	testFile := func(filename string) error {
		src, err := loadTest(filename)
		if err != nil {
			return err
		}
		for i := len(src) - 1; i >= 0; i-- {
			testErr := func(err error) error {
				return fmt.Errorf(`[%d] of %s  %v`, src[i].Line, filename, err)
			}
			fmt.Println(`LINE START`, src[i].Line, filename)
			unitID, err := workspace.Compile(src[i].Src, ``)
			if err != nil && err.Error() != src[i].Want {
				return testErr(err)
			}
			if err != nil {
				continue
			}
			linked, err := compiler.Link(workspace.Workspace, unitID)
			if err != nil {
				return testErr(err)
			}
			result, err := vm.Run(linked, vm.Settings{})
			if err == nil {
				if err = getWant(result, src[i].Want); err != nil {
					fmt.Println(`LINE`, src[i].Line)
					return testErr(err)
				}
			} else if err.Error() != src[i].Want {
				fmt.Println(`ERR`, err.Error(), `!=`, src[i].Want)
				return testErr(err)
			}
		}
		return nil
	}
	for _, name := range []string{`run_test`, `err_test`} {
		if err := testFile(name); err != nil {
			t.Error(err)
			return
		}
	}
	files, err := ioutil.ReadDir(filepath.Join("../tests", "stdlib"))
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) < 8 {
		t.Error(`stdlib tests cannot be found`)
		return
	}
	return
	for _, file := range files {
		if err := testFile(filepath.Join(`stdlib`, file.Name())); err != nil {
			t.Error(err)
			return
		}
	}
	if runtime.GOOS == `linux` {
		for _, name := range []string{`linux_test`} {
			if err := testFile(name); err != nil {
				t.Error(err)
				return
			}
		}
	}
	scriptName := filepath.Join(`../tests`, filepath.Join(`scripts`, `const.g`))
	unitID, err := workspace.CompileFile(scriptName)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := workspace.Run(unitID)
	if err != nil {
		t.Error(err)
		return
	}
	if result != workspace.Version() {
		t.Errorf(`Wrong version %v`, result)
		return
	}
	if !strings.HasSuffix(workspace.Unit(unitID).Name, scriptName) {
		t.Errorf(`Wrong unit name %v`, workspace.Unit(unitID).Name)
		return
	}
}

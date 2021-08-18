//
// Copyright 2021 Three Guys Labs, LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
package shell_test

import (
	"github.com/threeguys/golang-ezshell/shell"
	"github.com/threeguys/golang-toolkit/objects"
	"os"
	"strings"
	"testing"
)

func TestNewCommandShell(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cs := shell.NewCommandShell("foo", []*shell.Command{})
	assert.NotNil(cs)
	assert.Equal(cs.Out, os.Stdout)
	assert.False(cs.Echo)
	assert.Equal(cs.Global, cs.Mode)
	assert.Equal(0, len(cs.Global.Commands))
	assert.NotNil(cs.Global.Delegate)
	assert.Equal("foo", cs.Prompt)

	bazCmd := &shell.Command {
		Name:        "baz",
		Description: "a command",
		Handler:     shell.NoOpHandler(),
	}

	cs = shell.NewCommandShell("bar", []*shell.Command{ bazCmd })

	assert.NotNil(cs)
	assert.Equal(1, len(cs.Global.Commands))
	assert.Equal(bazCmd, cs.Global.Commands[0])
}

func TestShell_SwitchMode(t *testing.T) {

	bazCmd := &shell.Command {
		Name:        "baz",
		Description: "a command",
		Handler:     shell.NoOpHandler(),
	}

	batCmd := &shell.Command{
		Name:        "bat",
		Description: "a different command",
		Flags:       1,
		Handler:     shell.NoOpHandler(),
	}

	batMode := &shell.CommandMode{
		Name:        "bat-mode",
		Description: "nana-nana-nana bat mode!",
		Commands:    []*shell.Command { batCmd },
		Delegate:    nil,
	}

	assert := objects.NewTestAssertions(t)
	cs := shell.NewCommandShell("boh", []*shell.Command{ bazCmd }, batMode)
	assert.NotNil(cs)
	assert.Equal(1, len(cs.Global.Commands))
	assert.Equal(batMode, cs.Mode)
	assert.Equal(cs.Global, batMode.Delegate)

	assert.Nil(cs.SwitchMode("global"))
	assert.Equal(cs.Global, cs.Mode)
	assert.Nil(cs.SwitchMode("bat-mode"))
	assert.Equal(batMode, cs.Mode)
	assert.Equal(shell.ErrNoMatch, cs.SwitchMode("not found mode"))
}

func TestShell_Print(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cs := shell.NewCommandShell("bunk", []*shell.Command{})
	out := makeTempLog(t)
	assert.Nil(out.Close())
	cs.Out = out

	cs.Printf("This will not 'fail' per se, but there will be log message")
	cs.Println("what should happen?")
}

func TestNewCommandShell_Help(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cs := shell.NewCommandShell("bunk", []*shell.Command{})
	out := makeTempLog(t)
	defer func() { assert.Nil(out.Close()) }()
	cs.Out = out

	assert.Nil(cs.RunCommand([]string { "help" }))
	resetTempFile(t, out)
	assert.False(strings.EqualFold("", strings.TrimSpace(getLogData(t, out))))
}

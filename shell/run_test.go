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
	"bytes"
	"errors"
	"fmt"
	"github.com/threeguys/golang-ezshell/shell"
	"github.com/threeguys/golang-toolkit/objects"
	"io"
	"os"
	"strings"
	"testing"
)

var (
	mockErrMade = errors.New("i made an error")
)

func createMockTestShell() *shell.Shell {
	makeErr := &shell.Command{
		Name:        "err",
		Description: "makes an error",
		Flags:       13,
		Handler: func(strings []string) error {
			return mockErrMade
		},
	}

	okMode := &shell.CommandMode{
		Name:        "ok-mode",
		Description: "this only has ok commands",
		Commands:    []*shell.Command{
			{
				Name:        "noop",
				Description: "this command is ok",
				Flags:       42,
				Handler:     shell.NoOpHandler(),
			},
		},
		Delegate:    nil,
	}

	cs := shell.NewCommandShell("# ", []*shell.Command { makeErr }, okMode)
	cs.Echo = true
	return cs
}

func generateExpectedLog(t *testing.T, commands [][]string) string{
	assert := objects.NewTestAssertions(t)
	buffer := new(bytes.Buffer)

	responses := map[string]string {
		"err": fmt.Sprintf("ERROR: %s", mockErrMade.Error()),
		"noop": "SUCCESS",
	}

	cmdPrinter := func(cmd string) {
		resp, ok := responses[cmd]
		if !ok {
			resp = fmt.Sprintf("ERROR: %s", shell.ErrNoMatch)
		}
		_, err := fmt.Fprintf(buffer, "%s\n%s\n# ", cmd, resp)
		assert.Nil(err)
	}

	_, err := fmt.Fprintf(buffer, "# ")
	assert.Nil(err)

	for _, cmd := range commands {
		cmdPrinter(strings.Join(cmd, " "))
	}

	return buffer.String()
}

func TestShell_RunSupplier(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cmdList := [][]string{ { "err" }, { "noop" }, { "command not found" } }
	supplier := shell.NewListCommandSupplier(cmdList...)

	cs := createMockTestShell()
	cs.Out = makeTempLog(t)
	defer func() { assert.Nil(cs.Out.Close()) }()

	assert.Equal(io.EOF, cs.RunSupplier(supplier))

	expected := generateExpectedLog(t, cmdList)
	logs := getLogData(t, cs.Out)
	assert.Equal(expected, logs)
}

func TestShell_RunFile(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cmdList := [][]string { { "noop" }, { "noop" }, { "err" } }

	in := makeTempLog(t)
	defer func() { assert.Nil(in.Close()) }()

	for _, cl := range cmdList {
		_, err := fmt.Fprintln(in, cl[0])
		assert.Nil(err)
	}
	resetTempFile(t, in)

	cs := createMockTestShell()
	cs.Out = makeTempLog(t)
	defer func() { assert.Nil(cs.Out.Close()) }()

	assert.Equal(io.EOF, cs.RunFile(in))

	expected := generateExpectedLog(t, cmdList)
	logs := getLogData(t, cs.Out)
	assert.Equal(expected, logs)
}

func TestShell_Run(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cmdList := [][]string { { "noop" } }

	in := makeTempLog(t)
	defer func() { assert.Nil(in.Close()) }()

	for _, cl := range cmdList {
		_, err := fmt.Fprintln(in, cl[0])
		assert.Nil(err)
	}
	resetTempFile(t, in)

	cs := createMockTestShell()
	cs.Out = makeTempLog(t)
	defer func() { assert.Nil(cs.Out.Close()) }()

	oldStdin := os.Stdin
	os.Stdin = in
	err := cs.Run()
	os.Stdin = oldStdin
	assert.Equal(io.EOF, err)

	expected := generateExpectedLog(t, cmdList)
	logs := getLogData(t, cs.Out)
	assert.Equal(expected, logs)
}

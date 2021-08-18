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
	"errors"
	"fmt"
	"github.com/threeguys/golang-ezshell/shell"
	"github.com/threeguys/golang-toolkit/objects"
	"testing"
)

func TestHandlerWrap(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	called := false
	wrapped := shell.HandlerWrap(func() { called = true })
	assert.Nil(wrapped([]string { "foo" }))
	assert.True(called)
}

func TestNoOpHandler(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	handler := shell.NoOpHandler()
	assert.Nil(handler(nil))
	assert.Nil(handler([]string{}))
	assert.Nil(handler([]string { "a word" }))
}

func createTestNameList(count int) []string {
	names := make([]string, 0)
	for i := 0; i < count; i++ {
		names = append(names, fmt.Sprintf("cmd-%d", i))
	}
	return names
}

func createTestCommandList(count int) []*shell.Command {
	commands := make([]*shell.Command, 0)
	for i, name := range createTestNameList(count) {
		commands = append(commands, &shell.Command{
			Name:        name,
			Description: fmt.Sprintf("desc-%d", i),
			Flags:       uint32(i),
			Handler:     func(_ []string) error { return errors.New(fmt.Sprintf("err-%d", i)) },
		})
	}
	return commands
}

func createTestCommandMode(count int) *shell.CommandMode {
	return &shell.CommandMode{
		Name:        "test",
		Description: "my test mode",
		Commands:    createTestCommandList(count),
		Delegate:    nil,
	}
}

func TestCommandMode_Match_NoDelegate(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cm := createTestCommandMode(5)

	for i, name := range createTestNameList(5) {
		cmd, err := cm.Match(name)
		assert.Equal(cm.Commands[i], cmd)
		assert.Nil(err)
	}

	cmd, err := cm.Match("cmd-not-found")
	assert.Nil(cmd)
	assert.Equal(err, shell.ErrNoMatch)
}

func TestCommandMode_Match_WithDelegate(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cm := createTestCommandMode(1)
	cm.Delegate = &shell.CommandMode{
		Name:        "the-delegate",
		Description: "delegate desc",
		Commands:    []*shell.Command{
			{
				Name:        "delegated",
				Description: "a delegated command",
				Flags:       17,
				Handler:     shell.NoOpHandler(),
			},
		},
		Delegate:    nil,
	}

	checkCmd := func(cmdName string, cmd *shell.Command) {
		found, err := cm.Match(cmdName)
		assert.Nil(err)
		assert.Equal(cmd, found)
	}

	checkCmd("cmd-0", cm.Commands[0])
	checkCmd("delegated", cm.Delegate.Commands[0])
	_, err := cm.Match("not-gonna-be-there")
	assert.Equal(shell.ErrNoMatch, err)
}

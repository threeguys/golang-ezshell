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
	"github.com/threeguys/golang-ezshell/shell"
	"github.com/threeguys/golang-toolkit/objects"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func assertRegexp(assert *objects.Assertions, re *regexp.Regexp, value string) {
	assert.True(re.MatchString(value))
}

func TestNewHelpMode(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	var args []string
	h := shell.NewHelpMode(func(theArgs []string) error {
		args = theArgs
		return errors.New("another test error")
	})

	assert.NotNil(h)
	assert.Equal("help", h.Name)
	assert.True(len(h.Description) > 0)
	assert.NotNil(h.Commands)
	assert.Equal(1, len(h.Commands))
	assert.Nil(h.Delegate)

	cmd := h.Commands[0]
	assert.Equal("help", cmd.Name)
	assert.True(len(cmd.Description) > 0)
	assert.Equal(uint32(0), cmd.Flags)
	assert.NotNil(cmd.Handler)

	assert.Equal(errors.New("another test error"), cmd.Run([]string { "foo", "bar" }))
	assert.Equal([]string { "foo", "bar" }, args)
}

func makeTempLog(t *testing.T) *os.File {
	out, err := ioutil.TempFile(t.TempDir(), "help-out")
	if err != nil {
		t.Fatal("Could not create temp file", err)
	}
	return out
}

func resetTempFile(t *testing.T, f *os.File) {
	assert := objects.NewTestAssertions(t)
	pos, err := f.Seek(0, 0)
	if err != nil {
		t.Fatal("Could not reset temp file", err)
	}
	assert.Equal(int64(0), pos)
}

func getLogData(t *testing.T, out *os.File) string {
	assert := objects.NewTestAssertions(t)
	resetTempFile(t, out)

	data, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatal("Could not read temp file", err)
	}
	assert.Nil(err)
	return string(data)
}

func TestShell_PrintHelp_NoModes(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	cmd1 := &shell.Command{
		Name:        "testcmd",
		Description: "testcmddesc",
		Flags:       3,
		Handler:     shell.NoOpHandler(),
	}
	cs := shell.NewCommandShell("test1", []*shell.Command{cmd1 })
	out := makeTempLog(t)
	defer func(){ assert.Nil(out.Close()) }()

	cs.Out = out
	cs.PrintHelp()

	expr := `(?s:.*global.* testcmd [^\n]+ testcmddesc[ \n].*)`
	assertRegexp(assert, regexp.MustCompile(expr), getLogData(t, out))
}

func TestShell_PrintHelp_WithModes(t *testing.T) {
	assert := objects.NewTestAssertions(t)
	myErr := errors.New("this-is-a-test")
	cmd1 := &shell.Command{
		Name:        "testcmd",
		Description: "testcmddesc",
		Flags:       3,
		Handler:     shell.NoOpHandler(),
	}
	cmd2 := &shell.Command{
		Name:        "modecmd",
		Description: "fake-mode-desc",
		Flags:       42,
		Handler:     func(_ []string) error { return myErr },
	}
	mode := &shell.CommandMode{
		Name:        "testmode",
		Description: "test-mode-comment",
		Commands:    []*shell.Command{cmd2 },
		Delegate:    nil,
	}
	cs := shell.NewCommandShell("test1", []*shell.Command{cmd1 }, mode)
	out := makeTempLog(t)
	defer func(){ assert.Nil(out.Close()) }()
	cs.Out = out

	cs.PrintHelp()

	expr := `(?s:.*global.* testcmd [^\n]+ testcmddesc[ \n].*testmode.*test-mode-comment.* modecmd [^\n]+ fake-mode-desc[ \n].*)`
	assertRegexp(assert, regexp.MustCompile(expr), getLogData(t, out))
}

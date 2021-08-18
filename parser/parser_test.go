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
package parser_test

import (
	"errors"
	"fmt"
	"github.com/threeguys/golang-ezshell/parser"
	"github.com/threeguys/golang-toolkit/objects"
	"io"
	"strings"
	"testing"
)

func TestChangeState(t *testing.T) {
	// Test bad start state
	state, capturing, err := parser.ChangeState(parser.StateEOF, 0, 0)
	assert := objects.NewTestAssertions(t)
	assert.Equal(parser.StateParseError, state)
	assert.False(capturing)
	assert.NotNil(err)
	assert.Equal(fmt.Sprintf("unknown parser state (%d)", parser.StateEOF), err.Error())
}

type mockBadReader struct {
	Data []byte
	Error error
}

func (mbr *mockBadReader) Read(data []byte) (int, error) {
	if mbr.Data != nil {
		return copy(data, mbr.Data), mbr.Error
	} else {
		return 0, mbr.Error
	}
}

func TestCommandReader_Read_BadReader(t *testing.T) {
	mbr := &mockBadReader{
		Data:  []byte("abc"),
		Error: errors.New("a test error"),
	}

	assert := objects.NewTestAssertions(t)
	cr := parser.NewCommandReader(mbr)
	cmd, err := cr.Read()
	assert.Nil(cmd)
	assert.NotNil(err)
	assert.Equal("a test error", err.Error())
}

func TestCommandReader_Read(t *testing.T) {

	type commandReaderTest struct {
		Input string
		Output []string
		Error error
	}

	tests := []*commandReaderTest{
		{
			Input: `"ab cd" not in 'quotes now'`,
			Output: []string{ "ab cd", "not", "in", "quotes now"},
			Error: io.EOF,
		},
		{
			Input: `"abcd \" e" some more`,
			Output: []string{ "abcd \" e", "some", "more"},
			Error: io.EOF,
		},
		{
			Input: `"here's" something with 'single \' quotes'`,
			Output: []string{ "here's", "something", "with", "single ' quotes"},
			Error: io.EOF,
		},
		{
			Input: `'start with' single quote`,
			Output: []string{ "start with", "single", "quote" },
			Error: io.EOF,
		},
		{
			Input: "  leading \tspaces",
			Output: []string{ "leading", "spaces" },
			Error: io.EOF,
		},
		{ Input: "", Output: []string{}, Error: io.EOF },
		{ Input: "   ", Output: []string{}, Error: io.EOF },
		{ Input: "\t ", Output: []string{}, Error: io.EOF },
		{ Input: " \t ", Output: []string{}, Error: io.EOF },
		{ Input: ` "" '' `, Output: []string{}, Error: io.EOF },
		{ Input: `"foo"bar`, Output: nil, Error: errors.New("expected ' ' at char 5") },
		{   Input: `can't have a quote in the middle of a word`,
			Output: nil,
			Error: errors.New("unexpected ['] at char 3"),
		},
	}

	assert := objects.NewTestAssertions(t)
	checkTest := func(test *commandReaderTest, cmds []string, err error) {
		fmt.Println("Checking [", test.Input, "]")
		if test.Error == nil {
			assert.Nil(err)
		} else {
			assert.Equal(test.Error, err)
		}
		if test.Output == nil {
			assert.Nil(cmds)
		} else {
			assert.Equal(test.Output, cmds)
		}
	}

	allCmds := make([]string, 0, len(tests))
	for _, test := range tests {
		allCmds = append(allCmds, test.Input)
		cr := parser.NewCommandReader(strings.NewReader(test.Input))

		cmds, err := cr.Read()
		checkTest(test, cmds, err)
	}
	allCmds = append(allCmds, "\n")

	// Run all together as multiple commands
	cr := parser.NewCommandReader(strings.NewReader(strings.Join(allCmds, "\n")))
	for _, test := range tests {
		cmds, err := cr.Read()
		if test.Error == io.EOF {
			test.Error = nil
		}
		checkTest(test, cmds, err)
	}

}

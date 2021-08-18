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
package shell

import (
	"strings"
)

type CommandMode struct {
	Name string
	Description string
	Commands []*Command
	Delegate *CommandMode
}

func HandlerWrap(op func()) CommandHandler {
	return func(_ []string) error {
		op()
		return nil
	}
}

func NoOpHandler() CommandHandler {
	return func(_ []string) error { return nil }
}

func newGlobalMode(help CommandHandler, cmds []*Command) *CommandMode {
	return &CommandMode{
		Name:        "global",
		Description: "Available commands",
		Commands:    cmds,
		Delegate:    NewHelpMode(help),
	}
}

func (cm *CommandMode) Match(cmd string) (*Command, error) {
	for _, c := range cm.Commands {
		if strings.Compare(c.Name, cmd) == 0 {
			return c, nil
		}
	}
	if cm.Delegate != nil {
		return cm.Delegate.Match(cmd)
	} else {
		return nil, ErrNoMatch
	}
}

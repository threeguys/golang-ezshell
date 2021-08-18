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
	"fmt"
	"log"
	"os"
	"strings"
)

type Printer interface {
	Println(out ... interface{})
	Printf(fmtStr string, vars ... interface{})
}

type Shell struct {
	modes []*CommandMode
	modeIndex map[string]*CommandMode
	Mode *CommandMode
	Global *CommandMode
	Prompt string
	Out *os.File
	Echo bool
	Quiet bool
}

func NewCommandShell(prompt string, global []*Command, cmd ... *CommandMode) *Shell {
	var cs *Shell
	globalMode := newGlobalMode(HandlerWrap(func() { cs.PrintHelp() }), global)

	for _, c := range cmd {
		c.Delegate = globalMode
	}

	defaultMode := globalMode
	if len(cmd) > 0 {
		defaultMode = cmd[0]
	}

	modeIndex := make(map[string]*CommandMode)
	modeIndex[globalMode.Name] = globalMode
	for _, cm := range cmd {
		modeIndex[cm.Name] = cm
	}

	cs = &Shell{
		modes:     cmd,
		modeIndex: modeIndex,
		Mode:      defaultMode,
		Global:    globalMode,
		Prompt:    prompt,
		Out:       os.Stdout,
		Echo:      false,
		Quiet:     false,
	}

	return cs
}

func (cs *Shell) Println(out ... interface{}) {
	if _, err := fmt.Fprintln(cs.Out, out...); err != nil {
		log.Println("Unable to write to output file", err)
	}
}

func (cs *Shell) Printf(fmtStr string, vars ... interface{}) {
	if _, err := fmt.Fprintf(cs.Out, fmtStr, vars...); err != nil {
		log.Println("Unable to write to output file", err)
	}
}

func (cs *Shell) SwitchMode(name string) error {
	if mode, ok := cs.modeIndex[name]; !ok {
		return ErrNoMatch
	} else {
		cs.Mode = mode
		return nil
	}
}

func (cs *Shell) RunCommand(parsed []string) error {
	if cs.Echo {
		cs.Printf("%s\n", strings.Join(parsed, " "))
	}
	if cmd, err := cs.Mode.Match(parsed[0]); err != nil {
		return err
	} else {
		return cmd.Run(parsed[1:])
	}
}

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
package main

import (
	"errors"
	"fmt"
	"github.com/threeguys/golang-ezshell/shell"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// EzBash is a simple, contrived example to illustrate how to use
// the golang-ezshell library. It consists of a couple of "modes" to
// show how to use those.
type EzBash struct {
	*shell.Shell
}

// Constructor for the EzBash shell
func NewEzBash() *EzBash {
	ezb := &EzBash{}

	// The list of "global" commands which are available
	// in all modes of the shell
	commands := []*shell.Command{
		{
			Name:        "ls",
			Description: "lists the files in the current directory",
			Flags:       shell.FlagOptionalArgs,
			Handler:     ezb.HandlerList,
		},
		{
			Name:        "pwd",
			Description: "prints the current directory",
			Flags:       0,
			Handler:     ezb.HandlerPwd,
		},
		{
			Name:        "cd",
			Description: "changes the current directory",
			Flags:       shell.FlagRequiresArgs,
			Handler:     func (args []string) error { return os.Chdir(args[0]) },
		},
		{
			Name:        "mode",
			Description: "set the mode to user/debug/admin",
			Flags:       shell.FlagOptionalArgs,
			Handler:     ezb.HandlerMode,
		},
		{
			Name:        "exit",
			Description: "exits the shell",
			Flags:       0,
			Handler:     ezb.HandlerExit,
		},
	}

	// The "user" mode doesn't have any special commands, it's just a way
	// to differentiate the debug and admin modes
	userMode := &shell.CommandMode{
		Name:        "user",
		Description: "allows changing modes",
		Commands:    []*shell.Command{},
	}

	// Debug mode allows you to get/set environment variables
	debugMode := &shell.CommandMode{
		Name:        "debug",
		Description: "allows getting/setting of environment variables",
		Commands:    []*shell.Command {
			{
				Name:        "get",
				Description: "gets a variable value",
				Flags:       shell.FlagRequiresArgs,
				Handler:     ezb.HandlerGet,
			},
			{
				Name:        "set",
				Description: "sets a variable value",
				Flags:       shell.FlagRequiresArgs,
				Handler:     ezb.HandlerSet,
			},
		},
	}

	// Admin mode allows you to dump all of the variables (but you can
	// only get/set in debug mode)
	adminMode := &shell.CommandMode{
		Name:        "admin",
		Description: "allows seeing all environment variable values",
		Commands:    []*shell.Command {
			{
				Name:        "dump",
				Description: "display all environment variables",
				Flags:       0,
				Handler:     ezb.HandlerDump,
			},
		},
	}

	ezb.Shell = shell.NewCommandShell("ezbash $ ", commands, userMode, debugMode, adminMode)
	return ezb
}

// Simple helper function to make returning "error" a one-liner for the other handlers
func (ezb *EzBash) ConsolePrintf(format string, args ... interface{}) error {
	ezb.Printf(format, args...)
	return nil
}

// Implements the "ls" command, an argument is optional. If none is given
// then the current directory is listed, otherwise the specified file/directory
// will be listed.
func (ezb *EzBash) HandlerList(args []string) error {
	var loc string
	if len(args) > 0 {
		loc = args[0]
	} else if dir, err := os.Getwd(); err != nil {
		return err
	} else {
		loc = dir
	}

	return filepath.Walk(loc, func (path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return ezb.ConsolePrintf("%s/\n", info.Name())
			} else {
				return ezb.ConsolePrintf("%s\n", info.Name())
			}
		})
}

// Handler to show the current working directory
func (ezb *EzBash) HandlerPwd(_ []string) error {
	if dir, err := os.Getwd(); err != nil {
		return err
	}  else {
		return ezb.ConsolePrintf("%s\n", dir)
	}
}

// Handler to get environment variables, only available in debug mode
func (ezb *EzBash) HandlerGet(args []string) error {
	return ezb.ConsolePrintf("%s\n", os.Getenv(args[0]))
}

// Handler to set environment variables, only available in debug mode
func (ezb *EzBash) HandlerSet(args []string) error {
	if re, err := regexp.Compile(`^([^=]+)=(.*)$`); err != nil {
		return err
	} else if matches := re.FindStringSubmatch(args[0]); matches == nil || len(matches) != 3 {
		return errors.New(fmt.Sprintf("could parse set statement [%s] [%v]", args[0], matches))
	} else {
		return ezb.ConsolePrintf("SETTING [%s] = [%s]\n", matches[1], matches[2])
	}
}

// Handler to implement setting the mode. The modes in this
// case just control which commands are available outside of the
// normal global mode
func (ezb *EzBash) HandlerMode(args []string) error {
	if len(args) == 0 {
		return ezb.ConsolePrintf("%s\n", ezb.Mode.Name)
	} else {
		mode := strings.ToLower(args[0])
		switch mode {
		case "user","debug","admin":
			return ezb.SwitchMode(mode)
		default:
			return errors.New(fmt.Sprintf("unknown mode [%s]", mode))
		}
	}
}

// This handler implements the "dump" command, which prints all
// of the environment variables out and is only available in "admin" mode
func (ezb *EzBash) HandlerDump(_ []string) error {
	for _, e := range os.Environ() {
		ezb.Printf("%s\n", e)
	}
	return nil
}

// This handler implements the "exit" command, which just
// quits the shell
func (ezb *EzBash) HandlerExit(_ []string) error {
	os.Exit(0)
	return nil
}

// Main entrypoint for the EzBash shell
func (ezb *EzBash) ShellMain() {
	if len(os.Args) > 1 {
		if f, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0); err != nil {
			log.Fatal(err)
		} else if err := ezb.RunFile(f); err != nil && err != io.EOF {
			log.Fatal(err)
		}
	} else {
		if err := ezb.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	ezb := NewEzBash()
	ezb.ShellMain()
}

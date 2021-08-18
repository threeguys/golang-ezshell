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
	"github.com/threeguys/golang-ezshell/shell"
	"log"
	"os"
)

func nullHelper(op func()) shell.CommandHandler {
	return func(_ []string) error {
		op()
		return nil
	}
}

func argHelper(op func(args []string)) shell.CommandHandler {
	return func(args []string) error {
		op(args)
		return nil
	}
}

func main() {
	var sh *shell.Shell
	sh = shell.NewCommandShell("$ ", []*shell.Command{
		{
			Name:        "hello",
			Description: "say hello",
			Handler:     argHelper(func(args []string) {
				if len(args) > 0 {
					sh.Printf("Hello, %s. How are you?\n", args[0])
				} else {
					sh.Printf("Hello! I am ezbash, type 'help' to see what I can do")
				}
			}),
		},
		{
			Name:        "bye",
			Description: "quits the shell",
			Flags:       0,
			Handler:     nullHelper(func() { os.Exit(0) }),
		},
		{
			Name:        "move",
			Description: "moves in a certain direction",
			Flags:       0,
			Handler:     argHelper(func(args []string) {
				if len(args) > 0 {
					sh.Printf("You move %s\n", args[0])
				} else {
					sh.Println("You move in a random direction")
				}
				sh.Println("You are eaten by a grue.\nPlease try to be more careful")
			}),
		},
	})
	log.Fatal(sh.Run())
}

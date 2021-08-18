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

func NewHelpMode(help CommandHandler) *CommandMode {
	return &CommandMode{
		Name:        "help",
		Description: "When you don't know how to do something",
		Commands:    []*Command {
			{
				Name:        "help",
				Description: "Display this message",
				Handler:     help,
			},
		},
	}
}

func (cs *Shell) helpHelper(mode *CommandMode) {
	cs.Printf("  [ (%s) :: %s ]\n    ---> Commands <---\n", mode.Name, mode.Description)
	for _, c := range mode.Commands {
		cs.Printf("      %s - %s\n", c.Name, c.Description)
	}
	cs.Println()
}

func (cs *Shell) PrintHelp() {
	cs.Println()
	cs.helpHelper(cs.Global)
	if len(cs.modes) > 0 {
		cs.Printf( "  Current mode: %s\n\n  == Modes ==\n\n", cs.Mode.Name)
		for _, m := range cs.modes {
			cs.helpHelper(m)
		}
	}
}

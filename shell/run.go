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
	"github.com/threeguys/golang-ezshell/parser"
	"io"
	"os"
)

type CommandSupplier interface {
	Read() ([]string, error)
}

type ListCommandSupplier struct {
	records [][]string
	index int
}

func NewListCommandSupplier(recs ... []string) CommandSupplier {
	return &ListCommandSupplier{
		records: recs,
		index:   0,
	}
}

func (lcs *ListCommandSupplier) Read() ([]string, error) {
	if lcs.index >= len(lcs.records) {
		return nil, io.EOF
	} else {
		rec := lcs.records[lcs.index]
		lcs.index++
		return rec, nil
	}
}

func (cs *Shell) RunSupplier(rdr CommandSupplier) error {
	for {
		if !cs.Quiet {
			cs.Printf(cs.Prompt)
		}
		if parsed, err := rdr.Read(); err != nil {
			return err
		} else if err := cs.RunCommand(parsed); err != nil {
			cs.Println("ERROR:", err)
		} else if !cs.Quiet {
			cs.Println("SUCCESS")
		}
	}
}

func (cs *Shell) RunFile(f *os.File) error {
	rdr := parser.NewCommandReader(f)
	return cs.RunSupplier(rdr)
}

func (cs *Shell) Run() error {
	return cs.RunFile(os.Stdin)
}

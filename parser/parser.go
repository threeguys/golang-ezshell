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
package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const (
	StateReading = iota
	StateInWord
	StateDblQuote
	StateInDblQuote
	StateEndDblQuote
	StateSglQuote
	StateInSglQuote
	StateEndSglQuote
	StateLineFeed
	StateDblEscape
	StateSglEscape
	StateEndWord
	StateEOF
	StateParseError
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(rdr io.Reader) *CommandReader {
	return &CommandReader{ bufio.NewReader(rdr) }
}

func (cr *CommandReader) lineFeedOrEOF() error {
	counter := 0
	for {
		if c, err := cr.reader.ReadByte(); err != nil {
			return err
		} else if c == '\n' {
			return nil
		} else {
			counter++
		}
	}
}

func ChangeState(state int, c byte, index int) (int, bool, error) {
	if c == '\n' {
		return StateLineFeed, false, nil
	}

	switch state {
	case StateReading, StateEndWord, StateEndSglQuote, StateEndDblQuote:
		switch c {
		case '"':
			return StateDblQuote, false, nil
		case '\'':
			return StateSglQuote, false, nil
		case ' ','\t':
			return StateReading, false, nil
		default:
			if state == StateEndSglQuote || state == StateEndDblQuote {
				return StateParseError, false, errors.New(fmt.Sprintf("expected ' ' at char %d", index))
			}
			return StateInWord, true, nil
		}

	case StateDblQuote:
		if c != '"' {
			return StateInDblQuote, true, nil
		} else {
			return StateEndDblQuote, false, nil
		}

	case StateSglQuote:
		if c != '\'' {
			return StateInSglQuote, true, nil
		} else {
			return StateEndSglQuote, false, nil
		}

	case StateInWord:
		switch c {
		case '"','\'':
			return StateParseError, false, errors.New(fmt.Sprintf("unexpected [%c] at char %d", c, index))
		case ' ','\t':
			return StateEndWord, false, nil
		default:
			return StateInWord, true, nil
		}

	case StateInDblQuote:
		switch c {
		case '"':
			return StateEndDblQuote, false, nil
		case '\\':
			return StateDblEscape, false, nil
		default:
			return StateInDblQuote, true, nil
		}

	case StateDblEscape:
		return StateInDblQuote, true, nil

	case StateInSglQuote:
		switch c {
		case '\'':
			return StateEndSglQuote, false, nil

		case '\\':
			return StateSglEscape, false, nil

		default:
			return StateInSglQuote, true, nil
		}

	case StateSglEscape:
		return StateInSglQuote, true, nil

	default:
		return StateParseError, false, errors.New(fmt.Sprintf("unknown parser state (%d)", state))
	}
}

func (cr *CommandReader) Read() ([]string, error) {
	words := make([]string, 0)
	current := make([]byte, 0)
	capturing := false
	state := StateReading
	index := 0
	var parseErr error

	for {
		if c, err := cr.reader.ReadByte(); err != nil && err != io.EOF {
			return nil, err
		} else {

			if err == io.EOF {
				state = StateEOF

			} else if state, capturing, parseErr = ChangeState(state, c, index); parseErr != nil {
				_ = cr.lineFeedOrEOF()
				return nil, parseErr

			} else if capturing {
				current = append(current, c)
			}

			switch state {
			case StateLineFeed, StateEndWord, StateEndSglQuote, StateEndDblQuote, StateEOF:
				if len(current) > 0 {
					words = append(words, string(current))
					current = make([]byte, 0)
				}

				if state == StateLineFeed {
					return words, nil
				} else if state == StateEOF {
					return words, io.EOF
				}
			}

			index++
		}
	}
}

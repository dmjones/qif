//   Copyright 2018 Duncan Jones
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package qif

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
)

// Config defines the configuration of the reader.
type Config struct {

	// MonthFirst specifies whether to interpret dates as mm/dd/yy... or dd/mm/yy...
	MonthFirst bool
}

// DefaultConfig returns the default configuration used by NewReader:
//
//  Config{
//    MonthFirst: true,
//  }
func DefaultConfig() Config {
	return Config{
		MonthFirst: true,
	}
}

// A Reader consumes QIF data and returns parsed transactions.
type Reader interface {

	// Read returns the next transaction from the input data.
	Read() (Transaction, error)

	// ReadAll returns all the remaining transactions from the input data.
	ReadAll() ([]Transaction, error)
}

type reader struct {
	in           *bufio.Scanner
	config       Config
	headerParsed bool
}

// NewReader creates a new Reader with a default configuration (see DefaultConfig).
func NewReader(r io.Reader) *reader {
	return NewReaderWithConfig(r, DefaultConfig())
}

// NewReaderWithConfig creates a new Reader with the specified configuration.
func NewReaderWithConfig(r io.Reader, config Config) *reader {
	return &reader{
		in:     bufio.NewScanner(r),
		config: config,
	}
}

func (r *reader) parseHeader() error {
	if !r.in.Scan() {
		if err := r.in.Err(); err != nil {
			return err
		}

		return errors.New("file header not found")
	}

	switch r.in.Text() {
	case "!Type:Bank", "!Type:Cash", "!Type:CCard":
		r.headerParsed = true
		return nil

	default:
		return errors.Errorf("unsupported header type '%s'", r.in.Text())
	}
}

func (r *reader) Read() (Transaction, error) {

	if !r.headerParsed {
		err := r.parseHeader()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse file header")
		}
	}

	for r.in.Scan() {

	}

	return nil, nil
}

func (r *reader) ReadAll() ([]Transaction, error) {
	return nil, nil
}

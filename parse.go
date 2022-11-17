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

const (
	bankHeader = "!Type:Bank"
	cashHeader = "!Type:Cash"
	cardHeader = "!Type:CCard"
	recordEnd  = "^"
)

// A Reader consumes QIF data and returns parsed transactions.
type Reader interface {

	// Read returns the next transaction from the input data. Returns nil if
	// the end of the input has been reached. If the input ends without a
	// terminating '^' symbol, the result will be the transaction data read
	// thus far and a RecordEndError.
	Read() (Transaction, error)

	// ReadAll returns all the remaining transactions from the input data. It
	// returns the same errors as Read.
	ReadAll() ([]Transaction, error)

	// Informs what is class of transactions returned by Read() or ReadAll(),
	// so that transactions can be typecast to BankingTransaction or to
	// InvestmentTransaction.
	ReadTransactionType() TransactionType
}

// reader implements Reader. Construct using NewReader or NewReaderWithConfig.
type reader struct {

	// in scans the input.
	in *bufio.Scanner

	// config defines the behaviour of the reader.
	config Config

	// headerParsed is true if the header line has been read from the input
	// data.
	headerParsed bool

	// type of Transactions to be parsed (determined during ParseHeader)
	transactionType TransactionType
}

// NewReader creates a new Reader with a default configuration (see
// DefaultConfig).
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

// parseHeader reads the first line of the input and validates the header. An
// error is returned if the input is empty or the wrong type of header is found.
func (r *reader) parseHeader() error {
	if !r.in.Scan() {
		if err := r.in.Err(); err != nil {
			return err
		}

		return errors.New("file header not found")
	}

	switch r.in.Text() {
	case bankHeader, cashHeader, cardHeader:
		r.transactionType = TransactionTypeBanking
		r.headerParsed = true
		return nil

	default:
		return errors.Errorf("unsupported header type '%s'", r.in.Text())
	}
}

// Read implements Reader.Read.
func (r *reader) Read() (Transaction, error) {
	var tx Transaction

	if !r.headerParsed {
		err := r.parseHeader()
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse file header")
		}
	}

	// Parse type based on r.transactionType (set by parseHeader)
	// Only one type supported at the moment
	switch r.transactionType {
	default:
		tx = &bankingTransaction{}
	}
	data := false

	for r.in.Scan() {
		data = true
		line := r.in.Text()

		if line == recordEnd {
			return tx, nil
		}

		err := tx.parseTransactionTypeField(r.in.Text(), r.config)
		if err != nil {
			return nil, err
		}
	}

	if !data {
		// We were at the end of the file
		return nil, nil
	}

	return nil, RecordEndError{Incomplete: tx}
}

// ReadAll implements Reader.ReadAll.
func (r *reader) ReadAll() ([]Transaction, error) {
	var result []Transaction

	for {
		tx, err := r.Read()
		if err != nil {
			return nil, err
		}

		if tx == nil {
			break
		}

		result = append(result, tx)
	}

	return result, nil
}

func (r *reader) ReadTransactionType() TransactionType {
	return r.transactionType
}

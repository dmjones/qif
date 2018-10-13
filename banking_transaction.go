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

import "github.com/pkg/errors"

type AccountType int

const (
	Cash       AccountType = iota
	Bank                   = iota
	CreditCard             = iota
	//Investment = iota
	//Asset = iota
	//Liability = iota
	//Invoice = iota
)

// A BankingTransaction contains the information associated with non-investment
// transactions (i.e. Cash, Bank and CCard account types).
type BankingTransaction interface {
	Transaction

	// Num contains the check or reference number for the transaction. Wikipedia
	// suggests this may also contain "Deposit", "Transfer", "Print", "ATM", or
	// "EFT".
	Num() string

	// Payee describes the recipient of the transaction.
	Payee() string

	// Address contains no more than five address lines for the payee. Wikipedia
	// suggests the first entry is usually the same as the Payee field.
	Address() []string

	// AddressMessage contains an additional message associated with the payee
	// address. This is only non-empty if the transaction address contained a
	// special sixth line.
	AddressMessage() string

	// Category of the transaction.
	Category() string

	// Splits contains zero or more fragments of the transaction (AFAIK).
	Splits() []Split
}

type bankingTransaction struct {
	transaction
	num            string
	payee          string
	address        []string
	addressMessage string
	category       string
	splits         []Split
}

func (t *bankingTransaction) Num() string {
	return t.num
}

func (t *bankingTransaction) Payee() string {
	return t.payee
}

func (t *bankingTransaction) Address() []string {
	return t.address
}

func (t *bankingTransaction) AddressMessage() string {
	return t.addressMessage
}

func (t *bankingTransaction) Category() string {
	return t.category
}

func (t *bankingTransaction) Splits() []Split {
	return t.splits
}

func (t *bankingTransaction) parseBankingTransactionField(line string,
	config Config) error {
	if line == "" {
		return errors.New("line is empty")
	}

	err := t.parseTransactionField(line, config)
	if err == nil {
		// Must have been a field from our embedded struct
		return nil
	}

	if _, ok := err.(UnsupportedFieldError); !ok {
		// An actual error happened
		return err
	}

	// Otherwise, try and parse it here

	switch line[0] {
	case 'N':
		t.num = line[1:]
		return nil
	case 'P':
		t.payee = line[1:]
		return nil
	case 'A':
		// How many do we have already?
		if len(t.address) >= 5 {
			t.addressMessage = line[1:]
		} else {
			t.address = append(t.address, line[1:])
		}
		return nil
	case 'L':
		t.category = line[1:]
		return nil

		// These split fields must be in order, based on statement "The
		// non-split items can be in any sequence" from the spec.

	case 'S': // Category
		split := Split{}
		cat := line[1:]
		split.Category = &cat
		t.splits = append(t.splits, split)
		return nil
	case 'E': // Memo
		// This could be the first element of a new split, but only if there
		// isn't an existing split, or the existing split already has an 'E' or
		// a '$' field.
		if len(t.splits) == 0 || t.splits[len(t.splits)-1].Memo != nil ||
			t.splits[len(t.splits)-1].Amount != nil {
			t.splits = append(t.splits, Split{})
		}

		memo := line[1:]
		t.splits[len(t.splits)-1].Memo = &memo
		return nil

	case '$': // Amount
		amt, err := parseAmount(line[1:])
		if err != nil {
			return errors.Wrap(err, "failed to parse split amount")
		}

		// This could be the first element of a new split, but only if there
		// isn't an existing split, or the existing split already has '$' field.
		if len(t.splits) == 0 || t.splits[len(t.splits)-1].Amount != nil {
			t.splits = append(t.splits, Split{})
		}

		t.splits[len(t.splits)-1].Amount = &amt
		return nil

	default:
		return UnsupportedFieldError(
			errors.Errorf("cannot process line '%s'", line))
	}
}

// A Split is used to tag part of a transaction with a separate category and
// description.
type Split struct {

	// Category of this transaction split.
	Category *string

	// Memo is a string description of the transaction split.
	Memo *string

	// Amount stores the transaction split value in minor currency units. For
	// instance, a $12.99 transaction will be 1299.
	Amount *int
}

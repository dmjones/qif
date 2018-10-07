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

// A BankingTransaction contains the information associated with non-investment transactions (i.e.
// Cash, Bank and CCard account types).
type BankingTransaction interface {
	Transaction

	// Num contains the check or reference number for the transaction. Wikipedia suggest this may also contain
	// "Deposit", "Transfer", "Print", "ATM", or "EFT".
	Num() string

	// Payee describes the recipient of the transaction.
	Payee() string

	// Address contains no more than five address lines for the payee. Wikipedia suggests the first entry is usually
	// the same as the Payee field.
	Address() []string

	// AddressMessage contains an additional message associated with the payee address. This is only non-empty if
	// the transaction address contained a special sixth line.
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
}

// A Split is used to tag part of a transaction with a separate category and description.
type Split interface {

	// Category of this transaction split.
	Category() string

	// Memo is a string description of the transaction split.
	Memo() string

	// Amount stores the transaction split value in minor currency units. For instance, a $12.99 transaction
	// will be 1299.
	Amount() int
}

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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func TestUnexpectedEOF(t *testing.T) {
	inputData := strings.Join([]string{
		bankHeader,
		"Mmemo",
		"T-99.50",
	}, "\n")

	r := NewReader(strings.NewReader(inputData))

	tx, err := r.ReadAll()
	assert.Nil(t, tx)

	e, ok := err.(RecordEndError)
	require.True(t, ok)

	assert.Equal(t, "memo", e.Incomplete.Memo())
	assert.Equal(t, -9950, e.Incomplete.Amount())
}

func TestBadHeader(t *testing.T) {
	inputData := strings.Join([]string{
		"!Type:Bonk",
		"Mmemo",
		"T-99.50",
	}, "\n")

	r := NewReader(strings.NewReader(inputData))

	tx, err := r.ReadAll()
	assert.Nil(t, tx)
	assert.Error(t, err)
}

func TestReadOne(t *testing.T) {
	inputData := strings.Join([]string{
		cardHeader,
		"Mmemo",
		"T-99.50",
		recordEnd,
		"Aaddress1",
		"Aaddress2",
		"T123.00",
		recordEnd,
	}, "\n")

	r := NewReader(strings.NewReader(inputData))

	tx, err := r.Read()
	assert.NoError(t, err)

	assert.Equal(t, "memo", tx.Memo())
	assert.Equal(t, -9950, tx.Amount())

	tx, err = r.Read()
	assert.NoError(t, err)

	btx := tx.(BankingTransaction)
	assert.Equal(t, []string{"address1", "address2"}, btx.Address())
	assert.Equal(t, 12300, btx.Amount())
}

func strptr(s string) *string {
	return &s
}

func intptr(i int) *int {
	return &i
}

func TestSpecExample1(t *testing.T) {
	// This file includes the example from the QIF spec.
	// Note: I think there is a mistake in the spec example. The second split in
	// the first example record is "$=746.36". I'm assuming this is a typo and
	// was meant to be "$-746.36". This spec page shows it without the typo:
	// http://moneymvps.org/articles/qifspecification.aspx

	input, err := os.Open("testdata/example1.qif")
	require.NoError(t, err)
	defer input.Close()

	expected1 := &bankingTransaction{
		num:      "1005",
		payee:    "Bank Of Mortgage",
		category: "[linda]",
		splits: []Split{
			{Category: strptr("[linda]"), Amount: intptr(-25364)},
			{Category: strptr("Mort Int"), Amount: intptr(-74636)},
		},
	}
	expected1.date, err = time.Parse("1/ 2/06", "6/ 1/94")
	require.NoError(t, err)
	expected1.amount = -100000
	expected1.amountDecimal,_ = decimal.NewFromString("-1000.00")

	expected2 := &bankingTransaction{
		payee: "Deposit",
	}
	expected2.date, err = time.Parse("1/ 2/06", "6/ 2/94")
	require.NoError(t, err)
	expected2.amount = 7500
	expected2.amountDecimal,_ = decimal.NewFromString("75.00")

	expected3 := &bankingTransaction{
		payee:    "Anthony Hopkins",
		address:  []string{"P.O. Box 27027", "Tucson, AZ", "85726", "", ""},
		category: "Entertain",
	}
	expected3.date, err = time.Parse("1/ 2/06", "6/ 3/94")
	require.NoError(t, err)
	expected3.amount = -1000
	expected3.amountDecimal,_ = decimal.NewFromString("-10.00")
	expected3.memo = "Film"

	expected := []Transaction{expected1, expected2, expected3}

	reader := NewReader(input)

	txs, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, txs)
}

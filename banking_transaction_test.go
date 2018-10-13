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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckNumber(t *testing.T) {
	tx := &bankingTransaction{}
	const checkNum = "num123"

	err := tx.parseBankingTransactionField("N"+checkNum, Config{})
	require.NoError(t, err)

	assert.Equal(t, checkNum, tx.Num())
}

func TestCheckPayee(t *testing.T) {
	tx := &bankingTransaction{}
	const payee = "fred"

	err := tx.parseBankingTransactionField("P"+payee, Config{})
	require.NoError(t, err)

	assert.Equal(t, payee, tx.Payee())
}

func TestCheckCategory(t *testing.T) {
	tx := &bankingTransaction{}
	const category = "cat"

	err := tx.parseBankingTransactionField("L"+category, Config{})
	require.NoError(t, err)

	assert.Equal(t, category, tx.Category())
}

func Test5LineAddress(t *testing.T) {
	tx := &bankingTransaction{}

	address := []string{"a1", "a2", "a3", "a4", "a5"}

	for _, a := range address {
		err := tx.parseBankingTransactionField("A"+a, Config{})
		require.NoError(t, err)
	}
	assert.Equal(t, address, tx.Address())
	assert.Equal(t, "", tx.AddressMessage())
}

func Test6LineAddress(t *testing.T) {
	tx := &bankingTransaction{}

	address := []string{"a1", "a2", "a3", "a4", "a5", "msg"}

	for _, a := range address {
		err := tx.parseBankingTransactionField("A"+a, Config{})
		require.NoError(t, err)
	}
	assert.Equal(t, address[:5], tx.Address())
	assert.Equal(t, address[5], tx.AddressMessage())
}

func TestSplits(t *testing.T) {
	tx := &bankingTransaction{}

	lines := []string{
		// split 1
		"Scat1",
		"Ememo1",
		"$12.99",

		// split 2
		"$3.99",

		// split 3
		"Ememo3",
	}

	for _, l := range lines {
		err := tx.parseBankingTransactionField(l, Config{})
		require.NoError(t, err)
	}
	require.Equal(t, 3, len(tx.Splits()))

	assert.Equal(t, "cat1", *tx.Splits()[0].Category)
	assert.Equal(t, "memo1", *tx.Splits()[0].Memo)
	assert.Equal(t, 1299, *tx.Splits()[0].Amount)

	assert.Nil(t, tx.Splits()[1].Category)
	assert.Nil(t, tx.Splits()[1].Memo)
	assert.Equal(t, 399, *tx.Splits()[1].Amount)

	assert.Nil(t, tx.Splits()[2].Category)
	assert.Equal(t, "memo3", *tx.Splits()[2].Memo)
	assert.Nil(t, tx.Splits()[2].Amount)
}

func TestTransactionField(t *testing.T) {
	tx := &bankingTransaction{}
	const memo = "memo"

	err := tx.parseBankingTransactionField("M"+memo, Config{})
	require.NoError(t, err)

	assert.Equal(t, memo, tx.Memo())
}

func TestEmptyLine(t *testing.T) {
	tx := &bankingTransaction{}
	err := tx.parseBankingTransactionField("", Config{})
	assert.Error(t, err)
}

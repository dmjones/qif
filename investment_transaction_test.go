// SPDX-License-Identifier: Apache-2.0

package qif

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckAction(t *testing.T) {
	tx := &investmentTransaction{}
	const action = "Sell"

	err := tx.parseInvestmentTransactionField("N"+action, Config{})
	require.NoError(t, err)

	assert.Equal(t, action, tx.ActionString())
	assert.Equal(t, ActionSell, tx.Action())
}

func TestCheckPrice(t *testing.T) {
	tx := &investmentTransaction{}
	const price = "23.49"

	err := tx.parseInvestmentTransactionField("I"+price, Config{})
	require.NoError(t, err)

	priceDecimal, _ := decimal.NewFromString(price)
	assert.Equal(t, true, priceDecimal.Equal(tx.Price()))
}

func TestCheckShares(t *testing.T) {
	tx := &investmentTransaction{}
	const shares = "401"

	err := tx.parseInvestmentTransactionField("Q"+shares, Config{})
	require.NoError(t, err)

	sharesDecimal, _ := decimal.NewFromString(shares)
	assert.Equal(t, true, sharesDecimal.Equal(tx.Shares()))
}

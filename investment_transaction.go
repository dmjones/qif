// SPDX-License-Identifier: Apache-2.0

package qif

import "github.com/pkg/errors"
import "github.com/shopspring/decimal"

type InvestmentAction uint

// Constant value for each Investment Action
const (
	ActionUndefined InvestmentAction = iota
	ActionBuy
	ActionBuyX
	ActionSell
	ActionSellX
	ActionCGLong
	ActionCGLongX
	ActionCGMid
	ActionCGMidX
	ActionCGShort
	ActionCGShortX
	ActionDiv
	ActionDivX
	ActionIntInc
	ActionIntIncX
	ActionReInvDiv
	ActionReInvInt
	ActionReInvLg
	ActionReInvMd
	ActionReInvSh
	ActionReprice
	ActionXIn
	ActionXOut
	ActionMiscExp
	ActionMiscExpX
	ActionMiscInc
	ActionMiscIncX
	ActionMarginInt
	ActionMarginIntX
	ActionReturnCap
	ActionReturnCapX
	ActionStockSplit
	ActionSharesOut
	ActionSharesIn
)

var (
    actionMap = map[string]InvestmentAction {
	"Buy":      ActionBuy,
	"BuyX":     ActionBuyX,
	"Sell":     ActionSell,
	"SellX":    ActionSellX,
	"CGLong":   ActionCGLong,
	"CGLongX":  ActionCGLongX,
	"CGMid":    ActionCGMid,
	"CGMidX":   ActionCGMidX,
	"CGShort":  ActionCGShort,
	"CGShortX": ActionCGShortX,
	"Div":      ActionDiv,
	"DivX":     ActionDivX,
	"IntInc":   ActionIntInc,
	"IntIncX":  ActionIntIncX,
	"ReInvDiv": ActionReInvDiv,
	"ReInvInt": ActionReInvInt,
	"ReInvLg":  ActionReInvLg,
	"ReInvMd":  ActionReInvMd,
	"ReInvSh":  ActionReInvSh,
	"Reprice":  ActionReprice,
	"XIn":      ActionXIn,
	"XOut":     ActionXOut,
	"MiscExp":  ActionMiscExp,
	"MiscExpX": ActionMiscExpX,
	"MiscInc":  ActionMiscInc,
	"MiscIncX": ActionMiscIncX,
	"MargInt":  ActionMarginInt,
	"MargIntX": ActionMarginIntX,
	"RtrnCap":  ActionReturnCap,
	"RtrnCapX": ActionReturnCapX,
	"StkSplit": ActionStockSplit,
	"ShrsOut":  ActionSharesOut,
	"ShrsIn":   ActionSharesIn,
    }
)

// An InvestmentTransaction contains the information associated with
// transactions for Investment accounts.
type InvestmentTransaction interface {
	Transaction

	// Investment Action (Buy, Sell, etc.) as type InvestmentAction.
	Action() InvestmentAction

	// Investment Action (Buy, Sell, etc.) as type string.
	ActionString() string

	// Name of Security that was bought or sold (name of Stock, Fund, etc).
	SecurityName() string

	// Quantity of Shares (or split ratio, if Action is StkSplit).
	Shares() decimal.Decimal

	// Price which Investment Action was executed at.
	Price() decimal.Decimal

	// Commission cost (generally trades are commission-free these days).
	Commission() decimal.Decimal
}

type investmentTransaction struct {
	transaction
	action		InvestmentAction
	actionString	string
	securityName	string
	shares		decimal.Decimal
	price		decimal.Decimal
	commission	decimal.Decimal
}

func (t *investmentTransaction) Action() InvestmentAction {
	return t.action
}

func (t *investmentTransaction) ActionString() string {
	return t.actionString
}

func (t *investmentTransaction) SecurityName() string {
	return t.securityName
}

func (t *investmentTransaction) Shares() decimal.Decimal {
	return t.shares
}

func (t *investmentTransaction) Price() decimal.Decimal {
	return t.price
}

func (t *investmentTransaction) Commission() decimal.Decimal {
	return t.commission
}

func (t *investmentTransaction) parseInvestmentTransactionField(line string,
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
		t.actionString = line[1:]
		t.action = actionMap[t.actionString]
		return nil
	case 'Y':
		t.securityName = line[1:]
		return nil
	case 'Q':
		shares := line[1:]
		t.shares, err = decimal.NewFromString(shares)
		return err
	case 'I':
		price := line[1:]
		t.price, err = decimal.NewFromString(price)
		return err
	case 'O':
		commission := line[1:]
		t.commission, err = decimal.NewFromString(commission)
		return err
	default:
		return UnsupportedFieldError(
			errors.Errorf("cannot process line '%s'", line))
	}
}

func (t *investmentTransaction) parseTransactionTypeField(line string,
							  config Config) error {
	return t.parseInvestmentTransactionField(line, config)
}

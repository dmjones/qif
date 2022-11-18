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
	"strconv"
	"time"

	"fmt"

	"regexp"

	"strings"

	"github.com/shopspring/decimal"
	"github.com/pkg/errors"
)

type ClearedStatus int

const (
	UnknownStatus ClearedStatus = iota
	Cleared                     = iota
	Reconciled                  = iota
	NotCleared                  = iota
)

type UnsupportedFieldError error

// A Transaction contains the fields common to all transaction types.
type Transaction interface {

	// Date contains the year, month and day of the transaction. All other
	// fields are zero.
	Date() time.Time

	// Amount stores the transaction value in minor currency units. For
	// instance, a $12.99 transaction will be 1299.
	Amount() int

	// Amount stored as decimal.Decimal
	AmountDecimal() decimal.Decimal

	// Memo is a string description of the transaction.
	Memo() string

	// Status indicates if the transaction is cleared. The value will be
	// UnknownStatus if the transaction data did not specify a value for this
	// field.
	Status() ClearedStatus
}

type transaction struct {
	date   time.Time
	amount int
	amountDecimal decimal.Decimal
	memo   string
	status ClearedStatus
}

func (t *transaction) Date() time.Time {
	return t.date
}

func (t *transaction) Amount() int {
	return t.amount
}

func (t *transaction) AmountDecimal() decimal.Decimal {
	return t.amountDecimal
}

func (t *transaction) Memo() string {
	return t.memo
}

func (t *transaction) Status() ClearedStatus {
	return t.status
}

func (t *transaction) parseTransactionField(line string, config Config) error {
	if line == "" {
		return errors.New("line is empty")
	}

	switch line[0] {
	case 'D':
		date, err := parseDate(line[1:], config.DayFirst)
		if err != nil {
			return errors.Wrap(err, "failed to parse date")
		}
		t.date = date
		return nil

	case 'T', 'U': // Wikipedia suggests 'U' is a synonym for 'T'
		amount := strings.Replace(line[1:], ",", "", -1)
		amt, err := parseAmount(amount)
		if err != nil {
			return errors.Wrap(err, "failed to parse amount")
		}
		t.amount = amt
		t.amountDecimal, err = decimal.NewFromString(amount)
		return err

	case 'M':
		t.memo = line[1:]
		return nil

	case 'C':
		status, err := parseClearedStatus(line[1:])
		if err != nil {
			return errors.Wrap(err, "failed to parse cleared status")
		}
		t.status = status
		return nil

	default:
		return UnsupportedFieldError(errors.Errorf("cannot process line '%s'", line))
	}
}

func parseClearedStatus(s string) (ClearedStatus, error) {
	switch s {
	case "*", "c":
		return Cleared, nil
	case "X", "R":
		return Reconciled, nil
	case "":
		return NotCleared, nil

	default:
		return UnknownStatus, errors.Errorf(`bad cleared status: "%s"`, s)
	}
}

// parseAmount converts an amount string (such as '12.99') into minor currency
// units.
func parseAmount(s string) (int, error) {

	sMod := strings.Replace(s, ",", "", -1)

	// Check format. Expect an optional minus or plus sign, then either a whole
	// number or a decimal with numbers either side of the point.
	re := regexp.MustCompile(`[\-\+]?\d+\.\d{1,2}`)

	if !re.MatchString(sMod) {
		return 0, errors.Errorf(`bad amount string "%s"`, s)
	}

	// If there's only one number after the decimal point, add a zero
	if strings.Index(sMod, ".") == len(sMod)-2 {
		sMod = sMod + "0"
	}

	return strconv.Atoi(strings.Replace(sMod, ".", "", 1))
}

// parseDate attempts to parse the given string with a variety of formats.
// monthFirst controls whether mm/dd or dd/mm formats are used.
func parseDate(s string, dayFirst bool) (date time.Time, err error) {
	// The spec is vague on date formats. Based on wikipedia and other sources,
	// this is a potential list of valid options (using Go reference time of
	// Mon Jan 2 15:04:05 -0700 MST 2006).

	// Dates and months may or may not have leading zeroes. To reduce
	// permutations, remove any leading zeros:
	re := regexp.MustCompile(`(^|[^\d])0`)
	sMod := re.ReplaceAllString(s, "$1")

	// Some of the examples have spaces between days and months, in numeric
	// form. Let's remove all spaces to be safe.
	re = regexp.MustCompile(`^[\d/ ]+$`)
	if re.MatchString(sMod) {
		sMod = strings.Replace(sMod, " ", "", -1)
	}

	// Go's date parsing doesn't seem to cope with single digit years.
	re = regexp.MustCompile(`'(\d)$`)
	if re.MatchString(sMod) {
		now := time.Now()
		decimal := (now.Year() - 2000) / 10
		sMod = re.ReplaceAllString(sMod, fmt.Sprintf("%d$1", decimal))
	}

	var first, second string
	if dayFirst {
		first = "2"
		second = "1"
	} else {
		first = "1"
		second = "2"
	}

	formats := []string{
		"2 January 2006",
		"2 January 06",

		fmt.Sprintf("%s/%s/2006", first, second),
		fmt.Sprintf("%s/%s/06", first, second),
	}

	for _, f := range formats {
		date, err = time.Parse(f, sMod)
		if err == nil {
			return
		}
	}

	err = errors.Errorf(`failed to parse date "%s"`, s)
	return
}

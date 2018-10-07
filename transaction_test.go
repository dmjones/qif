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
	"time"
)

func TestDateParse(t *testing.T) {
	expectedDate, err := time.Parse("02/01/2006", "01/03/2017")
	require.NoError(t, err)

	inputs := []string{
		"1 March 2017",
		"1 March 17",
		"1 March '7",
		"03/01/2017",
		"03/01/17",
		"03/01/'7",
		"3/ 1/2017",
		"03/1/2017",
	}

	for _, i := range inputs {
		date, err := parseDate(i, true)
		assert.NoError(t, err)
		assert.Equalf(t, expectedDate, date, "failed for input %s", i)
	}
}

func TestAmountParse(t *testing.T) {
	vectors := map[string]int{
		"12.99":  1299,
		"+12.99": 1299,
		"-12.99": -1299,
		"-12.9":  -1290,
	}

	for k, v := range vectors {
		res, err := parseAmount(k)
		assert.NoErrorf(t, err, "error processing '%s'", k)
		assert.Equalf(t, v, res, "error processing '%s", k)
	}

	badVectors := []string{
		"12.",
		"12",
		".9",
		"+-12.00",
	}

	for _, v := range badVectors {
		_, err := parseAmount(v)
		assert.Errorf(t, err, "error processing '%s'", v)
	}
}

func TestClearedStatus(t *testing.T) {

	vectors := map[string]ClearedStatus{
		"":  NotCleared,
		"*": Cleared,
		"c": Cleared,
		"X": Reconciled,
		"R": Reconciled,
	}

	for k, v := range vectors {
		res, err := parseClearedStatus(k)
		assert.NoError(t, err)
		assert.EqualValues(t, v, res)
	}

	_, err := parseClearedStatus("Z") // not real
	assert.Error(t, err)
}

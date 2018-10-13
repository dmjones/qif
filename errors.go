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

import "fmt"

// A RecordEndError is returned if the input data finishes without a terminating
// '^' character. All records should be terminated in a QIF file, but an
// application may wish to be forgiving if the last record is not terminated.
// The non-terminated transaction can be found as a field within the error.
type RecordEndError struct {

	// Incomplete is the transaction that was being parsed when the input
	// ended.
	Incomplete Transaction
}

func (RecordEndError) Error() string {
	return fmt.Sprintf("unexpected end of input")
}

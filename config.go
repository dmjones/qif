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

// Config defines the configuration of the reader.
type Config struct {

	// DayFirst specifies whether to interpret dates as mm/dd or dd/mm.
	DayFirst bool
}

// DefaultConfig returns the default configuration used by NewReader:
//
//  Config{
//    DayFirst: false,
//  }
func DefaultConfig() Config {
	return Config{
		DayFirst: false,
	}
}

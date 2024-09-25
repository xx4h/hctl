// Copyright 2024 Fabian `xx4h` Sylvester
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"regexp"
	"testing"
)

func Test_printVersion(t *testing.T) {

	var tests = map[string]struct {
		short  bool
		rexOut string
		rexErr string
	}{
		"normal version": {
			false,
			"(?m)^Version:.*\n^Commit:.*\n^Date:.*$",
			"",
		},
		"short version": {
			true,
			"(?m)^Version:.*\n^Commit:.*\n^Date:.*$",
			"",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			b := new(bytes.Buffer)
			printVersion(b, tt.short)
			o := b.String()
			if ok, _ := regexp.MatchString(tt.rexOut, o); !ok {
				t.Errorf("got %q, want %q", o, tt.rexOut)
			}
		})
	}
}

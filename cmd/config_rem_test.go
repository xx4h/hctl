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

package cmd

import (
	"testing"
)

func Test_newCmdConfigRem(t *testing.T) {
	h := newTestingHctl(t)

	var tests = map[string]cmdTest{
		"rem completion.short_names": {
			"config remove device_map.g",
			"(?m)^.*Option `device_map.g` successfully removed",
			"",
		},
	}

	testCmd(t, h, tests)
}

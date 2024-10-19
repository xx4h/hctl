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

	"github.com/xx4h/hctl/pkg/hctltest"
)

func Test_newCmdList(t *testing.T) {
	ms := hctltest.MockServer(t)
	h := newTestingHctl(t)
	if err := h.SetConfigValue("hub.url", ms.URL); err != nil {
		t.Error(err)
	}

	var tests = map[string]cmdTest{
		"list services": {
			"list services",
			`^.*Services`,
			"",
		},
		"list services with play_media": {
			"list services -s play_media",
			`(?m)^.*Services.*\n.*media_player.*\n.*play_media`,
			"",
		},
		"list entities": {
			"list entities",
			`^.*States`,
			"",
		},
		"list entities of domain media_player": {
			"list entities -d media_player",
			`^.*States.*\n.*media_player.*\n.*player1.*\n.*player2`,
			"",
		},
	}

	testCmd(t, h, tests)
}

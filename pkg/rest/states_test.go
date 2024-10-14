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

package rest

import (
	"reflect"
	"testing"

	"github.com/xx4h/hctl/pkg/hctltest"
)

const (
	statesCount = 10
)

func Test_GetStates(t *testing.T) {
	ms := hctltest.MockServer(t)
	h := &Hass{
		APIURL: ms.URL,
		Token:  "test_token",
	}
	defer ms.Close()
	s, err := h.GetStates()
	if err != nil {
		t.Errorf("Error getting states: %v", err)
	}
	st := reflect.TypeOf(s)
	wt := reflect.TypeOf([]HassState{})
	if st != wt {
		t.Errorf("got %s, want %s", st, wt)
	}
	cs := len(s)
	if cs != statesCount {
		t.Errorf("got %d, want %d", cs, statesCount)
	}
}

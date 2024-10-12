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
)

const (
	serviceCount = 47
)

func Test_GetServices(t *testing.T) {
	ms := mockServerGetDataFromFile(t, "services.json")
	h := &Hass{
		APIURL: ms.URL,
		Token:  "test_token",
	}
	defer ms.Close()
	s, err := h.GetServices()
	if err != nil {
		t.Errorf("Error getting services: %v", err)
	}
	st := reflect.TypeOf(s)
	wt := reflect.TypeOf([]HassService{})
	if st != wt {
		t.Errorf("got %s, want %s", st, wt)
	}
	cs := len(s)
	if cs != serviceCount {
		t.Errorf("got %d, want %d", cs, serviceCount)
	}
}

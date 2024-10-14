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
	"testing"

	"github.com/xx4h/hctl/pkg/hctltest"
)

func Test_Toggle(t *testing.T) {
	ms := hctltest.MockServer(t)
	h := &Hass{
		APIURL: ms.URL,
		Token:  "test_token",
	}
	defer ms.Close()
	obj, action, sub, err := h.Toggle("bedroom_main")
	if err != nil {
		t.Errorf("Error toggling: %v", err)
	}
	if obj != "bedroom_main" {
		t.Errorf("got %s, want bedroom_main", obj)
	}
	if action != "toggle" {
		t.Errorf("got %s, want toggle", action)
	}
	if sub != "light" {
		t.Errorf("got %s, want light", sub)
	}
}

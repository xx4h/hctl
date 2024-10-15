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

	"github.com/xx4h/hctl/pkg"
	"github.com/xx4h/hctl/pkg/hctltest"
)

func Test_compListStates(t *testing.T) {
	ms := hctltest.MockServer(t)
	h, err := pkg.NewHctl(false)
	if err != nil {
		t.Errorf("Error createing new Hctl instance: %+v", err)
	}
	if err := h.SetConfigValue("hub.url", ms.URL); err != nil {
		t.Errorf("Could not set hub.url to %s: %+v", ms.URL, err)
	}
	if err := h.SetConfigValue("hub.token", "testtoken"); err != nil {
		t.Errorf("Could not set hub.token to %s: %+v", ms.URL, err)
	}
	defer ms.Close()

	tests := map[string]struct {
		ignoredStates []string
		serviceCaps   []string
		attributes    []string
		state         string
		expectedCount int
	}{
		"brightness attrib": {
			nil,
			nil,
			[]string{"brightness"},
			"",
			5,
		},
		"serviceCap turn_on": {
			nil,
			[]string{"turn_on"},
			nil,
			"",
			10,
		},
		"serviceCap turn_on + state off": {
			nil,
			[]string{"turn_on"},
			nil,
			"off",
			4,
		},
		"serviceCap turn_off + state on": {
			nil,
			[]string{"turn_off"},
			nil,
			"on",
			6,
		},
		"serviceCap play_media + attrib device_class": {
			nil,
			[]string{"play_media"},
			[]string{"device_class"},
			"",
			1,
		},
		"serviceCap play_media + attrib device_class or video_out": {
			nil,
			[]string{"play_media"},
			[]string{"device_class", "video_out"},
			"",
			2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s, _ := compListStates("", tt.ignoredStates, tt.serviceCaps, tt.attributes, tt.state, h)
			t.Logf("Completion states found: %+v", s)
			if len(s) != tt.expectedCount {
				t.Errorf("got %d, want %d", len(s), tt.expectedCount)
			}
		})
	}
}

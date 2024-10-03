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

func Test_FilterDomainsFromStates(t *testing.T) {
	var testInStates = []HassState{
		{EntityID: "light.foo", State: "off"},
		{EntityID: "light.bar", State: "on"},
		{EntityID: "switch.boo", State: "off"},
		{EntityID: "switch.far", State: "off"},
		{EntityID: "other.foo", State: "on"},
	}

	var tests = map[string]struct {
		domains []string
		in      []HassState
		out     []HassState
	}{
		"Empty states no domains": {
			[]string{},
			[]HassState{},
			[]HassState{},
		},
		"Empty states domains": {
			[]string{"foo", "bar"},
			[]HassState{},
			[]HassState{},
		},
		"States no domains": {
			[]string{},
			testInStates,
			testInStates,
		},
		"States no matching domains": {
			[]string{"foo", "bar"},
			testInStates,
			[]HassState{},
		},
		"States one matching domain": {
			[]string{"switch"},
			testInStates,
			[]HassState{
				{EntityID: "switch.boo", State: "off"},
				{EntityID: "switch.far", State: "off"},
			},
		},
		"States multiple matching domains": {
			[]string{"light", "switch"},
			testInStates,
			[]HassState{
				{EntityID: "light.foo", State: "off"},
				{EntityID: "light.bar", State: "on"},
				{EntityID: "switch.boo", State: "off"},
				{EntityID: "switch.far", State: "off"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// ensure we have a deep copy of tt.in as we don't want to affect other tests
			in := []HassState{}
			in = append(in, tt.in...)
			states := FilterDomainsFromStates(in, tt.domains)
			if len(in) == len(tt.out) && len(states) != len(tt.out) {
				t.Errorf("got %d states, want %d", len(states), len(tt.out))
			} else if !reflect.DeepEqual(states, tt.out) {
				t.Errorf("got %q states, want %q", states, tt.out)
			}
		})
	}

}

func Test_FilterDomainsFromServices(t *testing.T) {
	var testInServices = []HassService{
		{
			Domain: "switch",
			Services: map[string]HassDomainService{
				"turn_on": {
					Name: "Turn on",
				},
				"turn_off": {
					Name: "Turn off",
				},
				"toggle": {
					Name: "Toggle",
				},
			},
		},
		{
			Domain: "light",
			Services: map[string]HassDomainService{
				"turn_on": {
					Name: "Turn on",
				},
				"turn_off": {
					Name: "Turn off",
				},
				"toggle": {
					Name: "Toggle",
				},
			},
		},
		{
			Domain: "tts",
			Services: map[string]HassDomainService{
				"speak": {
					Name: "Speak",
				},
			},
		},
	}

	var tests = map[string]struct {
		domains []string
		in      []HassService
		out     []HassService
	}{
		"No domains": {
			[]string{},
			testInServices,
			testInServices,
		},
		"One non-matching domain": {
			[]string{"foobar"},
			testInServices,
			[]HassService{},
		},
		"Two non-matching domains": {
			[]string{"foobar", "barfoo"},
			testInServices,
			[]HassService{},
		},
		"One matching domain": {
			[]string{"light"},
			testInServices,
			[]HassService{{
				Domain: "light",
				Services: map[string]HassDomainService{
					"turn_on": {
						Name: "Turn on",
					},
					"turn_off": {
						Name: "Turn off",
					},
					"toggle": {
						Name: "Toggle",
					},
				},
			}},
		},
		"Two matching domains": {
			[]string{"light", "tts"},
			testInServices,
			[]HassService{{
				Domain: "light",
				Services: map[string]HassDomainService{
					"turn_on": {
						Name: "Turn on",
					},
					"turn_off": {
						Name: "Turn off",
					},
					"toggle": {
						Name: "Toggle",
					},
				},
			}, {
				Domain: "tts",
				Services: map[string]HassDomainService{
					"speak": {
						Name: "Speak",
					},
				},
			}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// ensure we have a deep copy of tt.in, as we don't want to affect other tests
			in := []HassService{}
			in = append(in, tt.in...)
			services := FilterDomainsFromServices(in, tt.domains)
			if len(services) != len(tt.out) {
				t.Errorf("got %d services, want %d", len(services), len(tt.out))
			} else if !reflect.DeepEqual(services, tt.out) {
				t.Errorf("got %q services, want %q", services, tt.out)
			}
		})
	}
}

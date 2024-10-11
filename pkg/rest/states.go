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
	"encoding/json"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

type HassState struct {
	EntityID   string         `json:"entity_id"`
	State      string         `json:"state"`
	Attributes map[string]any `json:"attributes"`
}

// TODO: Add sort, configurable, maybe even something like GetSortedStates(order ...string)?
// Get all states from Hass
func (h *Hass) GetStates() ([]HassState, error) {
	if h.States != nil {
		log.Info().Msg("Using cached states.")
		return h.States, nil
	}

	res, err := h.api("GET", "/states", nil)
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		return nil, err
	}

	states := []HassState{}
	err = json.Unmarshal(res, &states)
	if err != nil {
		log.Debug().Caller().Msgf("Could not create Services Object: %+v", err)
		return nil, err
	}

	h.States = states

	return states, nil
}

func (h *Hass) GetFilteredStates(domains []string) ([]HassState, error) {
	states, err := h.GetStates()
	if err != nil {
		return nil, err
	}
	return FilterDomainsFromStates(states, domains), nil
}

func (h *Hass) GetFilteredStatesMap(domains []string) (map[string][]string, error) {
	states, err := h.GetFilteredStates(domains)
	if err != nil {
		return nil, err
	}
	t := make(map[string][]string)
	for state := range states {
		elist := strings.Split(states[state].EntityID, ".")
		t[elist[0]] = append(t[elist[0]], elist[1])
	}
	return t, nil
}

func (h *Hass) GetStatesWithService(service string) ([]HassState, error) {
	var domainsWithService []string
	var statesWithService []HassState

	states, err := h.GetStates()
	if err != nil {
		return statesWithService, err
	}
	services, err := h.GetServices()
	if err != nil {
		return statesWithService, err
	}

	for _, svc := range services {
		for s := range svc.Services {
			if s == service {
				domainsWithService = append(domainsWithService, svc.Domain)
			}
		}
	}

	for d := range states {
		s := strings.Split(states[d].EntityID, ".")
		if slices.Contains(domainsWithService, s[0]) {
			statesWithService = append(statesWithService, states[d])
		}
	}

	return statesWithService, nil
}

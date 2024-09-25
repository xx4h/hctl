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
	EntityId   string         `json:"entity_id"`
	State      string         `json:"state"`
	Attributes map[string]any `json:"attributes"`
}

// TODO: Add sort, configurable, maybe even something like GetSortedStates(order ...string)?
// Get all states from Hass
func (h *Hass) GetStates() []HassState {
	if h.States != nil {
		log.Info().Msg("Using cached states.")
		return h.States
	}

	res, err := h.api("GET", "/states", nil)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
	}

	states := []HassState{}
	err = json.Unmarshal(res, &states)
	if err != nil {
		log.Error().Msgf("Could not create Services Object: %v", err)
	}

	h.States = states

	return states
}

func (h *Hass) GetFilteredStates(domains []string) []HassState {
	states := h.GetStates()
	return RemoveDomainsFromStates(states, domains)
}

func (h *Hass) GetFilteredStatesMap(domains []string) map[string][]string {
	states := h.GetFilteredStates(domains)
	t := make(map[string][]string)
	for state := range states {
		elist := strings.Split(states[state].EntityId, ".")
		t[elist[0]] = append(t[elist[0]], elist[1])
	}
	return t
}

func (h *Hass) GetStatesWithService(service string) ([]HassState, error) {
	var domainsWithService []string
	var statesWithService []HassState

	states := h.GetStates()
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
		s := strings.Split(states[d].EntityId, ".")
		if slices.Contains(domainsWithService, s[0]) {
			statesWithService = append(statesWithService, states[d])
		}
	}

	return statesWithService, nil
}

func (h *Hass) hasEntityInDomain(state string, domain string) bool {
	states := h.GetStates()

	for d := range states {
		s := strings.Split(states[d].EntityId, ".")
		if domain == s[0] && state == s[1] {
			return true
		}
	}
	return false
}

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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rs/zerolog/log"
)

type Hass struct {
	APIURL    string
	Token     string
	Fuzz      bool
	States    []HassState
	Services  []HassService
	DeviceMap map[string]string

	Result HassResult
}

type HassResult struct {
	EntityID    string `json:"entity_id"`
	State       string `json:"state"`
	TargetState string
	Attributes  HassResponseAttributes `json:"attributes"`
}

type HassResponseAttributes struct {
	FriendlyName string `json:"friendly_name"`
}

func New(apiURL string, token string, fuzz bool, deviceMap map[string]string) *Hass {
	return &Hass{APIURL: apiURL, Token: token, Fuzz: fuzz, DeviceMap: deviceMap}
}

func (h *Hass) preflight() error {
	if h.APIURL == "" {
		return errors.New("no Hub URL found: Run `hctl init` or manually create config")
	}
	if h.Token == "" {
		return errors.New("no Hub Token found: Run `hctl init` or manually create config")
	}
	return nil
}

func (h *Hass) createRequest(meth string, path string, payload map[string]any) (*http.Request, error) {
	var r io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(b)
	}

	return http.NewRequest(meth, fmt.Sprintf("%s%s", h.APIURL, path), r)
}

func (h *Hass) api(meth string, path string, payload map[string]any) ([]byte, error) {
	if err := h.preflight(); err != nil {
		return nil, err
	}

	req, err := h.createRequest(meth, path, payload)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Requesting URL %s, Method %s, Payload: %#v", req.URL, req.Method, payload)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	rData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return rData, nil
}

// TODO: Rework to return result list and work with it
func (h *Hass) getResult(res []byte) error {
	var result []HassResult

	if err := json.Unmarshal(res, &result); err != nil {
		log.Debug().Caller().Msgf("Failed to Unmarshal: %+v", string(res))
		return err
	}

	log.Debug().Caller().Msgf("Result: %#v", result)
	return nil
}

func getFuzz(name string, names []string) (int, bool) {
	// get all fuzzy matches with ranks
	ranks := fuzzy.RankFind(name, names)
	log.Debug().Caller().Msgf("Found Fuzzy Matches: %+v", ranks)

	// As Levenshtein is only positive, starting distance -1 is also indicator if there is a match
	distance := -1
	var position int

	// go through found fuzzy matches with ranks
	for m := range ranks {
		// Levenshtein distance is lower than current
		if distance == -1 || ranks[m].Distance < distance {
			// save new current distance
			distance = ranks[m].Distance
			// save position of original index
			position = ranks[m].OriginalIndex
		}
	}
	return position, distance > -1
}

func (h *Hass) resolveMapping(domain, name string) (string, string) {
	if val, ok := h.DeviceMap[name]; ok && domain == "" {
		log.Debug().Caller().Msgf("Found `%s` in device_map: %s", name, val)
		domain, name = splitDomainAndName(val)
		h.Fuzz = false
	}
	return domain, name
}

// Find matching entity for provided service
// Return error if none has been found
func (h *Hass) findEntity(name string, domain string, service string) (string, string, error) {
	states, err := h.GetStatesWithService(service)
	if err != nil {
		return "", "", err
	}

	domain, name = h.resolveMapping(domain, name)

	var names []string

	for i := range states {
		d, n := splitDomainAndName(states[i].EntityID)
		// when domain is set, but not matching, we can continue
		if domain != "" && domain != d {
			continue
		}
		// domain unset or matching
		// directly return find when entity matches name
		if n == name {
			return d, n, nil
		}

		// add to fuzz checker names list when fuzz enabled
		if h.Fuzz {
			names = append(names, n)
		}
	}

	// when fuzz enabled
	if h.Fuzz {
		if p, ok := getFuzz(name, names); ok {
			d, n := splitDomainAndName(states[p].EntityID)
			// get domain and entity name from original states array by position
			return d, n, nil
		}
	}
	return "", "", fmt.Errorf("no Entity %s capable of %s", name, service)
}

func (h *Hass) entityArgHandler(args []string, service string) (string, string, error) {
	domain, name := splitDomainAndName(args[0])
	if len(args) == 1 {
		return h.findEntity(name, domain, service)
	} else if len(args) == 2 {
		if domain == "" {
			return args[0], args[1], nil // #nosec G602
		}
		return h.findEntity(name, domain, service)
	}
	return "", "", fmt.Errorf("entityArgHandler has to many entries in args: %d", len(args))
}

// Returns domain, name
func splitDomainAndName(s string) (string, string) {
	p := strings.Split(s, ".")
	// we don't have any dots, so we only have a name
	if len(p) == 1 {
		return "", p[0]
	}
	// at least one dot, first element is considered the domain
	return p[0], strings.Join(p[1:], ".")
}

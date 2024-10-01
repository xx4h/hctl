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

import "fmt"

func (h *Hass) turn(state string, sub string, obj string) error {
	hasDomain, err := h.hasDomainWithService(sub, fmt.Sprintf("turn_%s", state))
	if err != nil {
		return err
	} else if !hasDomain {
		return fmt.Errorf("No such Domain with Service: %s with %s", sub, fmt.Sprintf("turn_%s", state))
	}
	if !h.hasEntityInDomain(obj, sub) {
		return fmt.Errorf("No such Entity in Domain: %s in %s", obj, sub)
	}
	payload := map[string]string{"entity_id": fmt.Sprintf("%s.%s", sub, obj)}
	res, err := h.api("POST", fmt.Sprintf("/services/%s/turn_%s", sub, state), payload)
	if err != nil {
		return err
	}

	if err := h.getResult(res); err != nil {
		return err
	}

	return nil
}

func (h *Hass) TurnOff(args ...string) (string, string, string, error) {
	sub, obj, err := h.entityArgHandler(args, "turn_off")
	if err != nil {
		return "", "", "", err
	}
	return obj, "off", sub, h.turn("off", sub, obj)
}

func (h *Hass) TurnOn(args ...string) (string, string, string, error) {
	sub, obj, err := h.entityArgHandler(args, "turn_on")
	if err != nil {
		return "", "", "", err
	}
	return obj, "on", sub, h.turn("on", sub, obj)
}

func (h *Hass) TurnLightOff(obj string) (string, string, string, error) {
	return h.TurnOff("light", obj)
}

func (h *Hass) TurnLightOn(obj string) (string, string, string, error) {
	return h.TurnOn("light", obj)
}

func (h *Hass) ToggleLight(obj string) (string, string, string, error) {
	return h.Toggle("light", obj)
}

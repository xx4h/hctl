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
	"fmt"
	"strconv"
)

func (h *Hass) turn(state, domain, device, brightness string) error {
	// if err := h.checkEntity(sub, fmt.Sprintf("turn_%s", state), obj); err != nil {
	// 	return err
	// }
	payload := map[string]any{"entity_id": fmt.Sprintf("%s.%s", domain, device)}
	if brightness != "" {
		payload["brightness"] = brightness
	}
	res, err := h.api("POST", fmt.Sprintf("/services/%s/turn_%s", domain, state), payload)
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
	return obj, "off", sub, h.turn("off", sub, obj, "")
}

func (h *Hass) TurnOn(args ...string) (string, string, string, error) {
	sub, obj, err := h.entityArgHandler(args, "turn_on")
	if err != nil {
		return "", "", "", err
	}
	return obj, "on", sub, h.turn("on", sub, obj, "")
}

func (h *Hass) brightStep(domain, device, updown string) (string, error) {
	state, err := h.GetState(domain, device)
	if err != nil {
		return "", err
	}

	curany, ok := state.Attributes["brightness"]
	if !ok {
		return "", fmt.Errorf("state `%s.%s` has no attribute `brightness`", domain, device)
	}

	i, err := strconv.Atoi(fmt.Sprintf("%.0f", curany))
	if err != nil {
		return "", err
	}

	diff := i % 10
	switch updown {
	case "+":
		b := i + (10 - diff)
		if b == 100 {
			b = 99
		}
		return fmt.Sprintf("%d", b), nil
	case "-":
		b := i - diff
		if diff == 0 {
			b = i - 10
		}
		if b == 0 {
			b = 1
		}
		return fmt.Sprintf("%d", b), nil
	default:
		return "", fmt.Errorf("no such brightStep: %s", updown)
	}
}

func (h *Hass) TurnLightOnBrightness(device, brightness string) (string, string, string, error) {
	domain, device, err := h.entityArgHandler([]string{device}, "turn_on")
	switch brightness {
	case "-":
		brightness, err = h.brightStep(domain, device, "-")
	case "+":
		brightness, err = h.brightStep(domain, device, "+")
	case "min":
		brightness = "1"
	case "mid":
		brightness = "50"
	case "max":
		brightness = "99"
	}
	if err != nil {
		return "", "", "", err
	}
	return device, "on", domain, h.turn("on", domain, device, brightness)
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

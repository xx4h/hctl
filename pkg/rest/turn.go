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
	"math"
	"strconv"
	"strings"
)

func kelvinToMired(k int) int {
	return int(math.Round(1_000_000.0 / float64(k)))
}

func (h *Hass) turn(state, domain, device, brightness string, rgb []int, colorTemp int) error {
	// if err := h.checkEntity(sub, fmt.Sprintf("turn_%s", state), obj); err != nil {
	// 	return err
	// }
	payload := map[string]any{"entity_id": fmt.Sprintf("%s.%s", domain, device)}
	if brightness != "" {
		payload["brightness"] = brightness
	}

	if len(rgb) == 3 {
		payload["rgb_color"] = rgb
	}

	if colorTemp >= 153 && colorTemp <= 500 {
		payload["color_temp"] = colorTemp
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
	return obj, "off", sub, h.turn("off", sub, obj, "", nil, 0)
}

func (h *Hass) TurnOn(args ...string) (string, string, string, error) {
	sub, obj, err := h.entityArgHandler(args, "turn_on")
	if err != nil {
		return "", "", "", err
	}
	return obj, "on", sub, h.turn("on", sub, obj, "", nil, 0)
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

	// handle nil state (e.g. when device is off)
	if curany == nil {
		return "", nil
	}

	i, err := strconv.Atoi(fmt.Sprintf("%.0f", curany))
	if err != nil {
		return "", err
	}

	// Convert raw brightness (0-255) to 1-99 percentage scale
	p := int(math.Round(float64(i) / 255.0 * 99.0))
	if p < 1 {
		p = 1
	} else if p > 99 {
		p = 99
	}

	diff := p % 10
	switch updown {
	case "+":
		b := p + (10 - diff)
		if b > 99 {
			b = 99
		}
		return fmt.Sprintf("%d", b), nil
	case "-":
		b := p - diff
		if diff == 0 {
			b = p - 10
		}
		if b < 1 {
			b = 1
		}
		return fmt.Sprintf("%d", b), nil
	default:
		return "", fmt.Errorf("no such brightStep: %s", updown)
	}
}

func parseRGB(color string) ([]int, error) {
	var rgb []int
	parts := strings.Split(color, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("color must be in format R,G,B")
	}
	for _, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || val < 0 || val > 255 {
			return nil, fmt.Errorf("invalid RGB value: %s", part)
		}
		rgb = append(rgb, val)
	}
	return rgb, nil
}

func scaleBrightness(percent string) (string, error) {
	val, err := strconv.Atoi(percent)
	if err != nil || val < 1 || val > 99 {
		return "", fmt.Errorf("Invalid brightness percentage: %s", percent)
	}
	scaled := int(math.Round(float64(val) / 99.0 * 255.0))
	return fmt.Sprintf("%d", scaled), nil
}

func (h *Hass) TurnLightOnCustom(device, brightness string, color string, colorTemp int) (string, string, string, error) {
	domain, device, err := h.entityArgHandler([]string{device}, "turn_on")

	if color != "" && colorTemp != 0 {
		return "", "", "", fmt.Errorf("Cannot specify both RGB color and color temperature at the same time")
	}

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

	var rgb []int
	if color != "" {
		rgb, err = parseRGB(color)
		if err != nil {
			return "", "", "", err
		}
	}

	var brightnessScaled string
	if brightness != "" {
		brightnessScaled, err = scaleBrightness(brightness)
		if err != nil {
			return "", "", "", err
		}
	}

	if colorTemp != 0 {
		if colorTemp < 1000 || colorTemp > 10000 {
			return "", "", "", fmt.Errorf("color temperature must be 1000-10000 K")
		}
		colorTemp = kelvinToMired(colorTemp)
	}

	return device, "on", domain, h.turn("on", domain, device, brightnessScaled, rgb, colorTemp)
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

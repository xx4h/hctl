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

func (h *Hass) VolumeSet(obj string, volume int) (string, string, string, error) {
	svc := "volume_set"
	sub, obj, err := h.entityArgHandler([]string{obj}, svc)
	if err != nil {
		return "", "", "", err
	}
	hasDomain, err := h.hasDomainWithService(sub, svc)
	if err != nil {
		return "", "", "", err
	} else if !hasDomain {
		return "", "", "", fmt.Errorf("No such Domain with Service: %s with %s", sub, svc)
	}
	if !h.hasEntityInDomain(obj, sub) {
		return "", "", "", fmt.Errorf("No such Entity in Domain: %s in %s", obj, sub)
	}
	payload := map[string]any{"entity_id": fmt.Sprintf("%s.%s", sub, obj), "volume_level": fmt.Sprintf("%.2f", float32(volume)/100)}
	res, err := h.api("POST", fmt.Sprintf("/services/%s/%s", sub, svc), payload)
	if err != nil {
		return "", "", "", err
	}

	if err := h.getResult(res); err != nil {
		return "", "", "", err
	}

	return obj, fmt.Sprintf("volume set to %d%%", volume), sub, nil
}

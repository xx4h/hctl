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

func (h *Hass) play(mediaURL string, mediaType string, sub string, obj string) error {
	payload := map[string]any{
		"entity_id":          fmt.Sprintf("%s.%s", sub, obj),
		"media_content_id":   mediaURL,
		"media_content_type": mediaType,
	}

	res, err := h.api("POST", fmt.Sprintf("/services/%s/play_media", sub), payload)
	if err != nil {
		return err
	}

	if err := h.getResult(res); err != nil {
		return err
	}
	return nil
}

func (h *Hass) PlayMusic(obj string, mediaURL string, name string) (string, string, string, error) {
	sub, obj, err := h.entityArgHandler([]string{obj}, "play_media")
	if err != nil {
		return "", "", "", err
	}
	return obj, fmt.Sprintf("playing %s", name), sub, h.play(mediaURL, "music", sub, obj)
}

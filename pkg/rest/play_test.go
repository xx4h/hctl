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
	"testing"
)

func Test_PlayMusic(t *testing.T) {
	ms := mockServer(t)
	h := &Hass{
		APIURL: ms.URL,
		Token:  "test_token",
	}
	defer ms.Close()
	obj, action, sub, err := h.PlayMusic("player1", "testdata/fake.mp3", "fake.mp3")
	if err != nil {
		t.Errorf("Error playing: %v", err)
	}
	if obj != "player1" {
		t.Errorf("got %s, want player1", obj)
	}
	if action != "playing fake.mp3" {
		t.Errorf("got %s, want 'playing fake.mp3'", action)
	}
	if sub != "media_player" {
		t.Errorf("got %s, want media_player", sub)
	}
}

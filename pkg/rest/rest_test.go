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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xx4h/hctl/pkg/hctltest"
)

func Test_findEntity_Errors(t *testing.T) {
	ms := hctltest.MockServer(t)
	defer ms.Close()
	h := &Hass{
		APIURL: ms.URL,
		Token:  "test_token",
	}

	tests := map[string]struct {
		name    string
		domain  string
		service string
		wantErr string
	}{
		"entity does not exist": {
			name:    "nonexisting",
			domain:  "",
			service: "toggle",
			wantErr: "entity nonexisting does not exist",
		},
		"qualified entity does not exist": {
			name:    "nonexisting",
			domain:  "light",
			service: "toggle",
			wantErr: "entity light.nonexisting does not exist",
		},
		"entity exists but does not support service": {
			name:    "bedroom_main",
			domain:  "",
			service: "speak",
			wantErr: "entity bedroom_main exists but does not support speak",
		},
		"domain has no such service": {
			name:    "bedroom_main",
			domain:  "light",
			service: "speak",
			wantErr: "domain light has no service speak",
		},
		"happy path - no error": {
			name:    "bedroom_main",
			domain:  "",
			service: "toggle",
			wantErr: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d, n, err := h.findEntity(tt.name, tt.domain, tt.service)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if d != "light" || n != "bedroom_main" {
					t.Errorf("got %s.%s, want light.bedroom_main", d, n)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("got error %q, want %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func Test_findEntity_FuzzyErrorResolution(t *testing.T) {
	ms := hctltest.MockServer(t)
	defer ms.Close()
	h := &Hass{
		APIURL: ms.URL,
		Token:  "test_token",
		Fuzz:   true,
	}

	tests := map[string]struct {
		name    string
		domain  string
		service string
		wantErr string
	}{
		"fuzzy resolves name before reporting unsupported service": {
			name:    "bedroom_mai",
			domain:  "",
			service: "speak",
			wantErr: "entity bedroom_main exists but does not support speak",
		},
		"fuzzy resolves name before reporting nonexistent domain service": {
			name:    "bedroom_mai",
			domain:  "light",
			service: "speak",
			wantErr: "domain light has no service speak",
		},
		"fuzzy match not possible reports original name": {
			name:    "zzzzzzzzz",
			domain:  "",
			service: "speak",
			wantErr: "entity zzzzzzzzz does not exist",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, _, err := h.findEntity(tt.name, tt.domain, tt.service)
			if err == nil {
				t.Errorf("expected error %q, got nil", tt.wantErr)
			} else if err.Error() != tt.wantErr {
				t.Errorf("got error %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func Test_api_HTTPErrors(t *testing.T) {
	tests := map[string]struct {
		statusCode int
		wantErr    string
	}{
		"unauthorized": {
			statusCode: http.StatusUnauthorized,
			wantErr:    "authentication failed: invalid or expired token",
		},
		"not found": {
			statusCode: http.StatusNotFound,
			wantErr:    "API endpoint not found (404): /test",
		},
		"server error": {
			statusCode: http.StatusInternalServerError,
			wantErr:    "unexpected API response code 500 for /test",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				if _, err := w.Write([]byte(`{"message": "error"}`)); err != nil {
					t.Errorf("Error writing data: %v", err)
				}
			}))
			defer ms.Close()

			h := &Hass{
				APIURL: ms.URL,
				Token:  "test_token",
			}

			_, err := h.api("GET", "/test", nil)
			if err == nil {
				t.Errorf("expected error %q, got nil", tt.wantErr)
			} else if err.Error() != tt.wantErr {
				t.Errorf("got error %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

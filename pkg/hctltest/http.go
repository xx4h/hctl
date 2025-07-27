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

package hctltest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func MockServer(t testing.TB) *httptest.Server {
	_, filename, _, _ := runtime.Caller(0)
	testdir := filepath.Dir(filename)
	t.Helper()
	mux := http.NewServeMux()

	// handle default entry point and return API running
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "API running."}`)); err != nil {
			t.Errorf("Error writing data: %v", err)
		}
	})

	// get all services
	mux.HandleFunc("/services", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		data, err := os.ReadFile(fmt.Sprintf("%s/testdata/services.json", testdir))
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			t.Errorf("Error writing data: %v", err)
		}
	})

	// get all states
	mux.HandleFunc("/states", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		data, err := os.ReadFile(fmt.Sprintf("%s/testdata/states.json", testdir))
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			t.Errorf("Error writing data: %v", err)
		}
	})

	// Update entity in domain/service
	//nolint:govet
	mux.HandleFunc("POST /services/{domain}/{service}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading body: %v", err)
		}
		m := make(map[string]any)
		if err := json.Unmarshal(body, &m); err != nil {
			t.Errorf("Error Unmarshal: %v", err)
		}
		l := strings.Split(m["entity_id"].(string), ".")
		name := l[1]
		data, err := os.ReadFile(fmt.Sprintf("%s/testdata/%s_%s_%s_response.json", testdir, name, r.PathValue("domain"), r.PathValue("service")))
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			t.Errorf("Error writing data: %v", err)
		}
	})
	mockServer := httptest.NewServer(mux)
	return mockServer
}

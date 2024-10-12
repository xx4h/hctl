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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func mockServerGetDataFromFile(t testing.TB, testdatafile string) *httptest.Server {
	t.Helper()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		data, err := os.ReadFile(fmt.Sprintf("testdata/%s", testdatafile))
		if err != nil {
			t.Errorf("Error reading file: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			t.Errorf("Error writing data: %v", err)
		}
	}))
	return mockServer
}

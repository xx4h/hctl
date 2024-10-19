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

package cmd

import (
	"bytes"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/xx4h/hctl/pkg"
)

type cmdTest struct {
	input  string
	rexOut string
	rexErr string
}

func newTestingHctl(t *testing.T) *pkg.Hctl {
	t.Helper()
	h, err := pkg.NewHctl(true)
	if err != nil {
		t.Error(err)
	}

	// create tempdir for config
	tmpDir := t.TempDir()
	testdata, err := os.Open("testdata/hctl.yaml")
	if err != nil {
		t.Error(err)
	}
	config, err := os.Create(path.Join(tmpDir, "hctl.yaml"))
	if err != nil {
		t.Error(err)
	}
	if _, err := io.Copy(config, testdata); err != nil {
		t.Error(err)
	}

	if err := h.LoadConfig(path.Join(tmpDir, "hctl.yaml")); err != nil {
		t.Error(err)
	}
	return h
}

func testCmd(t *testing.T, h *pkg.Hctl, tests map[string]cmdTest) {
	t.Helper()
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			out := new(bytes.Buffer)
			errout := new(bytes.Buffer)
			in := strings.Split(tt.input, " ")
			rootCmd = newRootCmd(h, out, in)
			rootCmd.SetOut(out)
			rootCmd.SetErr(errout)
			rootCmd.SetArgs(in)
			if err := rootCmd.Execute(); err != nil {
				t.Error(err)
			}
			e := errout.String()
			o := out.String()
			rex, err := regexp.Compile(tt.rexErr)
			if err != nil {
				t.Error(err)
			}
			if ok := rex.MatchString(e); !ok {
				t.Errorf("got %s, want %s", o, tt.rexErr)
			}
			ok, err := regexp.MatchString(tt.rexOut, o)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("got %s, want %s", o, tt.rexOut)
			}
		})
	}
}

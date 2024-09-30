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

package main

import (
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/rs/zerolog/log"

	"github.com/xx4h/hctl/pkg"
)

// initCmd represents the init command
func newInitCmd(h *pkg.Hctl) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize and config hctl",
		Run: func(_ *cobra.Command, _ []string) {
			// TODO: find better way to get this path (needs to work together with cmd/root.go)
			userDir, err := os.UserHomeDir()
			if err != nil {
				log.Error().Msgf("Error: %v", err)
			}
			h.InitializeConfig(path.Join(userDir, ".config/hctl/hctl.yaml"))
		},
	}

	return cmd
}

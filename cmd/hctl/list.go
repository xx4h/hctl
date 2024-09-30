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
	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
)

// listCmd represents the list command
func newListCmd(h *pkg.Hctl) *cobra.Command {

	var domains []string
	var services []string

	cmd := &cobra.Command{
		Use:       "list [entities|services]",
		Short:     "List all existing entities or services",
		Aliases:   []string{"l"},
		ValidArgs: []string{"entities", "services"},
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "entities" {
				h.DumpStates(domains)
			} else if args[0] == "services" {
				h.DumpServices(domains, services)
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringArrayVarP(&domains, "domains", "d", []string{}, "Limit domains")
	cmd.PersistentFlags().StringArrayVarP(&services, "services", "s", []string{}, "Limit services")

	return cmd
}

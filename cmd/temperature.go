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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

func newTemperatureCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "temperature",
		Short:   "Set the temperature of a climate entity",
		Aliases: []string{"te", "temp"}, // codespell:ignore
		Args:    cobra.MatchAll(cobra.MinimumNArgs(2)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return compListStates(toComplete, args, []string{"set_temperature"}, nil, "", h)
			}
			// Offer devices for subsequent args (value is free-form, no completion needed)
			return compListStatesMulti(toComplete, args, []string{"set_temperature"}, nil, "", h)
		},
		Run: func(_ *cobra.Command, args []string) {
			value := args[len(args)-1]
			devices := args[:len(args)-1]
			var hasErr bool
			for _, device := range devices {
				obj, state, err := h.TemperatureSet(device, value)
				if err != nil {
					o.FprintErrorMsg(out, err)
					hasErr = true
				} else {
					o.FprintSuccessAction(out, obj, state)
				}
			}
			if hasErr {
				os.Exit(1)
			}
		},
	}

	return cmd
}

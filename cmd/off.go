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

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

func newOffCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	var transition float64

	cmd := &cobra.Command{
		Use:   "off [--transition seconds]",
		Short: "Switch or turn off a light or switch",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return compListStatesMulti(toComplete, args, []string{"turn_off"}, nil, "on", h)
		},
		Run: func(_ *cobra.Command, args []string) {
			c := h.GetRest()
			hasTransition := transition != 0
			var hasErr bool
			for _, device := range args {
				var obj, state, sub string
				var err error
				if hasTransition {
					obj, state, sub, err = c.TurnLightOffTransition(device, transition)
				} else {
					obj, state, sub, err = c.TurnOff(device)
				}
				if err != nil {
					o.FprintErrorMsg(out, err)
					hasErr = true
				} else {
					o.FprintSuccessAction(out, obj, state)
				}
				log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
			}
			if hasErr {
				os.Exit(1)
			}
		},
	}

	cmd.PersistentFlags().Float64VarP(&transition, "transition", "s", 0, "Set transition time in seconds (e.g. 1.5)")

	return cmd
}

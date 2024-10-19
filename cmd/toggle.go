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

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

// toggleCmd represents the toggle command
func newToggleCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "toggle",
		Short:   "Toggle on/off a light or switch",
		Aliases: []string{"t"},
		Args:    cobra.MatchAll(cobra.ExactArgs(1)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return noMoreArgsComp()
			}
			return compListStates(toComplete, args, []string{"toggle"}, nil, "", h)
		},
		Run: func(_ *cobra.Command, args []string) {
			c := h.GetRest()
			obj, state, sub, err := c.Toggle(args[0])
			if err != nil {
				o.FprintError(out, err)
			} else {
				o.FprintSuccessAction(out, obj, state)
			}
			log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
		},
	}

	return cmd
}

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
	"fmt"
	"io"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
	"github.com/xx4h/hctl/pkg/util"
)

var (
	brightnessRange = append([]string{"+", "-", "min", "mid", "max"}, util.MakeRangeString(1, 99)...)
)

func newBrightnessCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "brightness [+|-|min|max|1-99]",
		Short:   "Change brightness",
		Aliases: []string{"b", "br", "bright"},
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return compListStates(toComplete, args, []string{"turn_on", "turn_off"}, []string{"brightness"}, "", h)
			} else if len(args) == 1 {
				return brightnessRange, cobra.ShellCompDirectiveKeepOrder | cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateBrightness(args[1]); err != nil {
				return err
			}
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			c := h.GetRest()
			obj, state, sub, err := c.TurnLightOnBrightness(args[0], args[1])
			if err != nil {
				o.FprintError(out, err)
			} else {
				o.FprintSuccess(out, fmt.Sprintf("%s brightness set to %s%%", obj, args[1]))
			}
			log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
		},
	}

	return cmd
}

func validateBrightness(brightness string) error {
	if !slices.Contains(brightnessRange, brightness) {
		return fmt.Errorf("brightness needs to be 1-99, or +/-/min/max")
	}
	return nil
}

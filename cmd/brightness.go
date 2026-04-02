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
	"os"
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
		Args:    cobra.MatchAll(cobra.MinimumNArgs(2)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return compListStates(toComplete, args, []string{"turn_on", "turn_off"}, []string{"brightness"}, "", h)
			}
			// Offer both devices and brightness values for subsequent args
			devices, directive := compListStatesMulti(toComplete, args, []string{"turn_on", "turn_off"}, []string{"brightness"}, "", h)
			return append(devices, brightnessRange...), directive
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateBrightness(args[len(args)-1]); err != nil {
				return err
			}
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			c := h.GetRest()
			value := args[len(args)-1]
			devices := args[:len(args)-1]
			var hasErr bool
			for _, device := range devices {
				obj, state, sub, err := c.TurnLightOnCustom(device, value, "", 0, 0)
				if err != nil {
					o.FprintErrorMsg(out, err)
					hasErr = true
				} else {
					o.FprintSuccess(out, fmt.Sprintf("%s brightness set to %s%%", obj, value))
				}
				log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
			}
			if hasErr {
				os.Exit(1)
			}
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

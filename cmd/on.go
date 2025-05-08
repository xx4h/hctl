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

func newOnCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	var brightness string
	var color string
	var colorTemp int

	cmd := &cobra.Command{
		Use:   "on [-b|--brightness +|-|min|max|1-99] [-c|--color R,G,B] [-t|--color_temp 153-500]",
		Short: "Switch or turn on a light or switch",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return noMoreArgsComp()
			}
			return compListStates(toComplete, args, []string{"turn_on"}, nil, "off", h)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if brightness == "" {
				return nil
			}
			if err := validateBrightness(brightness); err != nil {
				return err
			}
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			c := h.GetRest()
			var obj, state, sub string
			var err error
			if brightness != "" || color != "" || colorTemp != 0 {
				obj, state, sub, err = c.TurnLightOnCustom(args[0], brightness, color, colorTemp)
			} else {
				obj, state, sub, err = c.TurnOn(args[0])
			}
			if err != nil {
				o.FprintError(out, err)
			} else {
				o.FprintSuccessAction(out, obj, state)
			}
			log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
		},
	}

	cmd.PersistentFlags().StringVarP(&brightness, "brightness", "b", "", "Set brightness")
	cmd.PersistentFlags().StringVarP(&color, "color", "c", "", "Set RGB color in format R,G,B")
	cmd.PersistentFlags().IntVarP(&colorTemp, "color_temp", "t", 0, "Set color temperature in mireds (153-500)")
	err := cmd.RegisterFlagCompletionFunc("brightness", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return brightnessRange, cobra.ShellCompDirectiveKeepOrder | cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		log.Error().Msgf("Could not register flag completion func for brightness: %+v", err)
	}

	return cmd
}

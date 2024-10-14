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

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
	"github.com/xx4h/hctl/pkg/util"
)

// toggleCmd represents the toggle command
func newVolumeCmd(h *pkg.Hctl, _ io.Writer) *cobra.Command {
	volRange := util.MakeRangeString(0, 100)
	cmd := &cobra.Command{
		Use:     "volume",
		Short:   "Set volume of e.g media player",
		Aliases: []string{"v"},
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return compListStates(toComplete, args, []string{"volume_set"}, nil, "", h)
			} else if len(args) == 1 {
				return volRange, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		Run: func(_ *cobra.Command, args []string) {
			if !slices.Contains(volRange, args[1]) {
				o.PrintError(fmt.Errorf("volume needs to be 1-100"))
			}
			h.VolumeSet(args[0], args[1])
		},
	}

	return cmd
}

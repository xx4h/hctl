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

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
)

// toggleCmd represents the toggle command
func newPlayCmd(h *pkg.Hctl, _ io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "play",
		Short:   "Play music from url on media player",
		Aliases: []string{"p"},
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 1 {
				return compMediaMap(toComplete, args, h)
			} else if len(args) == 0 {
				return compListStates(toComplete, args, []string{"play_media"}, nil, "", h)
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(_ *cobra.Command, args []string) {
			h.PlayMusic(args[0], args[1])
		},
	}

	return cmd
}

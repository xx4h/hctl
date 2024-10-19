// Copyright 2024 Fabian `xx4h` Sylvester
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

const (
	// editorconfig-checker-disable
	configRemExample = `
  # Remove config option
  hctl config rem device_map.a
  hctl config remove device_map.b
  `
	// editorconfig-checker-enable
)

func newConfigRemCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove PATH",
		Short:   "Set config variables",
		Aliases: []string{"r", "re", "rem", "remo"},
		Example: configRemExample,
		Args:    cobra.MatchAll(cobra.ExactArgs(1)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return compListConfig(toComplete, args, h)
		},
		Run: func(_ *cobra.Command, args []string) {
			if err := h.RemoveConfigOptionWrite(args[0]); err != nil {
				o.FprintError(out, err)
			}
			o.FprintSuccess(out, fmt.Sprintf("Option `%s` successfully removed.", args[0]))
		},
	}

	return cmd
}

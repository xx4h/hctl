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

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

const (
	// editorconfig-checker-disable
	configSetExample = `
  # Set config option
  hctl config set logging.log_level debug
  `
	// editorconfig-checker-enable
)

func newConfigSetCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set PATH VALUE",
		Short:   "Set config variables",
		Aliases: []string{"s", "se"},
		Example: configSetExample,
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}
			return compListConfig(toComplete, args, h)
		},
		Run: func(_ *cobra.Command, args []string) {
			if err := h.SetConfigValueWrite(args[0], args[1]); err != nil {
				o.FprintError(out, err)
			}
			o.FprintSuccess(out, fmt.Sprintf("Option `%s` successfully set to `%s`.", args[0], args[1]))
		},
	}

	return cmd
}

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
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

const (
	// editorconfig-checker-disable
	configGetExample = `
  # Get all config options
  hctl config get

  # Get all settings from a section
  hctl config get logging

  # Get a specific config option
  hctl config get hub.url
  `
	// editorconfig-checker-enable
)

func configGetAll(h *pkg.Hctl) [][]interface{} {
	a, _ := compListConfig("", []string{}, h)
	slices.Sort(a)
	var clist [][]interface{}
	for _, b := range a {
		v, err := h.GetConfigValue(b)
		if err == nil {
			clist = append(clist, []any{b, v})
		}
	}
	return clist
}

func configGetSection(h *pkg.Hctl, section string) [][]interface{} {
	prefix := section + "."
	var clist [][]interface{}
	// Use all config entries to find matching section paths
	all := configGetAll(h)
	for _, entry := range all {
		if path, ok := entry[0].(string); ok && strings.HasPrefix(path, prefix) {
			clist = append(clist, entry)
		}
	}
	return clist
}

func newConfigGetCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get [PATH]",
		Short:   "Get configuration parameters",
		Aliases: []string{"g", "ge"},
		Example: configGetExample,
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return noMoreArgsComp()
			}
			return compListConfigWithSections(toComplete, args, h)
		},
		Run: func(_ *cobra.Command, args []string) {
			var header = append([]any{}, "OPTION", "VALUE")
			var clist [][]interface{}
			if len(args) == 0 {
				clist = configGetAll(h)
			} else {
				// Try as section prefix first
				clist = configGetSection(h, args[0])
				// If no section matches, try as exact leaf path
				if len(clist) == 0 {
					v, err := h.GetConfigValue(args[0])
					if err != nil {
						o.FprintError(out, err)
					}
					clist = append(clist, []any{args[0], v})
				}
			}
			o.FprintSuccessListWithHeader(out, header, clist)
		},
	}

	return cmd
}

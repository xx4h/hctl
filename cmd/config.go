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

func newConfigCmd(h *pkg.Hctl, out io.Writer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Set and get config",
		Aliases: []string{"c", "co", "con", "conf"},
		Example: configGetExample + configSetExample,
	}

	cmd.AddCommand(
		newConfigGetCmd(h, out),
		newConfigSetCmd(h, out),
		newConfigRemCmd(h, out),
	)

	return cmd
}

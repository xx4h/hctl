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

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	o "github.com/xx4h/hctl/pkg/output"
)

// versionCmd represents the version command
func newVersionCmd(out io.Writer) *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Run: func(_ *cobra.Command, _ []string) {
			printVersion(out, short)
		},
	}

	cmd.PersistentFlags().BoolVarP(&short, "short", "s", false, "Prints version info in short format")

	return cmd
}

func printVersion(out io.Writer, short bool) {
	if !short {
		banner, err := o.GetBanner()
		if err != nil {
			log.Error().Msgf("Could not get banner: %v", err)
		}
		fmt.Fprint(out, banner)

	}
	const format = "%-10s %s\n"
	fmt.Fprintf(out, format, "Version:", version)
	fmt.Fprintf(out, format, "Commit:", commit)
	fmt.Fprintf(out, format, "Date:", date)
}

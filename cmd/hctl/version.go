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

package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
)

// versionCmd represents the version command
func newVersionCmd(hctl *pkg.Hctl, out io.Writer) *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion(out, short)
		},
	}

	cmd.PersistentFlags().BoolVarP(&short, "short", "s", false, "Prints version info in short format")

	return cmd
}

func printVersion(out io.Writer, short bool) {
	const format = "%-10s %s\n"
	fmt.Fprintf(out, format, "Version:", version)
	fmt.Fprintf(out, format, "Commit:", commit)
	fmt.Fprintf(out, format, "Date:", date)
}

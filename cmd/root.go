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
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
	o "github.com/xx4h/hctl/pkg/output"
)

var (
	appName               = "hctl"
	version, commit, date = "dev", "dev", "1970-01-01"
	rootCmd               *cobra.Command
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd(h *pkg.Hctl, out io.Writer, _ []string) *cobra.Command {
	var logLevel string

	banner, err := o.GetBanner()
	if err != nil {
		log.Warn().Msgf("Could not render banner: %v", err)
		banner = ""
	}

	cmd := &cobra.Command{
		Use:   appName,
		Short: "A command line tool to control your home automation",
		Long:  fmt.Sprintf("%s\nHctl is a CLI tool to control your home automation", banner),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if logLevel != "" {
				lvl, err := zerolog.ParseLevel(logLevel)
				if err != nil {
					return err
				}
				zerolog.SetGlobalLevel(lvl)
			}
			return nil
		},
	}

	cmd.AddCommand(
		newBrightnessCmd(h, out),
		newCompletionCmd(),
		newConfigCmd(h, out),
		newInitCmd(h),
		newListCmd(h, out),
		newOffCmd(h, out),
		newOnCmd(h, out),
		newPlayCmd(h, out),
		newToggleCmd(h, out),
		newVersionCmd(out),
		newVolumeCmd(h, out),
		newTemperatureCmd(h, out),
	)

	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "Set the log level")

	return cmd
}

func RunCmd() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	log.Logger = log.Output(output)
	h, err := pkg.NewHctl(false)
	if err != nil {
		log.Fatal().Caller().Msgf("Error: %v", err)
	}
	rootCmd = newRootCmd(h, os.Stdout, os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msgf("Error: %v", err)
	}
}

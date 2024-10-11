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
	"os"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/xx4h/hctl/pkg"
)

func newCompletionCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: fmt.Sprintf(`To load completions:

  Bash:

    $ source <(%[1]s completion bash)

    # To load completions for each session, execute once:
    # Linux:
    $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
    # macOS:
    $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

  Zsh:

    # If shell completion is not already enabled in your environment,
    # you will need to enable it.  You can execute the following once:

    $ echo "autoload -U compinit; compinit" >> ~/.zshrc

    # To load completions for each session, execute once:
    $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

    # You will need to start a new shell for this setup to take effect.

  fish:

    $ %[1]s completion fish | source

    # To load completions for each session, execute once:
    $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

  PowerShell:

    PS> %[1]s completion powershell | Out-String | Invoke-Expression

    # To load completions for every new session, run:
    PS> %[1]s completion powershell > %[1]s.ps1
    # and source this file from your PowerShell profile.
  `, appName),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				err := cmd.Root().GenBashCompletion(os.Stdout)
				if err != nil {
					log.Fatal().Msg("Could not generate Bash Completion")
				}
			case "zsh":
				err := cmd.Root().GenZshCompletion(os.Stdout)
				if err != nil {
					log.Fatal().Msg("Could not generate Bash Completion")
				}
			case "fish":
				err := cmd.Root().GenFishCompletion(os.Stdout, true)
				if err != nil {
					log.Fatal().Msg("Could not generate Bash Completion")
				}
			case "powershell":
				err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
				if err != nil {
					log.Fatal().Msg("Could not generate Bash Completion")
				}
			}
		},
	}

	return cmd
}

// noMoreArgsCompFunc deactivates file completion when doing argument shell completion.
// It also provides some ActiveHelp to indicate no more arguments are accepted.
// func noMoreArgsCompFunc(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
// 	return noMoreArgsComp()
// }

// noMoreArgsComp deactivates file completion when doing argument shell completion.
// It also provides some ActiveHelp to indicate no more arguments are accepted.
func noMoreArgsComp() ([]string, cobra.ShellCompDirective) {
	activeHelpMsg := "This command does not take any more arguments (but may accept flags)."
	return cobra.AppendActiveHelp(nil, activeHelpMsg), cobra.ShellCompDirectiveNoFileComp
}

// support function for completion
func compListStates(_ string, ignoredStates []string, service string, state string, h *pkg.Hctl) ([]string, cobra.ShellCompDirective) {
	states, err := h.GetStates()
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
	}
	services, err := h.GetServices()
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
	}

	filteredStates := filterStates(states, ignoredStates)
	filteredStates = filterCapable(filteredStates, services, service, state)

	var choices []string
	for _, rel := range filteredStates {
		if h.CompletionShortNamesEnabled() {
			s := strings.Split(rel.EntityID, ".")
			if slices.Contains(choices, s[1]) {
				// if we've more than one, we add the second-n with long name
				choices = append(choices, rel.EntityID)
			} else {
				choices = append(choices, s[1])
			}
		} else {
			choices = append(choices, rel.EntityID)
		}
	}

	return choices, cobra.ShellCompDirectiveNoFileComp
}

func compListConfig(_ string, _ []string, h *pkg.Hctl) ([]string, cobra.ShellCompDirective) {
	return h.GetConfigOptionsAsPaths(), cobra.ShellCompDirectiveNoFileComp
}

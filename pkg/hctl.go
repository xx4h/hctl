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

package pkg

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/xx4h/hctl/pkg/config"
	i "github.com/xx4h/hctl/pkg/init"
	o "github.com/xx4h/hctl/pkg/output"
	"github.com/xx4h/hctl/pkg/rest"
)

type Hctl struct {
	cfg *config.Config
	//out io.ReadWriteCloser
	//log *zerolog.Logger
}

func NewHctl() (*Hctl, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	return &Hctl{
		cfg: cfg,
	}, nil
}

func (h *Hctl) InitializeConfig(path string) {
	if err := i.InitializeConfig(h.cfg, path); err != nil {
		o.PrintError(err)
	} else {
		o.PrintSuccess(fmt.Sprintf("Successfully created config: %s\n", path))
	}
}

func (h *Hctl) CompletionShortNamesEnabled() bool {
	return h.cfg.Completion.ShortNames
}

func (h *Hctl) GetRest() *rest.Hass {
	return rest.New(h.cfg.Hub.Url, h.cfg.Hub.Token, h.cfg.Handling.Fuzz)
}

func (h *Hctl) GetServices() []rest.HassService {
	services, err := h.GetRest().GetServices()
	if err != nil {
		log.Fatal().Msgf("Error: %v", err)
		return nil
	}
	return services
}

func (h *Hctl) GetStates() []rest.HassState {
	states := h.GetRest().GetStates()
	return states
}

func (h *Hctl) GetFilteredServices(domains []string, services []string) []rest.HassService {
	s, err := h.GetRest().GetFilteredServices(domains, services)
	if err != nil {
		log.Fatal().Msgf("Error: %v", err)
		return nil
	}
	return s
}

func (h *Hctl) GetFilteredStates(domains []string) []rest.HassState {
	return h.GetRest().GetFilteredStates(domains)
}

func (h *Hctl) GetFilteredServicesMap(domains []string, services []string) map[string][]string {
	t, err := h.GetRest().GetFilteredServicesMap(domains, services)
	if err != nil {
		log.Fatal().Msgf("Error: %v", err)
		return nil
	}
	return t
}

func (h *Hctl) GetFilteredStatesMap(domains []string) map[string][]string {
	return h.GetRest().GetFilteredStatesMap(domains)
}

func (h *Hctl) DumpServices(domains []string, services []string) {
	t := h.GetFilteredServicesMap(domains, services)
	if err := o.PrintThreeLevelFlatTree("Services", t); err != nil {
		log.Error().Msgf("Error: %v", err)
	}
}

func (h *Hctl) DumpStates(domains []string) {
	t := h.GetFilteredStatesMap(domains)
	if err := o.PrintThreeLevelFlatTree("States", t); err != nil {
		log.Error().Msgf("Error: %v", err)
	}
}

func (h *Hctl) SetLogging(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)
	return nil
}

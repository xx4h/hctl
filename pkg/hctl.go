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
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/xx4h/hctl/pkg/config"
	i "github.com/xx4h/hctl/pkg/init"
	o "github.com/xx4h/hctl/pkg/output"
	"github.com/xx4h/hctl/pkg/rest"
	"github.com/xx4h/hctl/pkg/serve"
	"github.com/xx4h/hctl/pkg/util"
)

type Hctl struct {
	cfg *config.Config
	// out io.ReadWriteCloser
	// log *zerolog.Logger
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

func (h *Hctl) GetConfigValue(p string) any {
	v, err := h.cfg.GetValueByPath(p)
	if err != nil {
		o.PrintError(err)
	}
	return v
}

func (h *Hctl) SetConfigValue(p string, v string) error {
	err := h.cfg.SetValueByPath(p, v)
	return err
}

func (h *Hctl) GetConfigOptionsAsPaths() []string {
	return h.cfg.GetOptionsAsPaths()
}

func (h *Hctl) GetRest() *rest.Hass {
	return rest.New(h.cfg.Hub.URL, h.cfg.Hub.Token, h.cfg.Handling.Fuzz)
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

func (h *Hctl) PlayMusic(obj string, mediaURL string) {
	// handle url or file system path
	if ok := util.IsURL(mediaURL); ok {
		// if we already have a url, just play it

		obj, state, sub, err := h.GetRest().PlayMusic(obj, mediaURL, mediaURL)
		if err != nil {
			o.PrintError(err)
		}

		o.PrintSuccessAction(obj, state)
		log.Debug().Msgf("Result: %s(%s) to %s", obj, sub, state)

	} else {
		// if we don't have a url but a filepath

		// check if file exists
		_, err := os.Stat(mediaURL)
		if err != nil {
			o.PrintError(err)
		}

		// get new Media instance
		s := serve.NewMedia(h.cfg.GetServeIP(), h.cfg.GetServePort(), mediaURL)
		// start instance and wait until ready
		s.FileHandler()
		if err := s.WaitForHTTPReady(); err != nil {
			log.Fatal().Msgf("HTTP server ready error: %v", err)
		}

		// we are ready and send the url to play
		obj, state, sub, err := h.GetRest().PlayMusic(obj, s.GetURL(), s.GetMediaName())
		if err != nil {
			o.PrintError(err)
		}

		o.PrintSuccessAction(obj, state)
		log.Debug().Msgf("Result: %s(%s) to %s", obj, sub, state)
		s.WaitAndClose()
	}
}

func (h *Hctl) VolumeSet(obj string, volume string) {
	vint, err := strconv.Atoi(volume)
	if err != nil {
		o.PrintError(err)
	}
	obj, state, sub, err := h.GetRest().VolumeSet(obj, vint)
	if err != nil {
		o.PrintError(err)
	}
	o.PrintSuccessAction(obj, state)
	log.Debug().Msgf("Result: %s(%s) to %s", obj, sub, state)
}

func (h *Hctl) SetLogging(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)
	return nil
}

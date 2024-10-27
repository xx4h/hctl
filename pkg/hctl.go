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
	"io"
	"os"
	"strconv"
	"time"

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

func NewHctl(testing bool) (*Hctl, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	if !testing {
		err := cfg.LoadConfig("")
		if err != nil {
			return nil, err
		}
	}

	return &Hctl{
		cfg: cfg,
	}, nil
}

func (h *Hctl) LoadConfig(configPath string) error {
	err := h.cfg.LoadConfig(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hctl) InitializeConfig(path string) {
	if err := i.InitializeConfig(h.cfg, path); err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		o.PrintError(err)
	} else {
		o.PrintSuccess(fmt.Sprintf("Successfully created config: %s\n", path))
	}
}

func (h *Hctl) CompletionShortNamesEnabled() bool {
	return h.cfg.Completion.ShortNames
}

func (h *Hctl) GetConfigValue(p string) (any, error) {
	v, err := h.cfg.GetValueByPath(p)
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		return nil, err
	}
	return v, nil
}

func (h *Hctl) SetConfigValue(p string, v string) error {
	err := h.cfg.SetValueByPath(p, v)
	return err
}

func (h *Hctl) GetMap(k string) map[string]string {
	return h.cfg.Viper.GetStringMapString(k)
}

func (h *Hctl) RemoveConfigOption(p string) error {
	err := h.cfg.RemoveOptionByPath(p)
	return err
}

func (h *Hctl) RemoveConfigOptionWrite(p string) error {
	err := h.cfg.RemoveOptionByPathWrite(p)
	return err
}

func (h *Hctl) SetConfigValueWrite(p string, v string) error {
	err := h.cfg.SetValueByPathWrite(p, v)
	return err
}

func (h *Hctl) GetConfigOptionsAsPaths() []string {
	return h.cfg.GetOptionsAsPaths()
}

func (h *Hctl) GetRest() *rest.Hass {
	return rest.New(h.cfg.Hub.URL, h.cfg.Hub.Token, h.cfg.Handling.Fuzz, h.cfg.DeviceMap)
}

func (h *Hctl) GetServices() ([]rest.HassService, error) {
	services, err := h.GetRest().GetServices()
	if err != nil {
		log.Fatal().Msgf("Error: %+v", err)
		return nil, err
	}
	return services, nil
}

func (h *Hctl) GetStates() ([]rest.HassState, error) {
	states, err := h.GetRest().GetStates()
	if err != nil {
		return nil, err
	}
	return states, nil
}

func (h *Hctl) GetFilteredServices(domains []string, services []string) []rest.HassService {
	s, err := h.GetRest().GetFilteredServices(domains, services)
	if err != nil {
		log.Fatal().Msgf("Error: %+v", err)
		return nil
	}
	return s
}

func (h *Hctl) GetFilteredStates(domains []string) ([]rest.HassState, error) {
	return h.GetRest().GetFilteredStates(domains)
}

func (h *Hctl) GetFilteredServicesMap(domains []string, services []string) map[string][]string {
	t, err := h.GetRest().GetFilteredServicesMap(domains, services)
	if err != nil {
		log.Fatal().Msgf("Error: %+v", err)
		return nil
	}
	return t
}

func (h *Hctl) GetFilteredStatesMap(domains []string) (map[string][]string, error) {
	return h.GetRest().GetFilteredStatesMap(domains)
}

func (h *Hctl) DumpServices(out io.Writer, domains []string, services []string) {
	t := h.GetFilteredServicesMap(domains, services)
	if err := o.PrintThreeLevelFlatTree(out, "Services", t); err != nil {
		log.Error().Msgf("Error: %+v", err)
	}
}

func (h *Hctl) DumpStates(out io.Writer, domains []string) {
	t, err := h.GetFilteredStatesMap(domains)
	if err != nil {
		o.FprintError(out, err)
	}
	if err := o.PrintThreeLevelFlatTree(out, "States", t); err != nil {
		log.Error().Msgf("Error: %+v", err)
	}
}

func (h *Hctl) PlayMusic(out io.Writer, obj string, mediaURL string) {
	if mapURL, ok := h.cfg.MediaMap[mediaURL]; ok {
		mediaURL = mapURL
	}

	// handle url or file system path
	if ok := util.IsURL(mediaURL); ok {
		// if we already have a url, just play it

		obj, state, sub, err := h.GetRest().PlayMusic(obj, mediaURL, mediaURL)
		if err != nil {
			log.Debug().Caller().Msgf("Error: %+v", err)
			o.FprintError(out, err)
		}

		o.PrintSuccessAction(obj, state)
		log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)

	} else {
		// if we don't have a url but a filepath

		// check if file exists
		_, err := os.Stat(mediaURL)
		if err != nil {
			log.Debug().Caller().Msgf("Error: %+v", err)
			o.FprintError(out, err)
		}

		// get new Media instance
		s := serve.NewMedia(h.cfg.GetServeIP(), h.cfg.GetServePort(), mediaURL)
		// start instance and wait until ready
		s.FileHandler()
		if err := s.WaitForHTTPReady(); err != nil {
			log.Fatal().Msgf("HTTP server ready error: %+v", err)
		}

		// we are ready and send the url to play
		obj, state, sub, err := h.GetRest().PlayMusic(obj, s.GetURL(), s.GetMediaName())
		if err != nil {
			log.Debug().Caller().Msgf("Error: %+v", err)
			o.FprintError(out, err)
		}

		o.FprintSuccessAction(out, obj, state)
		log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
		// TODO: find better way to ensure we don't close the server before file has been served
		// -> RaceCondition
		// Problem: We could wait for initial connection (e.g. from media player) and only then
		// move on to WaitAndClose, but media player is not always (i !guess! it depends on time
		// between plays and filesize) re-requesting the file when we play one and the same media
		// twice in a short period of time
		time.Sleep(100 * time.Millisecond)
		s.WaitAndClose()
	}
}

func (h *Hctl) VolumeSet(obj string, volume string) (string, string, error) {
	vint, err := strconv.Atoi(volume)
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		return "", "", err
	}
	obj, state, sub, err := h.GetRest().VolumeSet(obj, vint)
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		return "", "", err
	}
	log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
	return obj, state, nil
}

func (h *Hctl) TemperatureSet(obj string, temp string) (string, string, error) {
	tint, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		return "", "", err
	}
	obj, state, sub, err := h.GetRest().TemperatureSet(obj, tint)
	if err != nil {
		log.Debug().Caller().Msgf("Error: %+v", err)
		return "", "", err
	}
	log.Debug().Caller().Msgf("Result: %s(%s) to %s", obj, sub, state)
	return obj, state, nil
}

func (h *Hctl) SetLogging(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(lvl)
	return nil
}

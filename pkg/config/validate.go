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

package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func validateIsSection(p []string) error {
	if len(p) == 1 {
		return fmt.Errorf("cannot set value for section: %s", p[0])
	}
	return nil
}

func validateNonEmptyKey(p []string) error {
	if p[len(p)-1] == "" {
		return fmt.Errorf("cannot use empty key in %s", p[0])
	}
	return nil
}

func validateLoggingString(value string) error {
	_, err := zerolog.ParseLevel(value)
	if err != nil {
		return fmt.Errorf("unknown log_level: %s (Supported: trace, debug, error, warn, info)", value)
	}
	return nil
}

func validateSetMediaMap(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("media_map value needs to be string")
	}
	if strings.HasPrefix(s, "~") {
		return fmt.Errorf("media_map does not support tilde path expansion yet")
	}
	return nil
}

func validateSetDeviceMap(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("device_map value needs to be string")
	}
	return nil
}

func validateSetLogging(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("device_map value needs to be string")
	}
	opt := path[len(path)-1]
	if opt == "log_level" {
		return validateLoggingString(s)
	}
	return fmt.Errorf("unknown config option for logging: %s", path[len(path)-1])
}

func validateSetHub(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	opt := path[len(path)-1]
	switch opt {
	case "type":
		if value != "hass" {
			return fmt.Errorf("unknown hub type: %s (Supported: hass)", value)
		}
	case "url":
	case "token":
	default:
		return fmt.Errorf("unknown config option for hub: %s", opt)
	}
	return nil
}

func validateSetHandling(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	opt := path[len(path)-1]
	switch opt {
	case "fuzz":
		s := value.(string)
		if _, err := strconv.ParseBool(s); err != nil {
			return fmt.Errorf("Handling fuzz needs to be true/false")
		}
	default:
		return fmt.Errorf("unknown config option for handling: %s", opt)
	}
	return nil
}

func validateSetCompletion(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	opt := path[len(path)-1]
	switch opt {
	case "short_names":
		s := value.(string)
		if _, err := strconv.ParseBool(s); err != nil {
			return fmt.Errorf("Completion short_names needs to be true/false")
		}
	default:
		return fmt.Errorf("unknown config option for completion: %s", opt)
	}
	return nil
}

func validateSetServe(path []string, value any) error {
	log.Debug().Caller().Msgf("Validating set for %s: %+v", path, value)
	opt := path[len(path)-1]
	switch opt {
	case "ip":
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("device_map value needs to be string")
		}
		if ip := net.ParseIP(s); ip == nil {
			return fmt.Errorf("Serve ip option need valid ip address")
		}
	case "port":
		s := value.(string)
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Serve port needs to be a number")
		}
		if port < 1024 {
			return fmt.Errorf("use a non-well-known port (>1023)")
		}
		if port > 65535 {
			return fmt.Errorf("use a valid port in the range 1024-65535")
		}
	default:
		return fmt.Errorf("unknown config option for serve: %s", opt)
	}
	return nil
}

func validateSet(path string, value any) error {
	p := strings.Split(path, ".")
	if err := validateIsSection(p); err != nil {
		return err
	}
	if err := validateNonEmptyKey(p); err != nil {
		return err
	}

	switch p[0] {
	case "media_map":
		return validateSetMediaMap(p, value)
	case "device_map":
		return validateSetDeviceMap(p, value)
	case "logging":
		return validateSetLogging(p, value)
	case "hub":
		return validateSetHub(p, value)
	case "handling":
		return validateSetHandling(p, value)
	case "completion":
		return validateSetCompletion(p, value)
	case "serve":
		return validateSetServe(p, value)
	default:
		return fmt.Errorf("unknown config option: %s", path)
	}
}

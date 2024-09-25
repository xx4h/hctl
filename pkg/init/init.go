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

package init

import (
	"errors"
	"fmt"
	"io/fs"
	u "net/url"
	"os"
	"path/filepath"
	"syscall"

	config "github.com/xx4h/hctl/pkg/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

func InitializeConfig(c *config.Config, configPath string) error {
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("Config already initialized, please use `hctl config` or edit %s directly.\n", configPath)
	} else if errors.Is(err, fs.ErrNotExist) {
		hub := new(config.ConfigHub)
		hub.Type = getHubType()
		hub.Url = getUrl()
		hub.Token = getToken()

		c.Viper.Set("hub", &hub)
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			log.Error().Msgf("Couldn't create config dir: %v\n", err)
		}

		fmt.Printf("\n\n")
		if err := c.Viper.WriteConfigAs(configPath); err != nil {
			log.Error().Msgf("Couldn't write config: %v\n", err)
		}
		return nil
	} else {
		return fmt.Errorf("Unknown Error, this should not happen: %v", err)
	}
	//  if configDirStat, err := os.Stat(configDir); errors(err, os.ErrNotExists)
}

func getHubType() string {
	hubType := "hass"
	// TODO: Enable as soon as we have more supported hub types
	// supported := []string{"hass"}
	// fmt.Printf("Which Hub Type are you using? (Supported: hass) [%s]: ", hubType)
	// _, err := fmt.Scanln(&hubType)
	// if err != nil && err.Error() != "unexpected newline" {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return getHubType()
	// }
	// if !slices.Contains(supported, hubType) {
	// 	fmt.Printf("Unsupported Hub Type: %s\n", hubType)
	// 	return getHubType()
	// }
	return hubType
}

func getUrl() string {
	var url string
	fmt.Print("Enter API URL of your hub (e.g. https://home-assistant.example.com/api): ")
	_, err := fmt.Scanln(&url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return getUrl()
	}
	up, err := u.Parse(url)
	if err != nil || up.Scheme == "" || up.Host == "" {
		fmt.Printf("Not a valid URL: %s\n", url)
		return getUrl()
	}
	return url
}

func getToken() string {
	var token string
	fmt.Print("Enter your hub token: ")
	byteToken, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return getToken()
	}
	token = string(byteToken)
	_, _, err = new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		fmt.Printf("\nNot a valid Token (JWT): %v\n", err)
		return getToken()
	}
	return token
}

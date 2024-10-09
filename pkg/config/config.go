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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Hub        Hub        `mapstructure:"hub" yaml:"hub" json:"hub"`
	Completion Completion `mapstructure:"completion" yaml:"completion" json:"completion"`
	Handling   Handling   `mapstructure:"handling" yaml:"handling" json:"handling"`
	Logging    Logging    `mapstructure:"logging" yaml:"logging" json:"logging"`
	Serve      Serve      `mapstructure:"serve" yaml:"serve" json:"serve"`
	Viper      *viper.Viper
}

type Hub struct {
	Type  string `mapstructure:"type" yaml:"type" json:"type"`
	URL   string `mapstructure:"url" yaml:"url" json:"url"`
	Token string `mapstructure:"token" yaml:"token" json:"token"`
}

type Completion struct {
	ShortNames bool `mapstructure:"short_names" yaml:"short_names" json:"short_names"`
}

type Handling struct {
	Fuzz bool `mapstructure:"fuzz" yaml:"fuzz" json:"fuzz"`
}

type Logging struct {
	LogLevel string `mapstructure:"log_level" yaml:"log_level" json:"log_level"`
}

type Serve struct {
	IP   string `mapstructure:"ip" yaml:"ip" json:"ip"`
	Port int    `mapstructure:"port" yaml:"port" json:"port"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("HCTL")
	v.AutomaticEnv()
	v.SetConfigType("yaml")
	v.SetConfigName("hctl")
	v.AddConfigPath("$HOME/.config/hctl")
	v.AddConfigPath(".")

	// create empty config and set defaults
	cfg := &Config{}
	cfg.Completion.ShortNames = true
	cfg.Handling.Fuzz = true
	cfg.Logging.LogLevel = "error"
	cfg.Serve.IP = ""
	cfg.Serve.Port = 1337
	cfg.Hub.Type = "hass"
	cfg.Hub.URL = ""
	cfg.Hub.Token = ""

	// use defaults for viper as well
	v.SetDefault("completion", &cfg.Completion)
	v.SetDefault("handling", &cfg.Handling)
	v.SetDefault("logging", &cfg.Logging)
	v.SetDefault("serve", &cfg.Serve)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug().Msgf("Config File not found! Please run `hctl init` or manually create %s", v.ConfigFileUsed())
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	logLevel := v.GetString("logging.log_level")
	if logLevel != "" {
		lvl, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Error().Msgf("Could not set log level: %v", err)
		}
		zerolog.SetGlobalLevel(lvl)
	}

	log.Debug().Msgf("Running with the following config: %+v", cfg)

	cfg.Viper = v

	return cfg, nil
}

func (c *Config) GetServeIP() string {
	return c.Serve.IP
}

func (c *Config) GetServePort() int {
	return c.Serve.Port
}

func (c *Config) getElement(p []string) (*reflect.Value, *reflect.Type, error) {
	// create reflects for value and type
	v := reflect.ValueOf(c)
	t := reflect.TypeOf(c)

	return getElementByYamlPath(p, v, t)
}

func getElementByYamlPath(p []string, v reflect.Value, t reflect.Type) (*reflect.Value, *reflect.Type, error) {
	// resolve pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// resolve pointer
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typ := t.Field(i)

		if typ.Tag.Get("yaml") == p[0] {
			if field.Kind() == reflect.Struct {
				return getElementByYamlPath(p[1:], field, typ.Type)
			} else if len(p) == 1 {
				return &field, &typ.Type, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("no such config option: %s", strings.Join(p, "."))
}

func (c *Config) GetOptionsAsPaths() []string {
	t := reflect.ValueOf(c.Viper.AllSettings())
	return toPathSlice(t, "", []string{})
}

func toPathSlice(t reflect.Value, name string, dst []string) []string {
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		return toPathSlice(t.Elem(), name, dst)

	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fname := t.Type().Field(i).Name
			dst = toPathSlice(t.Field(i), strings.TrimLeft(name+"."+fname, "."), dst)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < t.Len(); i++ {
			dst = toPathSlice(t.Index(i), strings.TrimLeft(name+"."+strconv.Itoa(i), "."), dst)
		}

	case reflect.Map:
		for _, key := range t.MapKeys() {
			value := t.MapIndex(key)
			dst = toPathSlice(value, strings.TrimLeft(name+"."+key.String(), "."), dst)
		}

	default:
		return append(dst, name)
	}
	return dst
}

func (c *Config) GetValueByPath(p string) (string, error) {
	log.Info().Msgf("Getting value for `%v`", p)
	s := strings.Split(p, ".")
	v, _, err := c.getElement(s)
	if err != nil {
		return "", err
	}
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool()), nil
	case reflect.Float32:
		return fmt.Sprintf("%f", v.Float()), nil
	case reflect.Float64:
		return fmt.Sprintf("%f", v.Float()), nil
	case reflect.Int:
		return fmt.Sprintf("%d", v.Int()), nil
	default:
		return "", fmt.Errorf("unexpected type: %v", v.Type())
	}
}

func (c *Config) SetValueByPath(p string, val any) error {
	// set config element by path p and value v
	log.Info().Msgf("Setting `%v` to `%v`", p, val)
	s := strings.Split(p, ".")
	v, _, err := c.getElement(s)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Config before change: %+v", c)

	switch v.Kind() {
	case reflect.String:
		v.SetString(val.(string))
	case reflect.Bool:
		e, err := strconv.ParseBool(val.(string))
		if err != nil {
			return err
		}
		v.SetBool(e)
	case reflect.Float32:
		e, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return err
		}
		v.SetFloat(e)
	case reflect.Float64:
		e, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return err
		}
		v.SetFloat(e)
	case reflect.Int:
		e, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(e)
	default:
		return fmt.Errorf("unexpected type: %v", v.Type())
	}

	log.Debug().Msgf("Config after change: %+v", c)

	// convert current config to byte slice
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// create new io.Reader from byte slice config
	reader := bytes.NewReader(b)

	// read in the byte slice config to viper instance
	if err := c.Viper.ReadConfig(reader); err != nil {
		return err
	}

	// finally write updated viper instance to file
	if err := c.Viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

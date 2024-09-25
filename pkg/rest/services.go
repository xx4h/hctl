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

package rest

import (
	"encoding/json"
)

type HassService struct {
	Domain   string                       `json:"domain"`
	Services map[string]HassDomainService `json:"services"`
}

type HassDomainService struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Hass) GetServices() ([]HassService, error) {
	if h.Services != nil {
		return h.Services, nil
	}

	services := []HassService{}
	res, err := h.api("GET", "/services", nil)
	if err != nil {
		return services, err
	}

	err = json.Unmarshal(res, &services)
	if err != nil {
		return services, err
	}

	h.Services = services

	return services, nil
}

func (h *Hass) GetFilteredServices(domains []string, services []string) ([]HassService, error) {
	s, err := h.GetServices()
	if err != nil {
		return nil, err
	}

	s = RemoveDomainsFromServices(s, domains)
	s = RemoveServicesFromServices(s, services)

	return s, nil
}

func (h *Hass) GetFilteredServicesMap(domains []string, services []string) (map[string][]string, error) {
	t := make(map[string][]string)
	s, err := h.GetFilteredServices(domains, services)
	if err != nil {
		return nil, err
	}
	for domain := range s {
		for svc := range s[domain].Services {
			t[s[domain].Domain] = append(t[s[domain].Domain], svc)
		}
	}
	return t, nil
}

// function to check if domain with service exists
func (h *Hass) hasDomainWithService(domain string, service string) (bool, error) {
	services, err := h.GetServices()
	if err != nil {
		return false, err
	}

	for _, svc := range services {
		if svc.Domain == domain {
			for n := range svc.Services {
				if n == service {
					return true, nil
				}
			}
			// domain exists only once
			return false, nil
		}
	}
	return false, nil
}

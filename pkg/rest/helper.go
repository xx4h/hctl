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
	"slices"
	"strings"
)

func FilterDomainsFromStates(s []HassState, domains []string) []HassState {
	if len(domains) == 0 {
		return s
	}

	k := 0
	for _, d := range s {
		e := strings.Split(d.EntityID, ".")
		if slices.Contains(domains, e[0]) {
			s[k] = d
			k++
		}
	}
	s = s[:k]
	return s
}

func FilterDomainsFromServices(s []HassService, domains []string) []HassService {
	if len(domains) == 0 {
		return s
	}

	k := 0
	for _, d := range s {
		if slices.Contains(domains, d.Domain) {
			s[k] = d
			k++
		}
	}
	s = s[:k]
	return s
}

func FilterServicesFromServices(s []HassService, services []string) []HassService {
	if len(services) == 0 {
		return s
	}

	for key := range s {
		for svc := range s[key].Services {
			if !slices.Contains(services, svc) {
				delete(s[key].Services, svc)
			}
		}
	}
	return s
}

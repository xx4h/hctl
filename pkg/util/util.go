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

package util

import (
	"crypto/sha256"
	"fmt"
	"net"
	u "net/url"

	"github.com/rs/zerolog/log"
)

// Return a list of keys from the provided map
// thanks to:
//
//	https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map
func GetStringKeys(m map[string]any) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func GetStringHash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func RemoveIndex(s []any, index int) []any {
	return append(s[:index], s[index+1:]...)
}

func IsURL(url string) bool {
	up, err := u.Parse(url)
	if err != nil || up.Scheme == "" || up.Host == "" {
		return false
	}
	return true
}

func GetLocalIP() string {
	log.Debug().Caller().Msg("Getting local IP")
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		log.Fatal().Msgf("Error getting local IP: %v", err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	log.Debug().Caller().Msgf("Using local IP: %s", localAddress.IP)
	return localAddress.IP.String()
}

func MakeRange(mini, maxi int) []int {
	a := make([]int, maxi-mini+1)
	for i := range a {
		a[i] = mini + i
	}
	return a
}

func MakeRangeString(mini, maxi int) []string {
	a := make([]string, maxi-mini+1)
	for i := range a {
		a[i] = fmt.Sprint(mini + i)
	}
	return a
}

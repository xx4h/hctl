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
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"testing"
)

func Test_isURL(t *testing.T) {
	var tests = map[string]struct {
		url   string
		isURL bool
	}{
		"empty no url": {
			"",
			false,
		},
		"http no url": {
			"http",
			false,
		},
		"http: no url": {
			"http:",
			false,
		},
		"http:/ no url": {
			"http:/",
			false,
		},
		"http:// no url": {
			"http://",
			false,
		},
		"http://a is url": {
			"http://a",
			true,
		},
		"http://a.b is url": {
			"http://a.b",
			true,
		},
		"http://a.b? is url": {
			"http://a.b?",
			true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			isURL := isURL(tt.url)
			if isURL != tt.isURL {
				t.Errorf("got %t, want %t", isURL, tt.isURL)
			}
		})
	}
}

func Test_isJwtToken(t *testing.T) {
	// issuer secret for hmac
	jwtMasterKey := []byte(`supersecretkey`)

	// jwt header and payload
	jwtHeader := b64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	jwtPayload := b64.RawURLEncoding.EncodeToString([]byte(`{"name":"hctl test","iss":"hctl issuer","iat":"1728723175","exp":"1728723175"}`))

	// create jwt and create dot separated `header.payload`
	var jwt []byte
	jwt = append(jwt, []byte(jwtHeader)...)
	jwt = append(jwt, []byte(`.`)...)
	jwt = append(jwt, []byte(jwtPayload)...)

	// create message auth hash for `header.payload`
	sig := hmac.New(sha256.New, jwtMasterKey)
	sig.Write(jwt)
	jwtHmac := b64.RawURLEncoding.EncodeToString([]byte(hex.EncodeToString(sig.Sum(nil))))

	// append hmac to `header.payload` getting `header.payload.hmac`
	jwt = append(jwt, []byte(".")...)
	jwt = append(jwt, []byte(jwtHmac)...)

	var tests = map[string]struct {
		jwt   []byte
		isJwt bool
	}{
		"empty no jwt": {
			[]byte(""),
			false,
		},
		"hctl in b64 no jwt": {
			// just hctl in base64
			[]byte(b64.RawURLEncoding.EncodeToString([]byte(`hctl`))),
			false,
		},
		"jwt secret is jwt": {
			jwt,
			true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			isJwt := isJwtToken(tt.jwt)
			if isJwt != tt.isJwt {
				t.Errorf("got %t, want %t", isJwt, tt.isJwt)
			}
		})
	}

}

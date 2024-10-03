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

package serve

import (
	"net"
	"net/http"
	"sync"
)

/*
Big thanks to https://stackoverflow.com/a/62766994/1922402
*/

type ConnectionWatcher struct {
	// mu protects remaining fields
	mu sync.Mutex

	// open connections are keys in the map
	m map[net.Conn]struct{}
}

// OnStateChange records open connections in response to connection
// state changes. Set net/http Server.ConnState to this method
// as value.
func (cw *ConnectionWatcher) OnStateChange(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		cw.mu.Lock()
		if cw.m == nil {
			cw.m = make(map[net.Conn]struct{})
		}
		cw.m[conn] = struct{}{}
		cw.mu.Unlock()
	case http.StateHijacked, http.StateClosed:
		cw.mu.Lock()
		delete(cw.m, conn)
		cw.mu.Unlock()
	}
}

// Connections returns the open connections at the time
// the call.
func (cw *ConnectionWatcher) Connections() []net.Conn {
	var result []net.Conn
	cw.mu.Lock()
	for conn := range cw.m {
		result = append(result, conn)
	}
	cw.mu.Unlock()
	return result
}

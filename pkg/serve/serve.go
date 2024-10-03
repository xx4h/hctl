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
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/xx4h/hctl/pkg/util"
)

type Media struct {
	ip        string
	port      int
	path      string
	url       string
	serveName string
	srv       *http.Server
	wg        *sync.WaitGroup
	cw        *ConnectionWatcher
}

func NewMedia(ip string, port int, path string) *Media {
	serveName := util.GetStringHash(path)
	if ip == "" {
		ip = util.GetLocalIP()
	}

	return &Media{
		cw:        &ConnectionWatcher{},
		wg:        &sync.WaitGroup{},
		ip:        ip,
		port:      port,
		path:      path,
		serveName: serveName,
		url:       fmt.Sprintf("http://%s:%d/%s", ip, port, serveName),
	}
}

func (m *Media) GetURL() string {
	return m.url
}

func (m *Media) GetMediaName() string {
	return filepath.Base(m.path)
}

func (m *Media) serveFile() {
	defer m.wg.Done()

	log.Debug().Msgf("Starting local HTTP server: %s:%d", m.ip, m.port)

	mux := http.NewServeMux()
	m.srv = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", m.ip, m.port),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           mux,
		ConnState:         m.cw.OnStateChange,
	}

	mux.HandleFunc(fmt.Sprintf("/%s", m.serveName), func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, m.path)
		log.Debug().Msgf("File requested: /%s -> %s", m.serveName, m.path)
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "ready")
	})

	go func() {
		if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("HTTP server error: %v", err)
		}
	}()
}

func (m *Media) FileHandler() {
	m.wg.Add(1)
	m.serveFile()
}

func (m *Media) WaitForHTTPReady() error {
	readyEndpoint := fmt.Sprintf("http://%s:%d/ready", m.ip, m.port)
	tries := 10
	var err error
	log.Debug().Msg("Waiting for HTTP server to be ready")
	for i := 0; i < tries; i++ {
		log.Debug().Msgf("Checking readiness endpoint: %s", readyEndpoint)
		resp, err := http.Get(readyEndpoint) // nolint:gosec
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		defer resp.Body.Close()
		log.Debug().Msg("HTTP server is ready")
		return nil
	}
	return err
}

func (m *Media) WaitAndClose() error {
	var count int
	for {
		conns := m.cw.Connections()
		time.Sleep(200 * time.Millisecond)
		if len(conns) == 0 {
			break
		}
		if count%2 == 0 {
			log.Debug().Msgf("Waiting for connection to close: %#v", conns[0].RemoteAddr().String())
		}
		count++
	}
	if err := m.srv.Close(); err != nil {
		return err
	}
	m.wg.Wait()
	return nil
}

package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mongche/gupload/config"
	"golang.org/x/net/http2"
)

const (
	MaxConnectionTimeout = 5
)

// The HTTP/2.0 server is used for file uploading.
type Http2Server struct {
	server *http.Server
}

// Establish a new HTTP/2.0 server.
func NewHttp2Server(cfg config.Http2ServerConfig) (s Http2Server, err error) {
	if cfg.Port == 0 {
		err = errors.New("port must be required")
		return
	}

	if cfg.Cert == "" {
		err = errors.New("cert must be required")
		return
	}

	if cfg.Key == "" {
		err = errors.New("key must be required")
		return
	}

	s.server = &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
	http2.ConfigureServer(s.server, nil)
	http.HandleFunc("/upload", s.Upload)

	if err = s.server.ListenAndServeTLS(cfg.Cert, cfg.Key); err != nil {
		err = fmt.Errorf("failed to listen HTTP/2.0 server, %v", err)
		return
	}

	defer s.server.Close()

	_, cancel := context.WithTimeout(context.Background(), MaxConnectionTimeout*time.Second)
	defer cancel()

	return
}

// Upload a file to the server.
func (s *Http2Server) Upload(w http.ResponseWriter, r *http.Request) {
	var (
		buf = new(bytes.Buffer)
	)

	content, err := io.Copy(buf, r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "%+v", err)
		return
	}

	log.Printf("bytes received from the client, %v", content)
	return
}

// Close the server connection.
func (s *Http2Server) Close() {
	return
}

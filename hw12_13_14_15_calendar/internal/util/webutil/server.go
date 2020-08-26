package webutil

import (
	"context"
	"net/http"
	"time"
)

type WebServer struct {
	server *http.Server
}

func NewWebServer(handler http.Handler, listenAddress string) *WebServer {
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &WebServer{server: server}
}

func (s WebServer) Start() error {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s WebServer) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

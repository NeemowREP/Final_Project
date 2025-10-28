package server

import (
	"Final_Project/internal/settings"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	cfg        settings.Config
}

func New(cfg settings.Config) *Server {
	mux := NewRouter(cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
	}

	return &Server{httpServer: srv, cfg: cfg}

}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

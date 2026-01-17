package server

import (
	"log"
	"net/http"
	"time"

	"Final_Project/internal/settings"
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
	log.Printf("Server is running on port %s", s.cfg.Port)
	return s.httpServer.ListenAndServe()
}

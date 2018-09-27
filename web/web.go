package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/vitalyo61/genkeys/config"
	"github.com/vitalyo61/genkeys/db"
	"github.com/vitalyo61/genkeys/generator"
)

type Server struct {
	srv     *http.Server
	address string
	timeout time.Duration
}

func Make(cfg *config.Server, b *db.DB, ch chan *generator.Data) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              cfg.Address,
			ReadTimeout:       time.Second * time.Duration(cfg.Timeout),
			ReadHeaderTimeout: time.Second * time.Duration(cfg.Timeout),
			WriteTimeout:      time.Second * time.Duration(cfg.Timeout),
			IdleTimeout:       time.Second * time.Duration(cfg.Timeout),
			Handler:           makeRouter(b, ch),
		},
	}
}

func (s *Server) Start() {
	log.Println("Server started...")
	log.Fatal(s.srv.ListenAndServe())
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.timeout))
	defer cancel()
	log.Println("Shutting down server")
	return s.srv.Shutdown(ctx)
}

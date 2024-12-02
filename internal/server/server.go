package server

import (
	"net/http"

	"github.com/elct9620/poc-raft-api/internal/store"
)

type Server struct {
	httpServer *http.Server
	store      *store.Store
}

func NewServer(
	store *store.Store,
) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: ":8080",
		},
		store: store,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}

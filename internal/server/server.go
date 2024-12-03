package server

import (
	"net/http"

	"github.com/hashicorp/raft"
)

type Server struct {
	httpServer *http.Server
	raft       *raft.Raft
}

func NewServer(
	raft *raft.Raft,
) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: ":8080",
		},
		raft: raft,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}

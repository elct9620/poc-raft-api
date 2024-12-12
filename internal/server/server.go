package server

import (
	"net/http"

	"github.com/hashicorp/raft"
)

type KvRepository interface {
	Set(key string, value string) error
}

type Server struct {
	httpServer *http.Server
	raft       *raft.Raft
	kvStore    KvRepository
}

func NewServer(
	raft *raft.Raft,
	kvStore KvRepository,
) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		raft:    raft,
		kvStore: kvStore,
	}

	mux.HandleFunc("POST /join", srv.PostJoin)
	mux.HandleFunc("PUT /value", srv.PutValue)

	return srv
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}

package server

import (
	"net/http"

	"github.com/hashicorp/raft"
)

type KvStore interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
}

type Server struct {
	httpServer *http.Server
	raft       *raft.Raft
	store      KvStore
}

func NewServer(
	raft *raft.Raft,
	store KvStore,
) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		raft:  raft,
		store: store,
	}

	mux.HandleFunc("POST /join", srv.PostJoin)

	return srv
}

func (s *Server) Start() error {

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}

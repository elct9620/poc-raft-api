package server

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/raft"
)

type JoinRequest struct {
	NodeID      string `json:"node_id"`
	NodeAddress string `json:"node_address"`
}

func (s *Server) PostJoin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var joinRequest JoinRequest
	if err := decoder.Decode(&joinRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if s.raft.State() != raft.Leader {
		http.Error(w, "not leader", http.StatusServiceUnavailable)
		return
	}

	if err := s.raft.AddVoter(raft.ServerID(joinRequest.NodeID), raft.ServerAddress(joinRequest.NodeAddress), 0, 0).Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

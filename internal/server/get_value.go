package server

import (
	"encoding/json"
	"net/http"
)

type GetValueResponse struct {
	Value string `json:"value"`
}

func (srv *Server) GetValue(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	value := srv.kvStore.Get(key)
	if value == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(GetValueResponse{Value: value}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

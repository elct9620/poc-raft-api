package server

import (
	"encoding/json"
	"net/http"
)

type PutValueRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutValueResponse struct {
	Success bool `json:"success"`
}

func (srv *Server) PutValue(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var putValueRequest PutValueRequest
	if err := decoder.Decode(&putValueRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := srv.kvStore.Set(putValueRequest.Key, putValueRequest.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(PutValueResponse{Success: true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

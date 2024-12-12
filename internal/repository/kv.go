package repository

import (
	"encoding/json"
	"time"

	"github.com/elct9620/poc-raft-api/internal/app"
	"github.com/hashicorp/raft"
)

type KeyValueRepository struct {
	raft  *raft.Raft
	state *app.State
}

func NewKeyValueRepository(raft *raft.Raft, state *app.State) *KeyValueRepository {
	return &KeyValueRepository{
		raft:  raft,
		state: state,
	}
}

func (r *KeyValueRepository) Set(key string, value string) error {
	cmd := app.StateCommand{
		Key:   key,
		Value: value,
	}

	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return r.raft.Apply(payload, 500*time.Millisecond).Error()
}

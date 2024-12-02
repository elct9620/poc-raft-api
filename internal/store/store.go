package store

import (
	"github.com/hashicorp/raft"
)

type Store struct {
	raft *raft.Raft
}

func NewStore() *Store {
	return &Store{}
}

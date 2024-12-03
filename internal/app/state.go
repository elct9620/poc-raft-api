package app

import (
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type noopSnapshot struct{}

func (n noopSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (n noopSnapshot) Release() {}

type state struct {
	db *sync.Map
}

func newState() *state {
	return &state{
		db: &sync.Map{},
	}
}

func (f *state) Apply(log *raft.Log) any {
	return nil
}

func (f *state) Restore(io.ReadCloser) error {
	return nil
}

func (f *state) Snapshot() (raft.FSMSnapshot, error) {
	return noopSnapshot{}, nil
}

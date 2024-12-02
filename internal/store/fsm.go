package store

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

type fsm struct {
	db *sync.Map
}

func newFsm() *fsm {
	return &fsm{
		db: &sync.Map{},
	}
}

func (f *fsm) Apply(log *raft.Log) any {
	return nil
}

func (f *fsm) Restore(io.ReadCloser) error {
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return noopSnapshot{}, nil
}

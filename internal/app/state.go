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

type State struct {
	db *sync.Map
}

func NewState() *State {
	return &State{
		db: &sync.Map{},
	}
}

func (f *State) Apply(log *raft.Log) any {
	return nil
}

func (f *State) Restore(io.ReadCloser) error {
	return nil
}

func (f *State) Snapshot() (raft.FSMSnapshot, error) {
	return noopSnapshot{}, nil
}

func (f *State) Get(key string) (any, bool) {
	v, ok := f.db.Load(key)
	if !ok {
		return nil, false
	}

	return v, true
}

func (f *State) Set(key string, value any) {
	f.db.Store(key, value)
}

func (f *State) Delete(key string) {
	f.db.Delete(key)
}

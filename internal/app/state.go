package app

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type noopSnapshot struct{}

func (n noopSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (n noopSnapshot) Release() {}

type StateCommand struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type State struct {
	db *sync.Map
}

func NewState() *State {
	return &State{
		db: &sync.Map{},
	}
}

func (f *State) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		var cmd StateCommand
		if err := json.Unmarshal(log.Data, &cmd); err != nil {
			return err
		}

		f.Set(cmd.Key, cmd.Value)
	default:
		return fmt.Errorf("unhandled log type %v", log.Type)
	}

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

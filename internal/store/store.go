package store

import (
	"net"
	"os"
	"path"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
)

const SnapshotRetainCount = 2
const RaftMaxPool = 10
const RaftTimeout = 10 * time.Second

type Store struct {
	raft *raft.Raft
}

func NewStore(nodeId string, rootDir string, raftAddress string) (*Store, error) {
	boltStore, err := raftboltdb.NewBoltStore(path.Join(rootDir, "boltdb"))
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(path.Join(rootDir, "snapshots"), SnapshotRetainCount, os.Stderr)
	if err != nil {
		return nil, err
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(raftAddress, tcpAddress, RaftMaxPool, RaftTimeout, os.Stderr)
	if err != nil {
		return nil, err
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeId)

	fsm := newFsm()

	r, err := raft.NewRaft(config, fsm, boltStore, boltStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	})

	return &Store{
		raft: r,
	}, nil
}

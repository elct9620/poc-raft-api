package app

import (
	"net"
	"os"
	"path"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const RaftSnapshotRetainCount = 2
const RaftMaxPool = 10
const RaftTimeout = 10 * time.Second

func NewRaft(nodeId string, rootDir string, raftAddress string) (*raft.Raft, error) {
	boltStore, err := raftboltdb.NewBoltStore(path.Join(rootDir, "boltdb"))
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(path.Join(rootDir, "snapshots"), RaftSnapshotRetainCount, os.Stderr)
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
	state := newState()

	r, err := raft.NewRaft(config, state, boltStore, boltStore, snapshots, transport)
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

	return r, nil
}

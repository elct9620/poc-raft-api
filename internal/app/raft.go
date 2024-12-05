package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/elct9620/poc-raft-api/internal/config"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const RaftSnapshotRetainCount = 2
const RaftMaxPool = 10
const RaftTimeout = 10 * time.Second

type JoinRequest struct {
	NodeID      string `json:"node_id"`
	NodeAddress string `json:"node_address"`
}

func NewRaft(cfg *config.Config, state raft.FSM) (*raft.Raft, error) {
	boltStore, err := raftboltdb.NewBoltStore(path.Join(cfg.DataDir(), "boltdb"))
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(path.Join(cfg.DataDir(), "snapshots"), RaftSnapshotRetainCount, os.Stderr)
	if err != nil {
		return nil, err
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", cfg.RaftAddress())
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(cfg.RaftAddress(), tcpAddress, RaftMaxPool, RaftTimeout, os.Stderr)
	if err != nil {
		return nil, err
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(cfg.Hostname())

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

	if cfg.IsLeader() {
		return r, nil
	}

	payload, err := json.Marshal(JoinRequest{
		NodeID:      cfg.Hostname(),
		NodeAddress: cfg.RaftAddress(),
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.RaftLeaderApi()+"/join", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to join cluster: %s", res.Status)
	}

	return r, nil
}

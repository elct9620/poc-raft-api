package config

import "os"

type Config struct {
	hostname      string
	dataDir       string
	raftAddress   string
	raftLeaderApi string
}

func New() *Config {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	raftAddress := os.Getenv("RAFT_ADDRESS")
	if raftAddress == "" {
		raftAddress = "localhost:2773"
	}

	raftLeaderApi := os.Getenv("RAFT_LEADER_API")

	return &Config{
		hostname:      hostname,
		dataDir:       dataDir,
		raftAddress:   raftAddress,
		raftLeaderApi: raftLeaderApi,
	}
}

func (c *Config) Hostname() string {
	return c.hostname
}

func (c *Config) DataDir() string {
	return c.dataDir
}

func (c *Config) RaftAddress() string {
	return c.raftAddress
}

func (c *Config) IsLeader() bool {
	return c.raftLeaderApi == ""
}

func (c *Config) RaftLeaderApi() string {
	return c.raftLeaderApi
}

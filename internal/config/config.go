package config

import "os"

type Config struct {
	hostname    string
	dataDir     string
	raftAddress string
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

	return &Config{
		hostname:    hostname,
		dataDir:     dataDir,
		raftAddress: raftAddress,
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

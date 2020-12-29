package tcpserver

import (
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"net"

	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
	"testing"
)

// TestTCPServer ...
type TestTCPServer struct {
	Server  net.Listener
	Config  *Config
	Chans   *NewChannels
	Store   store.Store
	Pool    map[*net.Conn]*models.Connection
	Working bool
}

// TestServer ...
func TestServer(t *testing.T, conf *Config) (Server, error) {

	return &TestTCPServer{
		// Pool:   pool,
		Config: conf,
		// Chans:  channels,
		// Store:  store,
	}, nil

}

// Run ...
func (t *TestTCPServer) Run() error {

	var err error

	t.Working = true

	for {
		t.Server, err = net.Listen(t.Config.Protocol, t.Config.Port)
		if err != nil {
			return err
		}
	}
	// return nil
}

// Close ...
func (t *TestTCPServer) Close() error {

	t.Working = false

	return t.Server.Close()
	//return nil
}

// Listener ...
func (t *TestTCPServer) Listener() error {

	return nil
}

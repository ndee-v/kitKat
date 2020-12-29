package tcpserver

import (
	testhistory "kitKat/internal/app/store/testHistory"
	"net"
	"netcat/internal/app/models"

	//testhistory "netcat/internal/app/store/testHistory"
	"testing"
	"time"
)

// TestConfig ...
func TestConfig(t *testing.T) *Config {

	return &Config{
		Protocol: "tcp",
		TimeReg:  time.Second * 30,
		TimeIdle: time.Second * 180,
		TimeOff:  time.Second * 2,
		Rooms:    4,
		Conns:    10,
		Port:     ":9000",
		Gui:      false,
		Testing:  true,
	}

}

// TestServerOld ....
func TestServerOld(t *testing.T, conf *Config) (Server, error) {

	channels := &NewChannels{
		Conns:   make(chan *net.Conn),
		Dconns:  make(chan *models.Message),
		Msgs:    make(chan *models.Message),
		GuiMsgs: make(chan *models.Message),
		Errors:  make(chan *models.ErrMes),
	}

	store := testhistory.New(t)

	if err := store.History().Create("newStory", conf.Rooms); err != nil {
		return nil, err
	}

	pool, err := store.Pool().Create(conf.Conns, conf.Rooms)
	if err != nil {
		return nil, err
	}

	return &NewTCPServer{
		Pool:   pool,
		Config: conf,
		Chans:  channels,
		Store:  store,
	}, nil

}

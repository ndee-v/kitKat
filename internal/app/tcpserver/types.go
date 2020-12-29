package tcpserver

import (
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"net"
	"time"
)

// Server defines ...
type Server interface {
	Run() error
	Close() error
	Listener() error
}

//changes after this message ...

// TCPServer holds the structure of our TCP
// implementation.
type TCPServer struct {
	Protocol string
	Port     string
	Server   net.Listener
	Config   *Config
	History  *models.History
	Pool     *map[net.Conn]*models.Connection //pool of all connections
	Chans    *Channels
}

// NewTCPServer using for refactoring code
// new functions and new methods
type NewTCPServer struct {
	Server  net.Listener
	Config  *Config
	Chans   *NewChannels
	Store   store.Store
	Pool    map[*net.Conn]*models.Connection
	Working bool
}

// Config struct using for prog configuration
/*
type Config struct {
	Protocol string // "tcp"
	//Listener net.Listener // listener for TCP chat
	Port string // port for listener
	// History  *models.History                 // struct with full addr to history files
	Rooms    int           // count of chat rooms
	Conns    int           // count of max connections
	TimeReg  time.Duration // time for registration
	TimeIdle time.Duration // time for idle timeout
	// Pool     *map[net.Conn]*model.Connection //pool of all connections
	Gui bool // GUI terminal on/off
}
*/
// Config struct using for prog configuration
type Config struct {
	Protocol string        // "tcp"
	Port     string        // port for listener
	Rooms    int           // count of chat rooms
	Conns    int           // count of max connections
	TimeReg  time.Duration // time for registration
	TimeIdle time.Duration // time for idle timeout
	TimeOff  time.Duration // time for server closing
	Gui      bool          // GUI terminal on/off
	Testing  bool          //use for testconfig
}

// Channels struct using for all channels in program
type Channels struct {
	Conns   chan net.Conn        //channel for connections
	Dconns  chan *models.Message //channel for disconnections
	Msgs    chan *models.Message //channel for messages
	GuiMsgs chan *models.Message //channel for GUI terminal (optional)
	Errors  chan error           // channel for errors from goroutines
}

// NewChannels struct using for all channels in program
type NewChannels struct {
	Conns   chan *net.Conn       //channel for connections
	Dconns  chan *models.Message //channel for disconnections
	Msgs    chan *models.Message //channel for messages
	GuiMsgs chan *models.Message //channel for GUI terminal (optional)
	Errors  chan *models.ErrMes  // channel for errors from goroutines
}

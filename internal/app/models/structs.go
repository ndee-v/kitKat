package models

import (
	"net"
	"os"
	"time"
)

//Connection struct...
type Connection struct {
	Test   bool
	Name   string    // name/login of user
	Action time.Time // time of last action
	Room   int       // chat room
	Conn   *net.Conn
	//History *History  // last messages from histories
	LastMsg []string //last received message
}

//Message ...
type Message struct {
	Author      string    //author of message
	Conn        *net.Conn //connection
	In          bool      //using if user entered to chat
	Out         bool      //using if user left from chat
	Kick        bool      //using if program kicks user from chat
	Move        bool      //using if user changing chat room
	ChangeName  string    //using if user changing name
	Text        string    //text of message
	Time        time.Time //time of posted
	Room        int       //chat room
	RoomFrom    int       //if user changed rooms
	CheckOnline bool      //using for get info about online users
	Help        bool      //using for get help options
	All         bool      //using for messages to all users
}

// Info ...
type Info struct {
	Conns int
}

// History ...
type History struct {
	Dir  string
	Room []*os.File
}

// ErrMes ...
type ErrMes struct {
	From  string
	Error error
}

package store

import (
	"kitKat/internal/app/models"
	"net"
	//"netcat/internal/app/models"
	"time"
)

// HistoryRepo ...
type HistoryRepo interface {
	Create(string, int) error             //initialize files for history needed string of directory and number of rooms
	AddInto(*models.Message) error        //save meggase to history file
	PrintToConn(*models.Connection) error // show messages to connections
}

// PoolRepo ...
type PoolRepo interface {
	//MainPool() *models.Info                                        //returns info struct
	Create(int, int) (map[*net.Conn]*models.Connection, error)     //initialize max number of connections
	Handler(*net.Conn, *time.Duration) (*models.Connection, error) //registration for each new connection
	Add(*models.Connection) error                                  //add new user to pool of connections
	Get(*net.Conn) (*models.Connection, error)                     //returns connection from pool
	DeleteByConn(*net.Conn) error                                  //deleting from pool
	RemoveUser(*models.Message) error                              //remove user from pool
	Change(*models.Message) error                                  // using if user send '-' at start of message
	RenameUser(*models.Message) error                              //changing name of connection
	ChangeRoom(*models.Message) error                              //changing room of chat
	WhoIsOnline(*models.Message) error                             //show online users to connection
	IdleTimer(*bool, *time.Duration, *chan *models.Message) error  //timer for auto kick by idle timeouts
	RegTimer(*net.Conn, *time.Duration) error                      //timer for auto kick if registration so long
	SendMsg(*models.Message) error                                 //sending message from user to users
}

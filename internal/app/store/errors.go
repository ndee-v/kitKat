package store

import "errors"

var (
	//ErrDirString ...
	ErrDirString = errors.New("invalid directory string")

	//ErrNumRooms is error number of rooms
	ErrNumRooms = errors.New("invalid numbers of rooms")

	//ErrNumConns is error number of rooms
	ErrNumConns = errors.New("invalid numbers of max connections")

	//ErrNotADir ...
	ErrNotADir = errors.New("is not a directory")

	//ErrHistoryNotExist ...
	ErrHistoryNotExist = errors.New("history is not exist")

	//ErrRoomNotFound ...
	ErrRoomNotFound = errors.New("room not found")

	//ErrConnNotInit ...
	ErrConnNotInit = errors.New("connection does not init")

	//ErrMaxConn ...
	ErrMaxConn = errors.New("maximum number of connections has been reached")

	//ErrDelNothing ...
	ErrDelNothing = errors.New("nothing to delete")

	//ErrConnNotAdded ...
	ErrConnNotAdded = errors.New("connection not in pool")

	//ErrInvalidName ...
	ErrInvalidName = errors.New("invalid name")

	//ErrNameUsed ...
	ErrNameUsed = errors.New("this name has been using")

	//ErrInvalidInput ...
	ErrInvalidInput = errors.New("invalid argument(s)")
)

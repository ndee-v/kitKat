package pool

import (
	"net"
	"netcat/internal/app/models"
	"netcat/internal/app/store"
)

// Repo ...
type Repo struct {
	Store    store.Store
	Capacity int
	Rooms    int
	Pool     map[*net.Conn]*models.Connection
}

package pool

import (
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"net"
)

// Repo ...
type Repo struct {
	Store    store.Store
	Capacity int
	Rooms    int
	Pool     map[*net.Conn]*models.Connection
}

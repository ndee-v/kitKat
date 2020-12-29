package pool

import (
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"net"
	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
)

// // MainPool ...
// func (r *Repo) MainPool() *models.Info {

// 	return &models.Info{
// 		Conns: r.Capacity,
// 	}
// }

//Create ...
func (r *Repo) Create(cap, rooms int) (map[*net.Conn]*models.Connection, error) {

	if cap <= 0 {
		return nil, store.ErrNumConns
	}
	if rooms <= 0 {
		return nil, store.ErrNumRooms
	}

	r.Rooms = rooms
	r.Capacity = cap

	return r.Pool, nil
}

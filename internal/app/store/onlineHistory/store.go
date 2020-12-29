package onlinehistory

import (
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"kitKat/internal/app/store/pool"
	"net"

	// "netcat/internal/app/models"
	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
	// "netcat/internal/app/store/pool"
	"os"
)

// Store ...
type Store struct {
	store       *store.Store
	historyRepo *HistoryRepo
	poolRepo    *pool.Repo
}

// New ...
func New() *Store {
	return &Store{}
}

//History ...
func (s *Store) History() store.HistoryRepo {
	if s.historyRepo != nil {
		return s.historyRepo
	}
	s.historyRepo = &HistoryRepo{

		History: make(map[int]*os.File),
	}
	return s.historyRepo
}

//Pool ...
func (s *Store) Pool() store.PoolRepo {
	if s.poolRepo != nil {
		return s.poolRepo
	}

	s.poolRepo = &pool.Repo{
		Store: s,
		Pool:  make(map[*net.Conn]*models.Connection),
	}
	return s.poolRepo
}

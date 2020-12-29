package testhistory

import (
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"kitKat/internal/app/store/pool"
	"net"

	//"netcat/internal/app/models"
	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
	// "netcat/internal/app/store/pool"
	"testing"
)

// Store ...
type Store struct {
	store       *store.Store
	historyRepo *HistoryRepo
	poolRepo    *pool.Repo
}

// New ...
func New(t *testing.T) *Store {

	return &Store{}

}

//History ...
func (s *Store) History() store.HistoryRepo {
	if s.historyRepo != nil {
		return s.historyRepo
	}

	s.historyRepo = &HistoryRepo{

		History: make(map[int][]string),
	}
	return s.historyRepo
}

// Pool ...
func (s *Store) Pool() store.PoolRepo {
	if s.poolRepo != nil {
		return s.poolRepo
	}
	rep := &pool.Repo{
		Store: s,
		//History: s.historyRepo,
		Pool: make(map[*net.Conn]*models.Connection),
	}
	s.poolRepo = rep
	return rep
}

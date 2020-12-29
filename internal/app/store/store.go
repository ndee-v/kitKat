package store

// Store ...
type Store interface {
	History() HistoryRepo
	Pool() PoolRepo
}

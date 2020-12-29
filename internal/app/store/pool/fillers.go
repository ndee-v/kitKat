package pool

import (
	"fmt"
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"log"
	"net"
	//	"netcat/internal/app/models"
	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
)

// Add ...
func (r *Repo) Add(m *models.Connection) error {

	if r.Capacity <= 0 {
		return store.ErrNumConns
	}

	if m.Conn == nil {
		return store.ErrConnNotInit
	}

	if r.Capacity <= len(r.Pool) {

		str1 := fmt.Sprintf("\n[MAXIMUM CONNECTIONS HAS BEEN REACHED]")

		if err := m.Write(str1); err != nil {
			return err
		}
		return store.ErrMaxConn
	}

	str2 := fmt.Sprintf("\n[USER %v JOIN TO CHAT]", m.Name)

	r.Pool[m.Conn] = m

	for _, val := range r.Pool {

		if val.Room == m.Room && val.Name != m.Name {

			if err := val.Write(str2); err != nil {
				return err
			}
			if err := val.Prefix(); err != nil {
				return err
			}

		}

	}

	return nil
}

// Get ...
func (r *Repo) Get(c *net.Conn) (*models.Connection, error) {

	if c == nil {
		return nil, store.ErrConnNotInit
	}

	val, ok := r.Pool[c]
	if !ok {
		return nil, store.ErrConnNotAdded
	}

	return val, nil
}

// DeleteByConn ...
func (r *Repo) DeleteByConn(c *net.Conn) error {
	if c == nil {
		return store.ErrDelNothing
	}

	delete(r.Pool, c)

	if err := (*c).Close(); err != nil {
		return err
	}

	return nil
}

// RemoveUser ...
func (r *Repo) RemoveUser(m *models.Message) error {

	if m == nil {
		return store.ErrInvalidInput
	}
	if m.Conn == nil {
		return store.ErrConnNotInit
	}

	if !m.Kick && !m.Out {
		return store.ErrInvalidInput
	}

	reason := fmt.Sprintf("\n[%v HAS LEFT OUR CHAT]", m.Author)
	if m.Kick {
		reason = fmt.Sprintf("\n[%v HAS KICKED BY IDLE]", m.Author)
		if _, err := (*m.Conn).Write([]byte("\n[YOU WERE REMOVED FROM CHAT BY IDLE TIMEOUT]\n")); err != nil {
			return err
		}

	}

	if err := r.DeleteByConn(m.Conn); err != nil {
		return err
	}

	logReason := "by himself"
	if m.Kick {
		logReason = "by server"
	}
	log.Printf("%v disconnected %v", (*m.Conn).RemoteAddr().String(), logReason)

	for _, val := range r.Pool {
		if val.Room == m.Room && val.Name != m.Author {
			if err := val.Write(reason); err != nil {
				return err
			}

			if err := val.Prefix(); err != nil {
				return err
			}
		}
	}

	return nil
}

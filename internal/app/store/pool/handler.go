package pool

import (
	"bufio"
	"kitKat/internal/app/models"
	"kitKat/internal/app/utils"
	"log"
	"net"

	// "netcat/internal/app/models"
	// "netcat/internal/app/utils"
	"time"
)

// Handler ...
func (r *Repo) Handler(c *net.Conn, t *time.Duration) (*models.Connection, error) {

	// if len(r.Pool) >= r.Capacity {
	// 	if _, err := (*c).Write([]byte(store.ErrMaxConn.Error())); err != nil {
	// 		return err
	// 	}
	// 	return (*c).Close()
	// }
	var (
		err  error
		text string
	)

	if err = utils.Hello(c); err != nil {
		return nil, err
	}

	go func() {
		if err = r.RegTimer(c, t); err != nil {
			log.Println(err)
		}
	}()

	for {

		busy := false

		rd := bufio.NewReader(*c)

		text, err = rd.ReadString('\n')
		if err != nil {

			break
			//return err
		}
		text = text[:len(text)-1]

		if !utils.ValidName(text) {
			if _, err = (*c).Write([]byte("[INCORRECT NAME]:")); err != nil {
				return nil, err
			}
			continue
		}
		for _, val := range r.Pool {
			if val.Name == text {
				if _, err = (*c).Write([]byte("[NAME ALREADY USING]:")); err != nil {
					return nil, err
				}
				busy = true
				break
			}
		}

		if !busy {

			m := &models.Connection{
				Name:    text,
				Conn:    c,
				Action:  time.Now(),
				Room:    1,
				LastMsg: make([]string, r.Rooms+1),
			}

			if err := r.Add(m); err != nil {
				if err2 := (*c).Close(); err2 != nil {
					return nil, err2
				}
				return nil, err
			}

			if err := m.Write("[FOR HELP USAGE]: --help"); err != nil {
				return nil, err
			}

			if err := r.Store.History().PrintToConn(m); err != nil {
				return nil, err
			}

			// if err := (*r.Store).History().PrintToConn(m); err != nil {
			// 	return err
			// }

			// if err := m.Prefix(); err != nil {
			// 	return
			// }
			//return m.Prefix()
			return m, nil

		}

	}

	return nil, err
}

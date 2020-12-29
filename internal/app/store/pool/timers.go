package pool

import (
	"errors"
	"fmt"
	"kitKat/internal/app/models"
	"log"
	"net"
	"time"
)

// IdleTimer ...
func (r *Repo) IdleTimer(working *bool, t *time.Duration, dconns *chan *models.Message) error {

	var err error

	for {

		time.Sleep(*t)

		if !(*working) {
			fmt.Printf("cleaner stopping ..., at time:%v\n", time.Now())
			break
		}

		date := time.Now().Add((*t) * -1)

		for _, val := range r.Pool {

			if val.Action.Before(date) {

				mes := &models.Message{
					Author: val.Name,
					Conn:   val.Conn,
					Kick:   true,
					Room:   val.Room,
				}

				select {
				case _, ok := <-*dconns:
					if !ok {
						return errors.New("channel closed")
					}
					// case <-time.After(100 * time.Millisecond):
					// 	*dconns <- mes
				default:
					*dconns <- mes
				}

			}

		}

	}
	return err
}

// RegTimer ...
func (r *Repo) RegTimer(c *net.Conn, t *time.Duration) error {

	time.Sleep(*t)

	_, ok := (r.Pool)[c]
	if !ok {
		if _, err := (*c).Write([]byte("\n[YOU HAS BEEN KICKED BY IDLE TIMEOUT]\n")); err != nil {
			return err
		}

		log.Printf("Connection closing by idle timeout: %v", (*c).RemoteAddr().String())
		return (*c).Close()
	}
	return nil
}

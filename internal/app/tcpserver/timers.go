package tcpserver

import (
	"kitKat/internal/app/models"
	"log"
	"net"
	"time"
)

// IdleTimer func checking idle timeout
// then auto kick, if timeout is
// over than three minutes
func idleTimer(t *TCPServer) {

	for {
		time.Sleep(time.Second * 10)

		for k, v := range *t.Pool {
			if v.Action.Before(time.Now().Add(t.Config.TimeIdle * -1)) {
				if _, err := k.Write([]byte("[YOU HAS BEEN KICKED BY IDLE TIMEOUT]\n")); err != nil {
					t.Chans.Errors <- err
					return
				}
				log.Printf("Connection closing by idle timeout: %v", k.RemoteAddr().String())

				mes := &models.Message{
					Author: v.Name,
					Conn:   v.Conn,
					Kick:   true,
				}

				t.Chans.Dconns <- mes
			}
		}
	}
}

// regTimer func using as auto kick connection
// if user connected, but doing nothing over than 30 seconds
func regTimer(c *net.Conn, t *TCPServer) error {

	time.Sleep(t.Config.TimeReg)

	_, ok := (*t.Pool)[*c]
	if !ok {
		if _, err := (*c).Write([]byte("\n[YOU HAS BEEN KICKED BY IDLE TIMEOUT]\n")); err != nil {
			return nil
		}
		log.Printf("Connection closing by idle timeout: %v", (*c).RemoteAddr().String())

		return (*c).Close()
	}
	return nil
}

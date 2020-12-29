package tcpserver

import (
	"bufio"
	"kitKat/internal/app/models"
	"kitKat/internal/app/utils"
	"net"
	"strings"
	"time"
)

// Registration func reading data from connection
// checks data for included errors
// checks name for valid, using
// if all is normal prints history of messages to user
// then add connection to pool of active connections
// and send message about new user to other users
func registration(conn *net.Conn, serv *TCPServer) error {

	channels := serv.Chans
	conf := serv.Config

	for {

		rd := bufio.NewReader(*conn)

		text, err := rd.ReadString('\n')
		if err != nil {
			return err
		}

		text = strings.SplitN(text, "\n", 1)[0]
		text = text[:len(text)-1]
		if !utils.ValidName(text) {
			if _, err = (*conn).Write([]byte("[INCORRECT NAME]:")); err != nil {
				return err
			}
			continue
		}

		busy := false
		for _, v := range *serv.Pool {
			if v.Name == text {
				if _, err = (*conn).Write([]byte("[NAME ALREADY IN USE]:")); err != nil {
					return err
				}
				busy = true
				break
			}
		}
		if busy {
			continue
		}

		//(*conn).Write([]byte(*history))
		connection := &models.Connection{
			Name:    text,
			Action:  time.Now(),
			Room:    1,
			LastMsg: make([]string, conf.Rooms+1),
			//History: &History{},
		}

		(*serv.Pool)[*conn] = connection

		m := &models.Message{
			Author: text,
			In:     true,
			Room:   connection.Room,
			Conn:   conn,
		}

		channels.Msgs <- m

		if conf.Gui {
			channels.GuiMsgs <- m
		}

		if _, err = (*conn).Write([]byte("[FOR HELP USAGE]: --help\n")); err != nil {
			return err
		}
		//log.Printf("user use name: %v from: %v", text, (*conn).RemoteAddr().String())

		if err := PrintHistory(conn, serv); err != nil {
			return err
		}

		channels.Conns <- *conn
		//*conns <- *conn
		break
	}
	return nil
}

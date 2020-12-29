package tcpserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"netcat/internal/app/models"
	"netcat/internal/app/utils"
	"regexp"
	"strconv"
	"time"
)

// NewServer returns tcp server ...
func NewServer(conf *Config) (Server, error) {

	channels := &Channels{
		Conns:   make(chan net.Conn),
		Dconns:  make(chan *models.Message),
		Msgs:    make(chan *models.Message),
		GuiMsgs: make(chan *models.Message),
		Errors:  make(chan error),
	}

	//pool is storage of connections
	pool := make(map[net.Conn]*models.Connection)

	h, err := HistoryCreate("history", conf.Rooms)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		Protocol: conf.Protocol,
		Port:     conf.Port,
		Pool:     &pool,
		History:  h,
		Config:   conf,
		Chans:    channels,
	}, nil
}

// Run ...
func (t *TCPServer) Run() error {

	var err error

	t.Server, err = net.Listen(t.Protocol, t.Port)
	if err != nil {
		return err
	}
	fmt.Printf("listening at localhost: %v\n", t.Config.Port)

	go idleTimer(t)

	if t.Config.Gui {
		go guiTerm(t)
	}

	go func() {

		for {

			c, err := t.Server.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}

			log.Printf("Connected from: %v", c.RemoteAddr().String())

			go func(conn *net.Conn) {

				//limit of maximum active connections
				if len(*t.Pool) >= t.Config.Conns {
					if _, err = (*conn).Write([]byte("[SORRY, max connections online. try again later]\n")); err != nil {
						log.Println(err)
					}
					if err = (*conn).Close(); err != nil {
						log.Println(err)
					}
				}

				if err := utils.Hello(conn); err != nil {
					log.Fatal(err)
				}

				//utils.Hello(conn)

				go regTimer(conn, t)

				if err := registration(conn, t); err != nil {
					log.Println(err)
					if err = (*conn).Close(); err != nil {
						log.Println(err)
					}

				}

				// if err := utils.Registration(conn, t); err != nil {
				// 	log.Println(err)
				// 	if err = (*conn).Close(); err != nil {
				// 		log.Println(err)
				// 	}
				// }

			}(&c)

		}

	}()

	r := regexp.MustCompile(`^-`)

	for {
		select {

		case c := <-t.Chans.Conns:

			go func(cn net.Conn) {

				for {
					data := (*t.Pool)[cn]

					pre := utils.PreText(data.Name)

					if _, err = cn.Write([]byte(pre)); err != nil {
						t.Chans.Errors <- err
					}

					rd := bufio.NewReader(cn)

					text, err := rd.ReadString('\n')
					if err != nil {
						log.Printf("Connection closing by %v from: %v", err.Error(), cn.RemoteAddr().String())
						break
					}
					text, err = utils.TextFilter(text)
					if err == nil {

						(*t.Pool)[cn].Action = time.Now()

						if r.MatchString(text) {

							if err := FlagCheck(&text, &cn, t); err != nil {
								t.Chans.Errors <- err
							}

						} else {
							t.Chans.Msgs <- &models.Message{
								Author: data.Name,
								Text:   text,
								Room:   data.Room,
								Conn:   &cn,
							}
						}

					}
				}
				//dconns <- &utils.Message{
				t.Chans.Dconns <- &models.Message{
					Conn: &cn,
					Out:  true,
				}
			}(c)

		case m := <-t.Chans.Msgs:

			if m.In {

				for k, v := range *t.Pool {
					if v.Name != m.Author && m.Room == v.Room {

						if _, err = k.Write([]byte("\n" + m.Author + " has join our chat\n")); err != nil {
							t.Chans.Errors <- err
						}

						if _, err = k.Write([]byte(utils.PreText(v.Name))); err != nil {
							t.Chans.Errors <- err
						}

					}
				}
				if t.Config.Gui {
					t.Chans.GuiMsgs <- m
				}
			} else if m.Out || m.Kick {

				reason := " has left our chat\n"
				if m.Kick {
					reason = " has been kicked from our chat for idle\n"
				}

				for k, v := range *t.Pool {
					if v.Name != m.Author && m.Room == v.Room {

						if _, err = k.Write([]byte("\n" + m.Author + reason)); err != nil {
							t.Chans.Errors <- err
						}
						if _, err = k.Write([]byte(utils.PreText(v.Name))); err != nil {
							t.Chans.Errors <- err
						}

					}
				}
				if t.Config.Gui {
					t.Chans.GuiMsgs <- m
				}
			} else if m.Move {

				for k, v := range *t.Pool {
					if m.RoomFrom == v.Room {
						if _, err = k.Write([]byte("\n" + m.Author + " switched to chat room " + strconv.Itoa(m.Room) + "\n")); err != nil {
							t.Chans.Errors <- err
						}
						if _, err = k.Write([]byte(utils.PreText(v.Name))); err != nil {
							t.Chans.Errors <- err
						}

					} else if m.Room == v.Room && m.Author != v.Name {
						if _, err = k.Write([]byte("\n" + m.Author + " join chat room" + "\n")); err != nil {
							t.Chans.Errors <- err
						}
						if _, err = k.Write([]byte(utils.PreText(v.Name))); err != nil {
							t.Chans.Errors <- err
						}

					}

				}
				if t.Config.Gui {
					t.Chans.GuiMsgs <- m
				}
			} else if m.ChangeName != "" {

				for k, v := range *t.Pool {
					if v.Name != m.ChangeName && m.Room == v.Room {
						if _, err = k.Write([]byte("\n" + m.Author + " changed name to " + m.ChangeName + "\n")); err != nil {
							t.Chans.Errors <- err
						}
						if _, err = k.Write([]byte(utils.PreText(v.Name))); err != nil {
							t.Chans.Errors <- err
						}
					}
				}
				if t.Config.Gui {
					t.Chans.GuiMsgs <- m
				}
			} else {

				//	str := "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + m.Author + "]: "

				pre := utils.PreText(m.Author)
				for k, v := range *t.Pool {
					if v.Name != m.Author && m.Room == v.Room {
						_, err := k.Write([]byte("\n" + pre + m.Text + "\n"))
						if err != nil {
							log.Println(err.Error())
						}
						if _, err = k.Write([]byte(utils.PreText(v.Name))); err != nil {
							t.Chans.Errors <- err
						}

					}
				}
				m := &models.Message{
					Text: pre + m.Text,
					Room: m.Room,
					Conn: m.Conn,
				}

				err = ToHistory(m, t)
				if err != nil {
					log.Printf(err.Error())
				}

				if t.Config.Gui {
					t.Chans.GuiMsgs <- m
				}

			}

		//case d := <-dconns:
		case d := <-t.Chans.Dconns:

			go func(c *models.Message) {

				data, ok := (*t.Pool)[*c.Conn]
				if ok {
					delete(*t.Pool, *c.Conn)
					if err := (*c.Conn).Close(); err != nil {
						t.Chans.Errors <- err
					}
					//msgs <- &utils.Message{
					t.Chans.Msgs <- &models.Message{
						Author: data.Name,
						Out:    c.Out,
						Kick:   c.Kick,
						Room:   data.Room,
					}

				}
			}(d)
		case e := <-t.Chans.Errors:
			log.Println(e)
		}
	}

	//return nil

}

// Close ...
func (t *TCPServer) Close() error {

	return t.Server.Close()
}

// Listener ...
func (t *TCPServer) Listener() error {

	return nil
}

package tcpserver

import (
	"bufio"
	"fmt"
	"io"
	"kitKat/internal/app/models"
	onlinehistory "kitKat/internal/app/store/onlineHistory"
	"log"
	"net"

	//	onlinehistory "netcat/internal/app/store/onlineHistory"
	"kitKat/internal/app/utils"
	"os"
	"regexp"
	"strings"
	"time"
)

// NewestServer ....
func NewestServer(conf *Config) (Server, error) {

	channels := &NewChannels{
		Conns:   make(chan *net.Conn),
		Dconns:  make(chan *models.Message),
		Msgs:    make(chan *models.Message),
		GuiMsgs: make(chan *models.Message),
		Errors:  make(chan *models.ErrMes),
	}

	store := onlinehistory.New()
	// store := testhistory.New()
	if err := store.History().Create("newStory", conf.Rooms); err != nil {
		return nil, err
	}

	pool, err := store.Pool().Create(conf.Conns, conf.Rooms)
	if err != nil {
		return nil, err
	}

	return &NewTCPServer{
		Pool:   pool,
		Config: conf,
		Chans:  channels,
		Store:  store,
	}, nil

}

// Run ...
func (t *NewTCPServer) Run() error {

	var err error

	if t.Config.Testing {
		_, err := utils.Log("testlogs")
		if err != nil {
			return err
		}

		defer func() {
			//logFile.Close()
			os.RemoveAll("testlogs")
		}()
	}

	t.Server, err = net.Listen(t.Config.Protocol, t.Config.Port)
	if err != nil {
		return err
	}

	fmt.Printf("listening at localhost: %v\n", t.Config.Port)

	t.Working = true

	go func() {
		if err := t.Store.Pool().IdleTimer(&t.Working, &t.Config.TimeIdle, &t.Chans.Dconns); err != nil {
			t.Chans.Errors <- &models.ErrMes{
				From:  "run idle timer",
				Error: err,
			}
		}
	}()

	if t.Config.Gui {
		go newGuiTerm(t)
	}

	go func() {

		for {

			c, err := t.Server.Accept()

			if err != nil {
				t.Chans.Errors <- &models.ErrMes{
					From:  "run server accept",
					Error: err,
				}
				break
			}

			log.Printf("Connected from: %v", c.RemoteAddr().String())

			go func() {
				conn, err := t.Store.Pool().Handler(&c, &t.Config.TimeReg)
				if err != nil {

					t.Chans.Errors <- &models.ErrMes{
						From:  fmt.Sprintf("%v", c.RemoteAddr().String()),
						Error: err,
					}
					return
				}

				if conn == nil {
					return
				}

				if t.Config.Gui {
					t.Chans.GuiMsgs <- &models.Message{
						Conn:   conn.Conn,
						In:     true,
						Author: conn.Name,
					}
				}

				t.Chans.Conns <- &c
			}()

		}

	}()

	//_ = regexp.MustCompile(`^-`)

	for {

		select {

		case c, ok := <-t.Chans.Conns:
			if !ok {
				break
			}

			go func() {

				reg := regexp.MustCompile(`^-`)

				for {

					conn, err := t.Store.Pool().Get(c)
					if err != nil {
						t.Chans.Errors <- &models.ErrMes{
							From:  "run selector pool get",
							Error: err,
						}
						break
					}

					if err = conn.Prefix(); err != nil {
						t.Chans.Errors <- &models.ErrMes{
							From:  "run selector prefix",
							Error: err,
						}
						break
					}
					rd := bufio.NewReader(*conn.Conn)

					text, err := rd.ReadString('\n')

					if err != nil {
						if err != io.EOF {
							t.Chans.Errors <- &models.ErrMes{
								From:  "run select readstring",
								Error: err,
							}
							return
						}
						break
					}

					text, err = utils.TextFilter(text)
					if err != nil {

						continue
					}
					conn.Action = time.Now()

					mesg := &models.Message{
						Author: conn.Name,
						Conn:   conn.Conn,
						Text:   text,
						Room:   conn.Room,
					}

					if reg.MatchString(text) {

						if err := t.Store.Pool().Change(mesg); err != nil {
							t.Chans.Errors <- &models.ErrMes{
								From:  "run select pool change",
								Error: err,
							}
							continue
						}

					}

					if t.Config.Gui {
						t.Chans.GuiMsgs <- mesg
					}

					if mesg.ChangeName != "" {
						if err := t.Store.Pool().RenameUser(mesg); err != nil {
							t.Chans.Errors <- &models.ErrMes{
								From:  "run select rename user",
								Error: err,
							}
						}
						continue
					}

					if mesg.Move {
						if err := t.Store.Pool().ChangeRoom(mesg); err != nil {
							t.Chans.Errors <- &models.ErrMes{
								From:  "run select change room",
								Error: err,
							}
						}
						continue
					}
					if mesg.CheckOnline {
						if err := t.Store.Pool().WhoIsOnline(mesg); err != nil {
							t.Chans.Errors <- &models.ErrMes{
								From:  "run select whoOnline",
								Error: err,
							}
						}
						continue
					}

					if mesg.Help {
						if _, err := (*mesg.Conn).Write([]byte(utils.HelpOptions + "\n")); err != nil {
							t.Chans.Errors <- &models.ErrMes{
								From:  "run select help options",
								Error: err,
							}
						}
						continue
					}

					t.Chans.Msgs <- mesg

				}

				connect, err := t.Store.Pool().Get(c)
				if err != nil {
					t.Chans.Errors <- &models.ErrMes{
						From:  "run conn routine get",
						Error: err,
					}
					return
				}

				t.Chans.Dconns <- &models.Message{
					Author: connect.Name,
					Conn:   connect.Conn,
					Out:    true,
					Room:   connect.Room,
				}

			}()

		case m, ok := <-t.Chans.Msgs:
			if !ok {

				break
			}

			go func() {

				m.Text = utils.Prefix(m.Author) + m.Text

				if err := t.Store.Pool().SendMsg(m); err != nil {
					t.Chans.Errors <- &models.ErrMes{
						From:  "run select send message",
						Error: err,
					}
				}
			}()

		case m, ok := <-t.Chans.Dconns:

			if !ok {
				break
			}
			go func() {

				if err := t.Store.Pool().RemoveUser(m); err != nil {
					t.Chans.Errors <- &models.ErrMes{
						From:  "run select remove user",
						Error: err,
					}
				}
				if t.Config.Gui {
					t.Chans.GuiMsgs <- m
				}
			}()

		case e, ok := <-t.Chans.Errors:

			if !ok {
				break
			}
			if !t.Config.Testing {
				go func() {
					log.Printf("error from: %v; error: %v\n", e.From, e.Error.Error())
				}()
			}

		}

		if !t.Working {
			break
		}

	}

	log.Printf("func RUN stoped")
	return nil
}

// Close ...
func (t *NewTCPServer) Close() error {

	//fmt.Fprintf(os.Stdout, "server closing...")
	fmt.Println("server closing...")

	if err := t.Server.Close(); err != nil {
		t.Chans.Errors <- &models.ErrMes{
			From:  "close serv close",
			Error: err,
		}
	}

	mes := &models.Message{
		All: true,
	}

	str := fmt.Sprintf("\n[SERVER MESSAGE]: server will stop after %v", t.Config.TimeOff)

	mes.Text = str

	t.Chans.Msgs <- mes

	time.Sleep(t.Config.TimeOff)

	t.Working = false

	close(t.Chans.Conns)
	close(t.Chans.Dconns)
	close(t.Chans.Msgs)
	close(t.Chans.GuiMsgs)
	//close(t.Chans.Errors)

	fmt.Fprintf(os.Stdout, "all closed\n")

	return nil
}

// Listener ...
func (t *NewTCPServer) Listener() error {

	exitMes := regexp.MustCompile(`^quit`)
	confMes := regexp.MustCompile(`^(y|yes)`)

	r := exitMes

	confirm := false

	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		text = strings.ToLower(text)

		if confirm && r.MatchString(text) {
			fmt.Fprintf(os.Stdout, "ok\n")
			break
		}

		if r.MatchString(text) {
			fmt.Fprintf(os.Stdout, "do you want to quit?\nY/N:")
			r = confMes
			confirm = true
			continue
		}

		r = exitMes
		confirm = false
		fmt.Println("for quit from prog send \"quit\"")
	}

	if err := t.Server.Close(); err != nil {
		return err
	}

	return nil
}

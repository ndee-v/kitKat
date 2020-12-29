package tcpserver

import (
	"errors"
	"fmt"
	"kitKat/internal/app/models"
	"sort"
	"strconv"

	"github.com/jroimartin/gocui"
)

// updaterGui ...
func updaterGui(g *gocui.Gui, s *TCPServer) {

	//time.Sleep(time.Second * 1)

	for {

		select {

		case m := <-(s.Chans.GuiMsgs):

			if m.In || m.Out || m.Kick || m.Move || m.ChangeName != "" {
				g.Update(func(gui *gocui.Gui) error {
					online := make([]int, s.Config.Rooms+1)
					users := make([][]string, s.Config.Rooms+1)

					for _, v := range *s.Pool {
						online[v.Room]++
						users[v.Room] = append(users[v.Room], v.Name)
					}

					view, err := g.View("rooms")
					if err != nil {
						s.Chans.Errors <- err
						return nil
					}

					view.Clear()

					view.Title = "online " + strconv.Itoa(len(*s.Pool))

					for i := 1; i <= s.Config.Rooms; i++ {
						if _, err = fmt.Fprintln(view, "room_"+strconv.Itoa(i)+"_["+strconv.Itoa(online[i])+"]"); err != nil {
							s.Chans.Errors <- err
						}
					}

					for i := 1; i <= s.Config.Rooms; i++ {

						sort.Strings(users[i])

						view, err = g.View("room_" + strconv.Itoa(i))

						if err != nil {
							s.Chans.Errors <- err
							return nil
						}

						view.Clear()

						for _, v := range users[i] {
							if _, err := view.Write([]byte(v + "\n")); err != nil {
								s.Chans.Errors <- err
								return nil
							}
						}

					}

					return nil
				})
			} else {
				g.Update(func(gui *gocui.Gui) error {

					view, err := g.View("chat_" + strconv.Itoa(m.Room))
					if err != nil {
						s.Chans.Errors <- err
						return nil
					}

					if _, err := view.Write([]byte(m.Text + "\n")); err != nil {
						s.Chans.Errors <- err
						return nil
					}

					return nil
				})
			}

			if m.In || m.Out || m.Kick || m.ChangeName != "" {
				g.Update(func(gui *gocui.Gui) error {
					if !m.In {
						if err := g.DeleteView("user: " + m.Author); err != nil {
							s.Chans.Errors <- err
							return nil
						}
					}
					if m.In || m.ChangeName != "" {
						X, Y := g.Size()

						name := m.Author
						if m.ChangeName != "" {
							name = m.ChangeName
						}

						v, err := g.SetView("user: "+name, 21+(X-20)/5+1, 3, 21+(X-20)/5+40, Y-3)
						if err != gocui.ErrUnknownView {
							s.Chans.Errors <- err
							return nil
						}
						if _, err = g.SetViewOnBottom("user: " + name); err != nil {
							s.Chans.Errors <- err
							return nil
						}

						data, ok := (*s.Pool)[*m.Conn]
						if !ok {
							return nil
						}

						v.Title = "user_information"
						//v.Highlight = true
						v.BgColor = gocui.ColorRed
						v.FgColor = gocui.ColorWhite
						if _, err := fmt.Fprintf(v, " name: %s\n"+
							" Remote addr: %s\n"+
							" Local addr: %s\n",
							data.Name, (*m.Conn).RemoteAddr().String(), (*m.Conn).LocalAddr().String()); err != nil {
							s.Chans.Errors <- err
							return nil
						}

						if _, err := g.SetViewOnBottom("user: " + m.ChangeName); err != nil {
							s.Chans.Errors <- err
							return nil
						}

					}

					return nil
				})
			}

		}
	}
}

// updaterGui ...
func newUpdaterGui(g *gocui.Gui, s *NewTCPServer) error {

	//time.Sleep(time.Second * 1)

	for s.Working {

		select {

		case m, ok := <-(s.Chans.GuiMsgs):

			if !ok {
				return errors.New("error from guimsg channel")
			}

			//left side of terminal && right side

			if m.In || m.Out || m.Kick || m.Move || m.ChangeName != "" {
				g.Update(func(gui *gocui.Gui) error {
					online := make([]int, s.Config.Rooms+1)
					users := make([][]string, s.Config.Rooms+1)

					for _, v := range s.Pool {
						online[v.Room]++
						users[v.Room] = append(users[v.Room], v.Name)
					}

					view, err := g.View("rooms")
					if err != nil {
						s.Chans.Errors <- &models.ErrMes{
							From:  "newGuiUpdater first case View",
							Error: err,
						}
						return nil
					}

					view.Clear()

					view.Title = "online " + strconv.Itoa(len(s.Pool))

					for i := 1; i <= s.Config.Rooms; i++ {
						if _, err = fmt.Fprintln(view, "room_"+strconv.Itoa(i)+"_["+strconv.Itoa(online[i])+"]"); err != nil {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater first case Fprint",
								Error: err,
							}
						}
					}

					for i := 1; i <= s.Config.Rooms; i++ {

						sort.Strings(users[i])

						view, err = g.View("room_" + strconv.Itoa(i))

						if err != nil {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater first case View2",
								Error: err,
							}
							return nil
						}

						view.Clear()

						for _, v := range users[i] {
							if _, err := view.Write([]byte(v + "\n")); err != nil {
								s.Chans.Errors <- &models.ErrMes{
									From:  "newGuiUpdater first case Write in view",
									Error: err,
								}
								return nil
							}
						}

					}

					return nil
				})
			} else if !m.Help && !m.CheckOnline {
				g.Update(func(gui *gocui.Gui) error {

					view, err := g.View("chat_" + strconv.Itoa(m.Room))
					if err != nil {
						s.Chans.Errors <- &models.ErrMes{
							From:  "newGuiUpdater first case update chat_ view",
							Error: err,
						}
						return nil
					}

					if _, err := view.Write([]byte(m.Text + "\n")); err != nil {
						s.Chans.Errors <- &models.ErrMes{
							From:  "newGuiUpdater first case update chat_ write to view",
							Error: err,
						}
						return nil
					}

					return nil
				})
			}

			//special views for each users (info, address etc.)

			if m.In || m.Out || m.Kick || m.ChangeName != "" {
				g.Update(func(gui *gocui.Gui) error {
					if !m.In {
						if err := g.DeleteView("user: " + m.Author); err != nil {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater first case update deleteView",
								Error: err,
							}
							return nil
						}
					}
					if m.In || m.ChangeName != "" {
						X, Y := g.Size()

						name := m.Author
						if m.ChangeName != "" {
							name = m.ChangeName
						}

						v, err := g.SetView("user: "+name, 21+(X-20)/5+1, 3, 21+(X-20)/5+40, Y-3)
						if err != gocui.ErrUnknownView {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater first case update setView user",
								Error: err,
							}
							return nil
						}
						if _, err = g.SetViewOnBottom("user: " + name); err != nil {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater first case update setViewOnBottom user",
								Error: err,
							}
							return nil
						}

						data, ok := (s.Pool)[m.Conn]
						if !ok {
							return nil
						}

						v.Title = "user_information"
						//v.Highlight = true
						v.BgColor = gocui.ColorRed
						v.FgColor = gocui.ColorWhite
						if _, err := fmt.Fprintf(v, " name: %s\n"+
							" Remote addr: %s\n"+
							" Local addr: %s\n",
							data.Name, (*m.Conn).RemoteAddr().String(), (*m.Conn).LocalAddr().String()); err != nil {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater FprintF",
								Error: err,
							}
							return nil
						}

						if _, err := g.SetViewOnBottom("user: " + name); err != nil {
							s.Chans.Errors <- &models.ErrMes{
								From:  "newGuiUpdater setViewOnBottom",
								Error: err,
							}
							return nil
						}

					}

					return nil
				})
			}

		}
	}
	return nil
}

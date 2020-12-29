package tcpserver

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jroimartin/gocui"
)

// Layout ...
func (gs *GuiScheme) Layout(g *gocui.Gui) error {

	// size of terminal

	X, Y := g.Size()

	//initialize backGround
	v, err := g.SetView("bg", 0, 0, X, Y)
	if err != gocui.ErrUnknownView {
		return err
	}
	if v != nil {
		v.Clear()
		v.BgColor = gocui.ColorBlack
		v.Frame = false
	}

	//initialize sign of program
	v, err = g.SetView("sign", 1, 0, X-1, 3)
	if err != gocui.ErrUnknownView {
		return err
	}
	if v != nil {
		v.Clear()
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorWhite
		v.Frame = false
		if _, err := fmt.Fprintf(v, "netCat program /// using like TCP-chat. listening at localhost: %v", gs.Port); err != nil {
			return err
		}
	}

	//initialize close button
	v, err = g.SetView("exit", X-15, 0, X-1, 2)
	if err != gocui.ErrUnknownView {
		return err
	}
	if v != nil {
		v.Clear()
		v.BgColor = gocui.ColorRed
		v.FgColor = gocui.ColorWhite
		v.Frame = false
		if _, err = v.Write([]byte("  CLOSE GUI  ")); err != nil {
			return err
		}
		if _, err = g.SetViewOnTop("exit"); err != nil {
			return err
		}
	}

	//initialize window for chat rooms
	v, err = g.SetView("rooms", 2, 3, 19, 3+gs.Rooms+1)
	if err != gocui.ErrUnknownView {

		return err
	}
	if v != nil {
		v.Title = "online 0"
		v.Highlight = true
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorWhite
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for i := 1; i <= gs.Rooms; i++ {
			if _, err := fmt.Fprintln(v, "room_"+strconv.Itoa(i)); err != nil {
				return err
			}
		}
	}

	//initialize different windows for each room ( to show users inside) and chats for each room
	for i := 1; i <= gs.Rooms; i++ {
		name := "room_" + strconv.Itoa(i)
		v, err = g.SetView(name, 21, 3, 20+(X-20)/5, 3+(gs.Users*2))
		if err != gocui.ErrUnknownView {
			return err
		}
		if v != nil {
			v.Title = name
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorYellow
			v.SelBgColor = gocui.ColorBlue
			if _, err = g.SetViewOnBottom(name); err != nil {
				log.Fatal(err)
			}
		}

		name = "chat_" + strconv.Itoa(i)
		v, err = g.SetView(name, 21+(X-20)/5+1, 3, X-2, Y-2)
		if err != gocui.ErrUnknownView {
			return err
		}
		if v != nil {
			v.Title = name
			v.BgColor = gocui.ColorBlue
			v.FgColor = gocui.ColorWhite
			v.Wrap = true
			v.Autoscroll = true
			if _, err := g.SetViewOnBottom(name); err != nil {
				return err
			}
		}

	}

	if _, err := g.SetViewOnBottom("sign"); err != nil {
		return err
	}
	if _, err := g.SetViewOnBottom("bg"); err != nil {
		return err
	}
	return nil
}

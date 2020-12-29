package tcpserver

import (
	"strings"

	"github.com/jroimartin/gocui"
)

// showRoom ...
func showRoom(g *gocui.Gui, v *gocui.View) error {

	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	l, err = v.Line(cy)
	if err != nil {
		//l = ""
		return nil
	}

	if l != "" {
		l = strings.Split(l, "_")[1]
		//if len(d)==2 {
		if _, err = g.SetViewOnTop("room_" + l); err != nil {
			return err
		}
		if _, err = g.SetViewOnTop("chat_" + l); err != nil {
			return err
		}
	}

	return err
}

// showUser ...
func showUser(g *gocui.Gui, v *gocui.View) error {

	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	l, err = v.Line(cy)
	if err != nil {
		//l = ""
		return nil
	}

	if l != "" {

		if _, err = g.SetViewOnTop("user: " + l); err != nil {
			return err
		}

	}

	return nil
}

// showChat ...
func showChat(g *gocui.Gui, v *gocui.View) error {

	if _, err := g.SetViewOnTop(v.Name()); err != nil {
		return err
	}

	return nil
}

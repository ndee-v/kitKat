package tcpserver

import (
	"kitKat/internal/app/models"
	"log"

	//"netcat/internal/app/models"
	//"netcat/internal/app/models"
	//"netcat/internal/app/models"
	//	"netcat/internal/app/models"
	"strconv"
	"time"

	"github.com/jroimartin/gocui"
)

// guiTerm ...
func guiTerm(serv *TCPServer) {

	time.Sleep(time.Second * 1)

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {

		serv.Chans.Errors <- err
		return
	}
	defer gui.Close()

	gui.SetManager(newGuiScheme(serv.Config))

	gui.Mouse = true
	gui.Cursor = true

	if err = gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		//log.Panicln(err)
		serv.Chans.Errors <- err
		//channels.Errors <- err
		return
	}

	if err = gui.SetKeybinding("rooms", gocui.MouseLeft, gocui.ModNone, showRoom); err != nil {
		//channels.Errors <- err
		serv.Chans.Errors <- err
		return
	}

	for i := 1; i <= serv.Config.Rooms; i++ {
		name := "room_" + strconv.Itoa(i)
		if err := gui.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, showUser); err != nil {
			//channels.Errors <- err
			serv.Chans.Errors <- err
			return
		}
		name = "chat_" + strconv.Itoa(i)
		if err := gui.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, showChat); err != nil {
			//channels.Errors <- err
			serv.Chans.Errors <- err
			return
		}
	}

	if err = gui.SetKeybinding("exit", gocui.MouseLeft, gocui.ModNone, closeGui); err != nil {
		serv.Chans.Errors <- err
		//channels.Errors <- err
		return
	}

	go updaterGui(gui, serv)

	if err = gui.MainLoop(); err != nil && err != gocui.ErrQuit {

		serv.Chans.Errors <- err

		return

	}
	return
}

// newGuiTerm ...
func newGuiTerm(serv *NewTCPServer) {

	log.Println("newGuiTerm starts")

	time.Sleep(time.Second * 1)

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {

		serv.Chans.Errors <- &models.ErrMes{
			From:  "newGuiTerm NewGui",
			Error: err,
		}
		return
	}
	defer gui.Close()

	gui.SetManager(newGuiScheme(serv.Config))

	gui.Mouse = true
	gui.Cursor = true

	if err = gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		//log.Panicln(err)
		serv.Chans.Errors <- &models.ErrMes{
			From:  "newGuiTerm setKeyBinding quit",
			Error: err,
		}
		//channels.Errors <- err
		return
	}

	if err = gui.SetKeybinding("rooms", gocui.MouseLeft, gocui.ModNone, showRoom); err != nil {
		//channels.Errors <- err
		serv.Chans.Errors <- &models.ErrMes{
			From:  "newGuiTerm setKeyBinding rooms",
			Error: err,
		}
		return
	}

	for i := 1; i <= serv.Config.Rooms; i++ {
		name := "room_" + strconv.Itoa(i)
		if err := gui.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, showUser); err != nil {
			//channels.Errors <- err
			serv.Chans.Errors <- &models.ErrMes{
				From:  "newGuiTerm setKeyBinding showUser",
				Error: err,
			}
			return
		}
		name = "chat_" + strconv.Itoa(i)
		if err := gui.SetKeybinding(name, gocui.MouseLeft, gocui.ModNone, showChat); err != nil {
			//channels.Errors <- err
			serv.Chans.Errors <- &models.ErrMes{
				From:  "newGuiTerm setKeyBinding showChat",
				Error: err,
			}
			return
		}
	}

	/*
		if err = gui.SetKeybinding("exit", gocui.MouseLeft, gocui.ModNone, closeGui); err != nil {
			serv.Chans.Errors <- err
			//channels.Errors <- err
			return
		}


	*/
	if err = gui.SetKeybinding("exit", gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {

		log.Printf("GUI closed by exit btn2")
		g.Close()

		go func() {
			if err := serv.Listener(); err != nil {
				log.Fatal(err)
			}

			if err := serv.Close(); err != nil {
				log.Fatal(err)
			}

		}()

		return nil
	}); err != nil {
		serv.Chans.Errors <- &models.ErrMes{
			From:  "newGuiTerm setKeyBinding exitButton",
			Error: err,
		}
		return
	}

	go func() {

		if err := newUpdaterGui(gui, serv); err != nil {
			serv.Chans.Errors <- &models.ErrMes{
				From:  "newGuiTerm newUpdaterGui",
				Error: err,
			}
			return
		}

	}()
	//go newUpdaterGui(gui, serv)

	//go newUpdaterGui(gui, serv)

	if err = gui.MainLoop(); err != nil && err != gocui.ErrQuit {

		serv.Chans.Errors <- &models.ErrMes{
			From:  "newGuiTerm MainLoop",
			Error: err,
		}

		return

	}
	return
}

// quit func ...
func quit(_ *gocui.Gui, _ *gocui.View) error {
	log.Printf("GUI closed by ^C")
	return gocui.ErrQuit
}

// closeGui func ...
func closeGui(g *gocui.Gui, _ *gocui.View) error {

	log.Printf("GUI closed by exit btn")
	g.Close()

	//go func() {
	//	if err := s.Listener(); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if err := s.Close(); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//}()

	return nil
}

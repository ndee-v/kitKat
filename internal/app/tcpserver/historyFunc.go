package tcpserver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"netcat/internal/app/models"
	"netcat/internal/app/utils"

	// "netcat/internal/app/models"
	// "netcat/internal/app/utils"

	// "netcat/internal/app/models"
	// "netcat/internal/app/utils"
	"os"
	"strconv"
)

// HistoryCreate make history files for each room of chat
func HistoryCreate(s string, num int) (*models.History, error) {

	if s == "" || num <= 0 {
		return nil, errors.New("invalid params")
	}

	src, err := os.Stat(s)
	if err != nil {
		if err = os.Mkdir(s, os.ModePerm); err != nil {
			//log.Fatal(err.Error())
			return nil, err
		}
	} else {
		if !src.IsDir() {
			//log.Fatal("history is not a directory")
			return nil, errors.New(s + " is not a directory")
		}
	}

	arr := make([]*os.File, num+1)
	ind := 1

	for i := 1; i <= num; i++ {

		//	name := PreText("room_" + strconv.Itoa(i))

		name := utils.PreText("room_" + strconv.Itoa(i))
		h, err := os.Create("./" + s + "/" + name + ".txt")
		if err != nil {
			return nil, err
		}

		arr[ind] = h
		ind++
	}

	return &models.History{
		Dir:  "./" + s + "/",
		Room: arr,
	}, nil
}

// PrintHistory func prints chat history
func PrintHistory(cn *net.Conn, serv *TCPServer) error {

	conn, ok := (*serv.Pool)[*cn]
	if !ok {
		//fmt.Println("can not found connection in pool")
		return errors.New("can not found connection in pool")
	}

	f := serv.History.Room[conn.Room]
	st, _ := f.Stat()
	file, err := os.OpenFile(serv.History.Dir+st.Name(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		serv.Chans.Errors <- file.Close()
	}()

	rd := bufio.NewReader(file)

	toPrint := false

	lastMess := conn.LastMsg[conn.Room]

	if lastMess == "" {
		toPrint = true
	}

	var lastPrinted string

	for {

		line, err := rd.ReadString('\n')

		if err != nil {
			if err == io.EOF {

				break
			}

			//log.Fatalf("read file line error: %v", err)
			return err
		}

		if lastMess == line {
			toPrint = true
			continue
		}
		if toPrint {
			if _, err := (*cn).Write([]byte(line)); err != nil {
				return err
			}
			lastPrinted = line
		}
	}
	(*serv.Pool)[*cn].LastMsg[conn.Room] = lastPrinted
	return nil

}

// ToHistory func save message to chat history
func ToHistory(msg *models.Message, t *TCPServer) error {

	hFile := t.History.Room[msg.Room]

	if _, err := hFile.WriteString(msg.Text + "\n"); err != nil {
		return err
	}

	for _, v := range *t.Pool {
		if msg.Room == v.Room {
			v.LastMsg[v.Room] = msg.Text + "\n"
		}
	}
	return nil
}

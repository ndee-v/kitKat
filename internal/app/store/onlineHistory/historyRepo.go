package onlinehistory

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"netcat/internal/app/models"
	"netcat/internal/app/store"
	"os"
	"regexp"
	"strconv"
	"time"
)

//HistoryRepo ...
type HistoryRepo struct {
	Dir     string
	Folder  string
	Rooms   int
	History map[int]*os.File
}

// Create files for save history
func (h *HistoryRepo) Create(dir string, num int) error {

	pattern := regexp.MustCompile(`^\w{1,10}$`)
	if !pattern.MatchString(dir) {
		return store.ErrDirString
	}

	if num <= 0 {
		return store.ErrNumRooms
	}

	src, err := os.Stat(dir)
	if err != nil {
		if err = os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
	} else {
		if !src.IsDir() {
			return store.ErrNotADir
		}
	}

	folder := time.Now().Format(time.RFC3339)
	src, err = os.Stat(dir + "/" + folder)
	if err != nil {
		if err = os.Mkdir(dir+"/"+folder, os.ModePerm); err != nil {
			return err
		}
	} else {
		if !src.IsDir() {
			return store.ErrNotADir
		}
	}

	h.Dir = dir
	h.Folder = folder
	h.Rooms = num

	for i := 1; i <= num; i++ {

		f, err := os.Create(fmt.Sprintf("./%v/%v/room_%v.txt", dir, folder, strconv.Itoa(i)))
		if err != nil {
			return err
		}
		h.History[i] = f
	}
	return nil
}

// AddInto adds message to history
func (h *HistoryRepo) AddInto(m *models.Message) error {

	file, ok := h.History[m.Room]
	if !ok {
		return store.ErrRoomNotFound
	}

	if _, err := file.WriteString(m.Text + "\n"); err != nil {
		return err
	}

	return nil
}

// PrintToConn ...
func (h *HistoryRepo) PrintToConn(c *models.Connection) error {

	if c == nil {
		return store.ErrConnNotInit
	}

	if c.Conn == nil {
		return store.ErrConnNotInit
	}

	f, ok := h.History[c.Room]
	if !ok {
		return store.ErrInvalidInput
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fmt.Sprintf("./%v/%v/%v", h.Dir, h.Folder, stat.Name()), os.O_RDONLY, os.ModePerm)

	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error file close: %v", err)
		}
	}()

	rd := bufio.NewReader(file)

	toPrint := false

	lastMess := c.LastMsg[c.Room]

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
			return err
		}

		if lastMess == line {
			toPrint = true
			continue
		}
		if toPrint {
			if _, err := (*c.Conn).Write([]byte(line)); err != nil {
				return err
			}
			lastPrinted = line
		}
	}
	c.LastMsg[c.Room] = lastPrinted

	return nil
}

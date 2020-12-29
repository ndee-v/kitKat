package testhistory

import (
	// "netcat/internal/app/models"
	// "netcat/internal/app/store"
	"kitKat/internal/app/models"
	"kitKat/internal/app/store"
	"regexp"
	"time"
)

// HistoryRepo ...
type HistoryRepo struct {
	Dir     string
	Folder  string
	Rooms   int
	History map[int][]string
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
	h.Dir = dir
	h.Folder = time.Now().Format(time.RFC3339)
	h.Rooms = num

	for i := 1; i <= num; i++ {
		h.History[i] = []string{}
	}
	return nil
}

// AddInto adds message to history
func (h *HistoryRepo) AddInto(m *models.Message) error {

	// if h.History == nil {
	// 	return store.ErrHistoryNotExist
	// }

	file, ok := h.History[m.Room]
	if !ok {
		return store.ErrRoomNotFound
	}

	file = append(file, m.Text)

	h.History[m.Room] = file

	return nil
}

// PrintToConn ...
func (h *HistoryRepo) PrintToConn(c *models.Connection) error {

	// if h.History == nil {
	// 	return store.ErrHistoryNotExist
	// }
	if c == nil {
		return store.ErrConnNotInit
	}

	if c.Conn == nil {
		return store.ErrConnNotInit
	}

	array, ok := h.History[c.Room]
	if !ok {
		return store.ErrHistoryNotExist
	}

	toPrint := false

	lastMess := c.LastMsg[c.Room]

	if lastMess == "" {
		toPrint = true
	}

	var lastPrinted string

	for _, line := range array {

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

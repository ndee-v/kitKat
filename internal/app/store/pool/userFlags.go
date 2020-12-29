package pool

import (
	"fmt"
	"netcat/internal/app/models"
	"netcat/internal/app/store"
	"netcat/internal/app/utils"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Change ...
func (r *Repo) Change(m *models.Message) error {

	if m == nil {
		return store.ErrInvalidInput
	}

	var (
		//help       = regexp.MustCompile(`(^--help)|(^-h)`)
		changeName = regexp.MustCompile(`^--name `)
		changeRoom = regexp.MustCompile(`^--room`)
		online     = regexp.MustCompile(`^--online`)
		//temp       = regexp.MustCompile(`^--`)
		//err        error
	)

	if changeName.MatchString(m.Text) {
		newName := strings.Split(m.Text, "--name ")[1]
		if newName != "" {
			m.ChangeName = newName
			return nil
		}
	}

	if changeRoom.MatchString(m.Text) {

		re := regexp.MustCompile(`^--room \d+$`)

		if re.MatchString(m.Text) {
			r := strings.Split(m.Text, "--room ")[1]

			room, err := strconv.Atoi(r)
			if err == nil {
				m.Move = true
				m.RoomFrom = m.Room
				m.Room = room
				return nil
			}
		}
	}

	if online.MatchString(m.Text) {
		m.CheckOnline = true
		return nil
	}

	m.Help = true

	return nil
}

//RenameUser ...
func (r *Repo) RenameUser(m *models.Message) error {

	if m == nil || m.ChangeName == "" {
		return store.ErrInvalidInput
	}

	if m.Conn == nil {
		return store.ErrConnNotInit
	}

	if !utils.ValidName(m.ChangeName) {
		if _, err := (*m.Conn).Write([]byte(utils.IncorrectName + "\n")); err != nil {
			return err
		}
		return nil
	}

	for _, val := range r.Pool {
		if val.Name == m.ChangeName {
			_, err := (*m.Conn).Write([]byte("[NAME ALREADY USING]\n"))
			if err != nil {
				return err
			}
			return nil
		}
	}

	str1 := fmt.Sprintf("[YOUR NAME CHANGED FROM  %v TO %v]", m.Author, m.ChangeName)
	str2 := fmt.Sprintf("\n[USER %v CHANGED HIS NAME TO %v]", m.Author, m.ChangeName)

	for _, val := range r.Pool {
		if val.Name == m.Author {
			val.Name = m.ChangeName
			if err := val.Write(str1); err != nil {
				return err
			}
			continue
		}
		if val.Room == m.Room {
			if err := val.Write(str2); err != nil {
				return err
			}

			if err := val.Prefix(); err != nil {
				return err
			}
		}
	}
	return nil
}

// WhoIsOnline ...
func (r *Repo) WhoIsOnline(m *models.Message) error {

	if m == nil || m.Conn == nil {
		return store.ErrInvalidInput
	}

	var err error

	users := make([][]string, r.Rooms+1)

	//count of all online users
	users[0] = append(users[0], strconv.Itoa(len(r.Pool)))

	//fill all slices if it possible
	for _, v := range r.Pool {
		users[v.Room] = append(users[v.Room], v.Name)
	}

	for i, v := range users {

		if i == 0 {
			if _, err = (*m.Conn).Write([]byte("└┬─── online: " + v[0] + "\n")); err != nil {
				return err
			}
			continue
		}

		sort.Strings(v)

		pre := " ├──"
		if i == len(users)-1 {
			pre = " └──"
		}

		pre2 := "──"
		if len(v) != 0 {
			pre2 = "┬─"
		}

		// time.Sleep(time.Millisecond * 100)
		if _, err = (*m.Conn).Write([]byte(pre + pre2 + " room " + strconv.Itoa(i) + "\n")); err != nil {
			return err
		}

		pre3 := " │ "
		if i == len(users)-1 {
			pre3 = "   "
		}

		// time.Sleep(time.Millisecond * 100)
		for ind, val := range v {

			pre4 := " ├───"
			if ind == len(v)-1 {
				pre4 = " └───"
			}
			if _, err = (*m.Conn).Write([]byte(pre3 + pre4 + " " + val + "\n")); err != nil {
				return err
			}
		}

	}

	return nil
}

//ChangeRoom ...
func (r *Repo) ChangeRoom(m *models.Message) error {

	if m == nil {
		return store.ErrInvalidInput
	}

	if !m.Move {
		return store.ErrInvalidInput
	}

	if m.Conn == nil {
		return store.ErrConnNotInit
	}

	if m.Room == m.RoomFrom {
		if _, err := (*m.Conn).Write([]byte(fmt.Sprintf("[YOU ARE IN ROOM %v]\n", m.Room))); err != nil {
			return err
		}
		return nil
	}
	if m.Room <= 0 || m.Room > r.Rooms {

		if _, err := (*m.Conn).Write([]byte(fmt.Sprintf("[THIS CHAT HAS %v ROOM(s)]\n", r.Capacity))); err != nil {
			return err
		}
		return nil
	}

	str1 := fmt.Sprintf("[YOU SWITCHED FROM ROOM %v to ROOM %v]", m.RoomFrom, m.Room)
	str2 := fmt.Sprintf("\n[USER %v CHANGED ROOM TO %v]", m.Author, m.Room)
	str3 := fmt.Sprintf("\n[USER %v JOIN TO CHAT]", m.Author)

	for key, val := range r.Pool {
		if key == m.Conn {
			r.Pool[key].Room = m.Room

			if err := val.Write(str1); err != nil {
				return err
			}

			if err := r.Store.History().PrintToConn(val); err != nil {
				return err
			}

			continue
		}
		if val.Room == m.RoomFrom {

			if err := val.Write(str2); err != nil {
				return err
			}
			if err := val.Prefix(); err != nil {
				return err
			}
		}
		if val.Room == m.Room {
			if err := val.Write(str3); err != nil {
				return err
			}
			if err := val.Prefix(); err != nil {
				return err
			}
		}
	}

	return nil
}

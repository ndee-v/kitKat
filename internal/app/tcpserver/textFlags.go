package tcpserver

import (
	"log"
	"net"
	"netcat/internal/app/models"
	"netcat/internal/app/utils"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FlagCheck func ...
func FlagCheck(s *string, c *net.Conn, serv *TCPServer) error {

	var err error

	channels := serv.Chans
	conf := serv.Config

	help := regexp.MustCompile(`(^--help)|(^-h)`)
	changeName := regexp.MustCompile(`^--name `)
	changeRoom := regexp.MustCompile(`^--room`)
	online := regexp.MustCompile(`^--online`)
	temp := regexp.MustCompile(`^--`)

	if changeName.MatchString(*s) {
		newName := strings.Split(*s, "--name ")[1]
		if !utils.ValidName(newName) {
			if _, err = (*c).Write([]byte(utils.IncorrectName + "\n")); err != nil {
				return err
			}
		}
		for _, v := range *serv.Pool {
			if v.Name == newName {
				if _, err = (*c).Write([]byte("[NAME IS ALREADY IN CHAT]" + "\n")); err != nil {
					return err
				}
			}
		}

		data := (*serv.Pool)[*c]
		if _, err = (*c).Write([]byte("[YOUR NAME CHANGED FROM \"" + data.Name + "\" to \"" + newName + "\"]\n")); err != nil {
			return err
		}
		log.Printf("user changed name: from %v to %v /// from: %v", data.Name, newName, (*c).RemoteAddr().String())

		// for k, v := range *pool {
		// 	if v.Name != data.Name {
		// 		k.Write([]byte(data.Name + " changed name to " + newName + "\n"))
		// 	}
		// }
		mes := &models.Message{
			Author:     data.Name,
			ChangeName: newName,
			Room:       data.Room,
			Conn:       c,
		}

		data.Name = newName
		(*serv.Pool)[*c] = data

		//(*msgs) <- mes
		channels.Msgs <- mes

		//if conf.Gui {
		//	channels.GuiMsgs <- mes
		//}

	} else if changeRoom.MatchString(*s) {

		re := regexp.MustCompile(`^--room \d+$`)

		if re.MatchString(*s) {

			data := (*serv.Pool)[*c]
			r := strings.Split(*s, "--room ")[1]
			room, _ := strconv.Atoi(r)

			if room == data.Room {

				str := "[YOU ARE IN ROOM]: " + strconv.Itoa(room)
				if _, err = (*c).Write([]byte(str + "\n")); err != nil {
					return err
				}

			} else if 0 < room && room <= conf.Rooms {
				temp := data.Room
				(*serv.Pool)[*c].Room = room

				str := "[YOU ARE SWITCHED TO ROOM]: " + strconv.Itoa(room)
				if _, err = (*c).Write([]byte(str + "\n")); err != nil {
					return err
				}

				log.Printf("user %v changed room to %v /// from: %v\n", data.Name, room, (*c).RemoteAddr().String())

				channels.Msgs <- &models.Message{
					Author:   data.Name,
					Move:     true,
					Room:     room,
					RoomFrom: temp,
				}
				if err := PrintHistory(c, serv); err != nil {
					return err
				}
			} else {
				str := "[CHAT HAS " + strconv.Itoa(conf.Rooms) + " CHAT ROOMS]"
				if _, err = (*c).Write([]byte(str + "\n")); err != nil {
					return err
				}
			}

		}
	} else if online.MatchString(*s) {

		users := make([][]string, conf.Rooms+1)

		//count of all online users
		users[0] = append(users[0], strconv.Itoa(len(*serv.Pool)))

		//fill all slices if it possible
		for _, v := range *serv.Pool {
			users[v.Room] = append(users[v.Room], v.Name)
		}

		for i, v := range users {

			if i == 0 {
				if _, err = (*c).Write([]byte("└┬─── online: " + v[0] + "\n")); err != nil {
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
			// if _, err = ; err != nil {
			// 	return err
			// }
			if _, err = (*c).Write([]byte(pre + pre2 + " room " + strconv.Itoa(i) + "\n")); err != nil {
				return err
			}

			pre3 := " │ "
			if i == len(users)-1 {
				pre3 = "   "
			}

			for ind, val := range v {

				pre4 := " ├───"
				if ind == len(v)-1 {
					pre4 = " └───"
				}
				if _, err = (*c).Write([]byte(pre3 + pre4 + " " + val + "\n")); err != nil {
					return err
				}
			}

		}

	} else if help.MatchString(*s) || temp.MatchString(*s) {

		if _, err = (*c).Write([]byte(utils.HelpOptions + "\n")); err != nil {
			return err
		}

	}

	return nil
}

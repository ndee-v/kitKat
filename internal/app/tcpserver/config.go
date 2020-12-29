package tcpserver

import (
	"errors"
	"fmt"
	"kitKat/internal/app/utils"
	"log"

	//	"netcat/internal/app/utils"

	"regexp"
	"strconv"
	"strings"
	"time"
)

// Prepare method fill config var from arguments
// if its was be recieved
func (c *Config) Prepare(args *[]string) error {

	// default params
	c.TimeReg = time.Second * 30
	c.TimeIdle = time.Second * 180
	c.TimeOff = time.Second * 15
	c.Rooms = 4
	c.Conns = 10
	c.Port = ":8989"
	c.Gui = false

	flags := regexp.MustCompile(`(^g=(true)|(on)$)|(^\d{4,5}$)|(((^tr)|(^ti)|(^r)|(^c))=\d+$)`)

	used := false

	for _, v := range (*args)[1:] {
		if flags.MatchString(v) {

			_, err := strconv.Atoi(v)
			if err == nil {
				c.Port = ":" + v
				continue
			}

			used = true

			d := strings.Split(v, "=")

			val, err := strconv.Atoi(d[1])
			if err != nil && d[0] != "g" {
				return err
			}

			if val <= 0 && d[0] != "g" {
				continue
			}

			switch d[0] {
			case "tr":
				c.TimeReg = time.Duration(val) * time.Second
			case "ti":
				c.TimeIdle = time.Duration(val) * time.Second
			case "r":
				c.Rooms = val
			case "c":
				c.Conns = val
			case "g":
				c.Gui = true
			}
		} else {
			return errors.New("[USAGE]: ./netCat $port")
		}
	}

	if !used && c.Protocol != "" {
		fmt.Println(utils.HelpMessage)
	}

	if c.Protocol != "" {
		log.Printf("\n"+
			"time for registration: %v;\n"+
			"time for idle timeout: %v;\n"+
			"count of rooms : %v;\n"+
			"count of maximum connections: %v;\n"+
			"gui terminal : %v;", c.TimeReg, c.TimeIdle, c.Rooms, c.Conns, c.Gui)
	}

	return nil
}

package main

import (
	"fmt"
	"kitKat/internal/app/tcpserver"
	"kitKat/internal/app/utils"
	"log"
	"os"
)

func main() {

	var logFile *os.File
	var err error

	logFile, err = utils.Log("log")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		if err = logFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	config := &tcpserver.Config{
		Protocol: "tcp",
	}

	if err = config.Prepare(&os.Args); err != nil {
		fmt.Println(err.Error())
		log.Printf("%v", err)
		return
	}

	//s, err := tcpserver.NewServer(config)
	//if err != nil {
	//	log.Fatal(err)
	//}

	s, err := tcpserver.NewestServer(config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if !config.Gui {
			if err := s.Listener(); err != nil {
				log.Fatal(err)
			}

			if err := s.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	log.Printf("TCP Chat stoped")

}

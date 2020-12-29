package utils

import (
	"errors"
	"io"
	"log"
	"os"
)

// Log func using for creation dir and file for logs
func Log(s string) (*os.File, error) {
	src, err := os.Stat(s)
	if err != nil {
		if err = os.Mkdir(s, os.ModePerm); err != nil {
			//log.Fatal(err.Error())
			return nil, err
		}
	} else {
		if !src.IsDir() {
			// log.Fatal("Logs is not a directory")
			// os.Exit(1)
			return nil, errors.New(s + " is not a directory")
		}
	}

	f, err := os.OpenFile("./"+s+"/"+PreText("logs")+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		// fmt.Println(err.Error())
		// log.Fatalf("error opening file: %v", err)
		return nil, err
	}
	//defer f.Close()

	// wrt := io.MultiWriter(os.Stdout, f)
	wrt := io.MultiWriter(f)

	log.SetOutput(wrt)

	return f, nil

}

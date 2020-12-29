package models_test

import (
	"netcat/internal/app/models"
	testhistory "netcat/internal/app/store/testHistory"
	"netcat/internal/app/tcpserver"
	"testing"
	"time"
)

func TestConnection_Prefix(t *testing.T) {

	config := tcpserver.TestConfig(t)
	server, err := tcpserver.TestServer(t, config)
	if err != nil {
		t.Fatal(err)
	}
	// s := testhistory.New(t)
	// s.Pool().Create(1, 1)

	go func() {
		_ = server.Run()
	}()

	time.Sleep(time.Millisecond * 100)

	conn, err := models.TestConnection(t, "John", config.Port, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	conn2, err := models.TestConnection(t, "Jane", config.Port, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	conn.Test = true

	err = conn.Prefix()
	if err != nil {
		t.Log("not print to testing connection")
		t.Fail()
	}
	conn.Test = false
	(*conn.Conn).Close()
	err = conn.Prefix()
	if err == nil {
		t.Log("conn closed")
		t.Fail()
	}

	conn.Conn = nil

	err = conn.Prefix()
	if err == nil {
		t.Log("conn not init")
		t.Fail()
	}

	err = conn2.Prefix()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestConnection_Write(t *testing.T) {

	config := tcpserver.TestConfig(t)
	server, err := tcpserver.TestServer(t, config)
	if err != nil {
		t.Fatal(err)
	}
	s := testhistory.New(t)
	s.Pool().Create(1, 1)

	go func() {
		_ = server.Run()
	}()

	time.Sleep(time.Millisecond * 100)

	conn, err := models.TestConnection(t, "John", config.Port, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	conn2, err := models.TestConnection(t, "Jane", config.Port, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	conn.Test = true
	err = conn.Write("")
	if err != nil {
		t.Log("not print to testing connection")
		t.Fail()
	}

	conn.Test = false
	(*conn.Conn).Close()
	err = conn.Write("")
	if err == nil {
		t.Log("conn closed")
		t.Fail()
	}

	conn.Conn = nil

	err = conn.Write("")
	if err == nil {
		t.Log("conn not init")
		t.Fail()
	}

	err = conn2.Write("")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

package pool_test

import (
	"kitKat/internal/app/models"
	testhistory "kitKat/internal/app/store/testHistory"
	"kitKat/internal/app/tcpserver"
	"testing"
)

func TestRepo_SendMsg(t *testing.T) {

	config := tcpserver.TestConfig(t)

	s := testhistory.New(t)

	if err := s.History().Create("test", config.Rooms); err != nil {

		t.Fatal(err)
	}

	pool, err := s.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	var m *models.Message

	if err := s.Pool().SendMsg(m); err == nil {
		t.Fatal("not send if m eq nil")
	}

	m = &models.Message{
		All:  true,
		Text: "test",
	}

	if err := s.Pool().SendMsg(m); err != nil {
		t.Fatal(err)
	}

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	c, err := models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}
	conn1.Conn = &c

	err = s.Pool().Add(conn1)
	if err != nil {
		t.Fatal(err)
	}

	//////

	if err = s.Pool().SendMsg(m); err == nil {
		t.Fatal("second Pre msg can not be send because conn closed after first msg")
	}

	if err = s.Pool().SendMsg(m); err == nil {
		t.Fatal("send text to closed conn")
	}

	delete(pool, conn1.Conn)

	if len(pool) != 0 {
		t.Fatal("pool must be empty")
	}

	///

	m = &models.Message{
		Text: "test",
	}

	if err = s.Pool().SendMsg(m); err == nil {
		t.Fatal("room of mess not init")
	}

	m.Room = 1

	if err = s.Pool().SendMsg(m); err != nil {
		t.Fatal(err)
	}

	conn1, err = models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	c, err = models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}
	conn1.Conn = &c

	err = s.Pool().Add(conn1)
	if err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().SendMsg(m); err == nil {
		t.Fatal("second Pre msg can not be send because conn closed after first msg")
	}

	if err = s.Pool().SendMsg(m); err == nil {
		t.Fatal("send text to closed conn")
	}

	delete(pool, conn1.Conn)

	if len(pool) != 0 {
		t.Fatal("pool must be empty")
	}

}

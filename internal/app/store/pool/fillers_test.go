package pool_test

import (
	"kitKat/internal/app/models"
	testhistory "kitKat/internal/app/store/testHistory"
	"kitKat/internal/app/tcpserver"
	"net"
	"reflect"
	"testing"
)

func TestRepo_Add(t *testing.T) {

	config := tcpserver.TestConfig(t)

	s := testhistory.New(t)

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}
	conn2, err := models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}
	conn3, err := models.NewTestConn(t, "Jakob", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}
	conn4, err := models.NewTestConn(t, "Jack", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}
	if err = (*conn4.Conn).Close(); err != nil {
		t.Fatal(err)
	}

	err = s.Pool().Add(conn1)
	if err == nil {
		t.Log("try to add if capacity = 0")
		t.Fail()
	}

	s.Pool().Create(1, config.Rooms)

	err = s.Pool().Add(&models.Connection{})
	if err == nil {
		t.Log("try to add if conn inside connection eq nil")
		t.Fail()
	}

	err = s.Pool().Add(conn1)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = s.Pool().Add(conn2)
	if err == nil {
		t.Log("try to add if pool is full")
		t.Fail()
	}

	err = s.Pool().Add(conn4)
	if err == nil {
		t.Log("try to add if pool is full and connection closed")
		t.Fail()
	}

	s.Pool().Create(3, config.Rooms)

	(*conn1.Conn).Close()

	err = s.Pool().Add(conn3)
	if err == nil {
		t.Log("err if one off connections closed")
		t.Fail()
	}

	//testing auto closing connection after one string incoming

	s = testhistory.New(t)

	s.Pool().Create(10, config.Rooms)

	conn1, err = models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err = models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().Add(conn1); err != nil {
		t.Fatal(err)
	}

	limitedConn, err := models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}
	conn1.Conn = &limitedConn

	// var wg sync.WaitGroup

	// wg.Add(1)

	// go func() {

	// 	buff := make([]byte, 100)

	// 	wg.Done()

	// 	time.Sleep(time.Millisecond * 80)

	// 	num, err := (*conn1.Conn).Read(buff)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	if num == 0 {
	// 		t.Fatal("num must be not eq zero")
	// 	}
	// 	if err = (*conn1.Conn).Close(); err != nil {
	// 		t.Fatal(err)
	// 	}

	// }()

	// wg.Wait()

	if err = s.Pool().Add(conn2); err == nil {
		t.Fatal("second Pre msg can not be send because conn closed after first msg")
	}

}

func TestRepo_Get(t *testing.T) {

	config := tcpserver.TestConfig(t)

	s := testhistory.New(t)

	s.Pool().Create(2, 2)

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Pool().Get(nil)
	if err == nil {
		t.Fatal("not init to input")
	}

	c, err := s.Pool().Get(conn1.Conn)
	if err == nil {
		t.Fatal("not added before")
	}
	if c != nil {
		t.Fatal("nothing to return")
	}

	err = s.Pool().Add(conn1)
	if err != nil {
		t.Fatal(err)
	}

	c, err = s.Pool().Get(conn1.Conn)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(conn1, c) {
		t.Fatal("must be equal")
	}

}

func TestRepo_DeleteByConn(t *testing.T) {

	config := tcpserver.TestConfig(t)
	s := testhistory.New(t)
	s.Pool().Create(2, 2)

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Pool().DeleteByConn(nil)
	if err == nil {
		t.Fatal("not init to input")
	}

	if err = s.Pool().Add(conn1); err != nil {
		t.Fatal(err)
	}

	err = s.Pool().DeleteByConn(conn1.Conn)
	if err != nil {
		t.Fatal(err)
	}

	c, err := s.Pool().Get(conn1.Conn)
	if c != nil {
		t.Fatal("nothing to get after deleting")
	}

	if err = s.Pool().Add(conn2); err != nil {
		t.Fatal(err)
	}

	if err = (*conn2.Conn).Close(); err != nil {
		t.Fatal(err)
	}

	err = s.Pool().DeleteByConn(conn2.Conn)
	if err == nil {
		t.Fatal("can not to close closed conn")
	}

}

func TestRepo_Removeuser(t *testing.T) {

	var (
		m   *models.Message
		err error
		c   net.Conn
	)

	config := tcpserver.TestConfig(t)

	s := testhistory.New(t)

	pool, err := s.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("invalid input")
	}

	m = &models.Message{}

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("conn not init")
	}

	c = models.TestNetConn(t)

	m.Conn = &c

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("kick or out not init")
	}

	if err = c.Close(); err != nil {
		t.Fatal(err)
	}

	m.Kick = true

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("send kick mess to closed conn")
	}

	m.Out = true
	m.Kick = false

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("delete closed conn")
	}

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err := models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().Add(conn1); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().Add(conn2); err != nil {
		t.Fatal(err)
	}

	m.Kick = true
	m.Out = false

	if err = (*conn1.Conn).Close(); err != nil {
		t.Fatal(err)
	}

	m.Conn = conn2.Conn
	m.Room = conn2.Room

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("send mess about deleting to closed conn")
	}

	delete(pool, conn1.Conn)

	if len(pool) != 0 {
		t.Fatal("pool must be empty")
	}

	conn1, err = models.NewTestConn(t, "Jill", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err = models.NewTestConn(t, "JacK", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().Add(conn1); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().Add(conn2); err != nil {
		t.Fatal(err)
	}

	c, err = models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}

	pool[&c] = pool[conn1.Conn]

	delete(pool, conn1.Conn)

	conn1.Conn = &c

	m.Conn = conn2.Conn
	m.Author = conn2.Name
	m.Room = conn2.Room

	if err = s.Pool().RemoveUser(m); err == nil {
		t.Fatal("prefix must be failed cuz conn 1 closed after first message")
	}

	c = models.TestNetConn(t)
	conn1.Conn = &c

	m.Conn = conn1.Conn
	m.Author = conn1.Name

	if err = s.Pool().RemoveUser(m); err != nil {
		t.Fatal(err)
	}

	if len(pool) != 0 {
		t.Fatal("pool must be empty")
	}

}

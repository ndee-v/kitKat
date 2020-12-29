package pool_test

import (
	"kitKat/internal/app/models"
	testhistory "kitKat/internal/app/store/testHistory"
	"kitKat/internal/app/tcpserver"
	"testing"
)

func TestRepo_Change(t *testing.T) {

	s := testhistory.New(t)

	var (
		mes *models.Message
	)

	err := s.Pool().Change(mes)
	if err == nil {
		t.Log("mes not init")
		t.Fail()
	}

	mes = &models.Message{
		Text: "hello",
	}

	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !mes.Help {
		t.Log("help bool must be true")
		t.Fail()
	}

	//change name

	mes = &models.Message{
		Text: "--name ",
	}
	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if mes.ChangeName != "" || !mes.Help {
		t.Log("not to change and help bool must be true")
		t.Fail()
	}

	mes = &models.Message{
		Text: "--name John",
	}
	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if mes.ChangeName != "John" || mes.Help {
		t.Log("changed name and help false")
		t.Fail()
	}

	//change room

	mes = &models.Message{
		Text: "--room",
	}
	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !mes.Help {
		t.Log("help bool must be true if rooms not init")
		t.Fail()
	}

	mes = &models.Message{
		Text: "--room 9223372036854775808",
	}
	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !mes.Help {
		t.Log("help bool must be true if rooms overflow")
		t.Fail()
	}

	mes = &models.Message{
		Text: "--room 9223372036854775807",
		Room: 25,
	}
	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if mes.Help || !mes.Move || mes.RoomFrom != 25 || mes.Room != 9223372036854775807 {
		t.Log("must be valid input")
		t.Fail()
	}

	// check online
	mes = &models.Message{
		Text: "--online",
	}
	err = s.Pool().Change(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if mes.Help || !mes.CheckOnline {
		t.Log("hepl must be false && chekonline must be true")
		t.Fail()
	}

}

func TestRepo_RenameUser(t *testing.T) {

	config := tcpserver.TestConfig(t)

	s := testhistory.New(t)

	pool, err := s.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

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

	err = s.Pool().Add(conn1)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = s.Pool().Add(conn2)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	var mes *models.Message

	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("not init input")
		t.Fail()
	}

	mes = &models.Message{
		Conn: conn2.Conn,
	}

	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("not init change name")
		t.Fail()
	}

	mes.Conn = nil
	mes.ChangeName = "John Doe"
	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("not init conn")
		t.Fail()
	}

	mes.Conn = conn2.Conn
	err = s.Pool().RenameUser(mes)
	if err != nil {
		t.Log("invalid new name")
		t.Fail()
	}

	(*conn2.Conn).Close()
	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("invalid name and write to closed conn")
		t.Fail()
	}

	mes.ChangeName = "John"
	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("used name and write to closed conn")
		t.Fail()
	}

	mes.ChangeName = "John"
	mes.Conn = conn3.Conn

	err = s.Pool().RenameUser(mes)
	if err != nil {
		t.Log("used name and write to open conn")
		t.Fail()
	}

	mes.ChangeName = "John_Doe"
	mes.Conn = conn3.Conn
	mes.Author = conn3.Name
	mes.Room = conn3.Room

	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("valid name send mess to others if one closed")
		t.Fail()
	}

	mes.ChangeName = "Jane_Doe"
	mes.Conn = conn2.Conn
	mes.Author = conn2.Name
	mes.Room = conn2.Room

	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Log("valid name send by closed conn")
		t.Fail()
	}

	delete(pool, conn2.Conn)

	mes.ChangeName = "Jane_Doe"
	mes.Conn = conn1.Conn
	mes.Author = conn1.Name
	mes.Room = conn1.Room

	err = s.Pool().RenameUser(mes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	/////

	s = testhistory.New(t)

	pool, err = s.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	conn1, err = models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}
	conn2, err = models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	mes.ChangeName = "Jane_Do"
	mes.Conn = conn1.Conn
	mes.Author = conn1.Name
	mes.Room = conn1.Room

	err = s.Pool().Add(conn1)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Pool().Add(conn2)
	if err != nil {
		t.Fatal(err)
	}

	nC, err := models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}

	pool[&nC] = pool[conn2.Conn]

	delete(pool, conn2.Conn)

	conn2.Conn = &nC

	err = s.Pool().RenameUser(mes)
	if err == nil {
		t.Fatal("conn2 must be closed before Prefix starts")
	}
}

func TestRepo_WhoIsOnline(t *testing.T) {

	var (
		msg *models.Message
		err error
	)

	s := testhistory.New(t)

	config := tcpserver.TestConfig(t)

	_, err = s.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if err = s.Pool().WhoIsOnline(msg); err == nil {
		t.Fatal("invalid arguments")
	}

	msg = &models.Message{
		Text: "Testing",
	}

	if err = s.Pool().WhoIsOnline(msg); err == nil {
		t.Fatal("invalid arguments")
	}

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	msg.Conn = conn1.Conn

	if err = s.Pool().WhoIsOnline(msg); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().Add(conn1); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().WhoIsOnline(msg); err != nil {
		t.Fatal(err)
	}

	if err = (*conn1.Conn).Close(); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().WhoIsOnline(msg); err == nil {
		t.Fatal("can not sending to closed conn")
	}

	//close after first input

	nC, err := models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}

	conn1.Conn = &nC
	msg.Conn = conn1.Conn

	if err = s.Pool().WhoIsOnline(msg); err == nil {
		t.Fatal("conn1 must be closed after first input")
	}

	//close after second input

	nC, err = models.TestNetConnWriteLimited(t, 2)
	if err != nil {

		t.Fatal(err)
	}
	conn1.Conn = &nC

	msg.Conn = conn1.Conn

	if err = s.Pool().WhoIsOnline(msg); err == nil {
		t.Fatal("conn1 must be closed after second input")
	}

}

func TestRepo_ChangeRoom(t *testing.T) {

	var (
		msg *models.Message
		err error
	)

	s := testhistory.New(t)

	config := tcpserver.TestConfig(t)

	conn1, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := s.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("nothing input")
	}

	msg = &models.Message{
		Text: "hello",
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("move boolean must be true")
	}

	msg.Move = true

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("mess conn must be init")
	}

	msg.Conn = conn1.Conn
	msg.RoomFrom = conn1.Room
	msg.Room = conn1.Room

	if err = s.Pool().ChangeRoom(msg); err != nil {
		t.Fatal(err)
	}

	if err = (*conn1.Conn).Close(); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("conn1 must be closed")
	}

	///new conn initialize
	nc := models.TestNetConn(t)
	conn1.Conn = &nc

	msg.Conn = conn1.Conn
	msg.Room = 0

	if err = s.Pool().ChangeRoom(msg); err != nil {
		t.Fatal(err)
	}

	msg.Room = config.Rooms + 1

	if err = (*conn1.Conn).Close(); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("conn1 must be closed")
	}

	///

	msg.Room = 2
	msg.RoomFrom = 1

	if err = s.Pool().Add(conn1); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("conn1 must be closed")
	}

	// nC, err := models.TestNetConnWriteLimited(t, 1)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	nC := models.TestNetConn(t)

	pool[&nC] = pool[conn1.Conn]

	delete(pool, conn1.Conn)

	conn1.Conn = &nC

	msg.Conn = &nC

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("history can not writes")
	}

	if err = s.History().Create("test", config.Rooms); err != nil {
		t.Fatal(err)
	}

	if err = s.Pool().ChangeRoom(msg); err != nil {
		t.Fatal(err)
	}

	delete(pool, conn1.Conn)

	if len(pool) != 0 {
		t.Fatal("pool must be empty")
	}

	///
	conn1, err = models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err := models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	nC, err = models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}
	conn2.Conn = &nC

	pool[conn1.Conn] = conn1
	pool[conn2.Conn] = conn2

	msg.Conn = conn1.Conn

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("can not write second string to conn from roomFrom")
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("can not write string to closed conn from roomFrom")
	}

	///
	nC, err = models.TestNetConnWriteLimited(t, 1)
	if err != nil {
		t.Fatal(err)
	}

	delete(pool, conn2.Conn)

	conn2.Conn = &nC

	conn2.Room = 2

	pool[&nC] = conn2

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("can not write second string to conn from newRoom")
	}

	if err = s.Pool().ChangeRoom(msg); err == nil {
		t.Fatal("can not write string to closed conn from newRoom")
	}

}

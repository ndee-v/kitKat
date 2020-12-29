package pool_test

import (
	"kitKat/internal/app/models"
	testhistory "kitKat/internal/app/store/testHistory"
	"kitKat/internal/app/tcpserver"
	"testing"
	"time"
)

func TestRepo_IdleTimer(t *testing.T) {

	config := tcpserver.TestConfig(t)

	conn, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	store := testhistory.New(t)

	_, err = store.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = store.Pool().Add(conn)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	worker := true
	ti := time.Duration(time.Millisecond * 20)
	ch := make(chan *models.Message, 2)

	close(ch)
	go func() {
		err = store.Pool().IdleTimer(&worker, &ti, &ch)
		if err == nil {
			t.Log("chan must be closed")
			t.Fail()
		}
	}()

	time.Sleep(time.Millisecond * 30)
	ch = make(chan *models.Message, 2)
	go func() {
		err = store.Pool().IdleTimer(&worker, &ti, &ch)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}()

	time.Sleep(time.Millisecond * 100)
	worker = false
}

func TestRepo_RegTimer(t *testing.T) {
	assert := Assert(t, true)
	config := tcpserver.TestConfig(t)

	conn, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	store := testhistory.New(t)

	_, err = store.Pool().Create(config.Conns, config.Rooms)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	ti := time.Duration(time.Millisecond * 20)

	err = store.Pool().RegTimer(conn.Conn, &ti)
	assert(err == nil, "conn must be opened")

	err = store.Pool().RegTimer(conn.Conn, &ti)
	assert(err != nil, "conn must be closed")

	err = store.Pool().Add(conn)
	assert(err == nil, "conn must be added")

	err = store.Pool().RegTimer(conn.Conn, &ti)
	assert(err == nil, "conn must be into pool")

}

func Assert(t *testing.T, bb bool) func(bool, ...interface{}) {
	if bb {
		return func(b bool, message ...interface{}) {
			if !b {
				t.Fatal(message)
			}
		}
	} else {
		return func(b bool, message ...interface{}) {
			if !b {
				t.Log(message)
			}
		}
	}

}

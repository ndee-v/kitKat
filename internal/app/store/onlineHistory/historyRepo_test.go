package onlinehistory_test

import (
	"kitKat/internal/app/models"
	onlinehistory "kitKat/internal/app/store/onlineHistory"
	"kitKat/internal/app/tcpserver"
	"log"
	"os"
	"testing"
)

func TestHistoryRepo_Create(t *testing.T) {

	type args struct {
		dir string
		num int
	}

	if err := os.Mkdir("closed_Dir", os.ModeDir); err != nil {
		log.Fatal(err)
	}

	_, err := os.Create("notDir")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll("closed_Dir"); err != nil {
			log.Fatal(err)
		}
		if err := os.RemoveAll("notDir"); err != nil {
			log.Fatal(err)
		}

	}()

	//store := New()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				"tempHisto",
				3,
			},
			wantErr: false,
		},
		{
			name: "invalid input#1",
			args: args{
				"",
				3,
			},
			wantErr: true,
		},
		{
			name: "invalid input#2",
			args: args{
				"longtext___",
				3,
			},
			wantErr: true,
		},
		{
			name: "invalid input#3",
			args: args{
				"tempDir",
				0,
			},
			wantErr: true,
		},
		{
			name: "invalid input#4",
			args: args{
				"tempDir",
				-1,
			},
			wantErr: true,
		},
		{
			name: "invalid input#5",
			args: args{
				"temp dir",
				4,
			},
			wantErr: true,
		},
		{
			name: "closed directory",
			args: args{
				"closed_Dir",
				4,
			},
			wantErr: true,
		},
		{
			name: "not a dir",
			args: args{
				"notDir",
				4,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := onlinehistory.New()
			if err := h.History().Create(tt.args.dir, tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("HistoryRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := os.RemoveAll(tt.args.dir); err != nil {
				log.Fatal(err)
			}
		})
	}
}

func TestHistoryRepo_AddInto(t *testing.T) {

	h := onlinehistory.New().History()

	m := &models.Message{
		Text: "Hello World",
		Room: 1,
	}

	if err := h.AddInto(m); err == nil {
		t.Log("add to file if it is not exist")
		t.Fail()
	}

	if err := h.Create("tempDir", 2); err != nil {
		t.Log(err)
		t.Fail()
	}
	defer func() {
		if err := os.RemoveAll("tempDir"); err != nil {
			t.Log(err)
			t.Fail()
		}
	}()

	if err := h.AddInto(m); err != nil {
		t.Log(err)
		t.Fail()
	}

	m.Room = 3
	if err := h.AddInto(m); err == nil {
		t.Log("num of room invalid")
		t.Fail()
	}
}

func TestHistoryRepo_PrintToConn(t *testing.T) {

	config := tcpserver.TestConfig(t)

	s := onlinehistory.New()

	if err := s.History().Create("tempDir", config.Rooms); err != nil {
		t.Log(err)
		t.Fail()
	}
	defer func() {
		if err := os.RemoveAll("tempDir"); err != nil {
			t.Log(err)
			t.Fail()
		}
	}()

	conn, err := models.NewTestConn(t, "John", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err := models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn2.Conn = nil

	conn3, err := models.NewTestConn(t, "Jane", config.Rooms)
	if err != nil {
		t.Fatal(err)
	}

	conn3.Room = config.Rooms + 1

	err = s.History().PrintToConn(conn)
	if err != nil {
		t.Fatal(err)
	}

	m1 := &models.Message{
		Text: "Test string1",
		Room: 1,
	}

	m3 := &models.Message{
		Text: "Test string3",
		Room: 1,
	}

	err = s.History().AddInto(m1)
	if err != nil {
		t.Fatal(err)
	}

	err = s.History().PrintToConn(conn)
	if err != nil {
		t.Fatal(err)
	}

	err = s.History().PrintToConn(nil)
	if err == nil {
		t.Fatal("conn is not init")
	}

	err = s.History().PrintToConn(conn2)
	if err == nil {
		t.Fatal("conn.Conn is not init")
	}

	err = s.History().AddInto(m3)
	if err != nil {
		t.Fatal(err)
	}
	(*conn.Conn).Close()
	err = s.History().PrintToConn(conn)
	if err == nil {
		t.Fatal("prints to closed conn")
	}

	err = s.History().PrintToConn(conn3)
	if err == nil {
		t.Fatal("invalid room")
	}
}

package models_test

import (
	"netcat/internal/app/models"
	testhistory "netcat/internal/app/store/testHistory"
	"netcat/internal/app/tcpserver"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	type args struct {
		te   *testing.T
		port string
	}

	config := tcpserver.TestConfig(t)
	server, err := tcpserver.TestServer(t, config)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		_ = server.Run()
	}()

	time.Sleep(time.Millisecond * 100)

	tests := []struct {
		name string
		args args
		// want    *net.Conn
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				te:   t,
				port: config.Port,
			},
			// want:    &net.Conn{},
			wantErr: false,
		},
		{
			name: "not normal 1",
			args: args{
				te:   t,
				port: "",
			},
			// want:    &net.Conn{},
			wantErr: true,
		},
		{
			name: "not normal 2",
			args: args{
				te:   t,
				port: "hello",
			},
			// want:    &net.Conn{},
			wantErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := models.TestConn(tt.args.te, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("TestConn() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestConnection(t *testing.T) {
	type args struct {
		te    *testing.T
		str   string
		port  string
		rooms int
	}

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

	tests := []struct {
		name string
		args args
		// want    *models.Connection
		wantErr bool
	}{
		{
			name: "valid//need to tcp server active",
			args: args{
				te:    t,
				str:   "John",
				port:  config.Port,
				rooms: config.Rooms,
			},
			//
			wantErr: false,
		},
		{
			name: "invalid 1",
			args: args{
				te:    t,
				str:   "",
				port:  config.Port,
				rooms: config.Rooms,
			},
			//
			wantErr: true,
		},
		{
			name: "invalid 2",
			args: args{
				te:    t,
				str:   "Jane",
				port:  "",
				rooms: config.Rooms,
			},
			//
			wantErr: true,
		},
		{
			name: "invalid 3",
			args: args{
				te:    t,
				str:   "Jane",
				port:  config.Port,
				rooms: 0,
			},
			//
			wantErr: true,
		},
		{
			name: "invalid 4",
			args: args{
				te:    t,
				str:   "Jane Doe",
				port:  config.Port,
				rooms: 3,
			},
			//
			wantErr: true,
		},
		{
			name: "invalid 5",
			args: args{
				te:    t,
				str:   "Jane_Doe",
				port:  ":0",
				rooms: config.Rooms,
			},
			//
			wantErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := models.TestConnection(tt.args.te, tt.args.str, tt.args.port, tt.args.rooms)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("TestConnection() = %v, want %v", got, tt.want)
			// }
		})
	}
}

package models_test

import (
	"io"
	"net"
	"netcat/internal/app/models"
	"reflect"
	"testing"
)

func TestNetConnection_Write(t *testing.T) {

	c := models.NetConnection{
		Open: true,
	}

	str := "Hello"

	num, err := c.Write([]byte(str))
	if err != nil {
		t.Fatal(err)
	}
	if num != len([]byte(str)) {
		t.Fatal("num must be eq len of string")
	}

	if err := c.Close(); err != nil {
		t.Fatal(err)
	}

	num, err = c.Write([]byte(str))
	if err == nil {
		t.Fatal("connection must be closed")
	}

	if num != 0 {
		t.Fatal("0 must be returns")
	}

}

func TestNetConnection_Close(t *testing.T) {

	//c := models.TestNetConn(t)

	c := models.NetConnection{
		Open: true,
	}

	if err := c.Close(); err != nil {
		t.Fatal(err)
	}

	if err := c.Close(); err == nil {
		t.Fatal("conn must be closed")
	}
}

func TestNetConnection_Read(t *testing.T) {

	// c := models.TestNetConn(t)
	c := &models.NetConnection{
		Open: true,
	}

	buf2 := make([]byte, 10)

	t.Logf("Test: %p\n", buf2)
	var num int
	var err error

	num, err = c.Read(buf2)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Buff) != 0 {
		t.Fatal("must be empty")
	}

	if num != 0 {
		t.Fatal("0 must be returned from empty buffer")
	}

	str := "Testing"
	_, err = c.Write([]byte(str))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Buff) == 0 {
		t.Fatal("must be not empty")
	}
	t.Logf("Test2: %p", buf2)

	num, err = c.Read(buf2)
	if err != nil {
		t.Fatal(err)
	}

	if num != len(str) {
		t.Fatal("all bytes must be returned")
	}

	buf2 = make([]byte, 2)
	_, err = c.Write([]byte(str))
	if err != nil {
		t.Fatal(err)
	}

	num, err = c.Read(buf2)
	if err != io.ErrUnexpectedEOF {
		t.Fatal("buffer must be so short for test")
	}
	if num != len(buf2) {
		t.Fatal("num of returns must be eq len of short buffer")
	}

	buf2 = make([]byte, 20)

	if err = c.Close(); err != nil {
		t.Fatal(err)
	}

	num, err = c.Read(buf2)
	if err == nil {
		t.Fatal("conn must be closed")
	}
	if num != 0 {
		t.Fatal("nothing must be return")
	}

}

func TestNetConnWriteLimited(t *testing.T) {
	type args struct {
		t   *testing.T
		num int
	}
	tests := []struct {
		name    string
		args    args
		want    net.Conn
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				t:   t,
				num: 1,
			},
			want: &models.NetConnection{
				Open:       true,
				WriteLimit: true,
				WritesLeft: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid 1",
			args: args{
				t:   t,
				num: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid 2",
			args: args{
				t:   t,
				num: 0,
			},
			want:    nil,
			wantErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := models.TestNetConnWriteLimited(tt.args.t, tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestNetConnWriteLimited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestNetConnWriteLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

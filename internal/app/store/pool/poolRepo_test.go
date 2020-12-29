package pool_test

import (
	"net"
	"netcat/internal/app/models"
	testhistory "netcat/internal/app/store/testHistory"
	"netcat/internal/app/tcpserver"
	"reflect"
	"testing"
)

// func TestRepo_Create(t *testing.T) {
// 	s := testhistory.New(t)
// 	conf := tcpserver.TestConfig(t)

// }

func TestRepo_Create(t *testing.T) {

	s := testhistory.New(t)
	conf := tcpserver.TestConfig(t)

	type args struct {
		cap   int
		rooms int
	}
	tests := []struct {
		name string
		//fields  fields
		args    args
		want    map[*net.Conn]*models.Connection
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				cap:   conf.Conns,
				rooms: conf.Rooms,
			},
			want:    make(map[*net.Conn]*models.Connection),
			wantErr: false,
		},
		{
			name: "not normal 1",
			args: args{
				cap:   0,
				rooms: conf.Rooms,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not normal 2",
			args: args{
				cap:   conf.Conns,
				rooms: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not normal 3",
			args: args{
				cap:   -1,
				rooms: conf.Rooms,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not normal 4",
			args: args{
				cap:   conf.Conns,
				rooms: -1,
			},
			want:    nil,
			wantErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Pool().Create(tt.args.cap, tt.args.rooms)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

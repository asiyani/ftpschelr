package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/asiyani/ftpschelr/ftpschelr"
	scribble "github.com/nanobox-io/golang-scribble"
)

var (
	now = time.Now().UTC()
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *DB
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDB_Store(t *testing.T) {
	type fields struct {
		Driver *scribble.Driver
	}
	type args struct {
		f ftpschelr.Connection
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{ftpschelr.Connection{ID: "1dbide1", Name: "FTP_TEST1", SerAddr: "my.ftp.serv.net:21", User: "anonymous", Pass: "pass12$*"}},
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{ftpschelr.Connection{ID: "1dbide2", Name: "FTP_TEST2", SerAddr: "speedtest.tele2.net:21", User: "anonymous", Pass: ""}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.f.CreateSchedules(".", "./ftpfiles", "10MB.zip", ftpschelr.Download, now, (10 * time.Second))
			d, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if err := d.Store(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("DB.Store() error = %+v, wantErr %+v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Restore(t *testing.T) {
	type fields struct {
		Driver *scribble.Driver
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ftpschelr.Connection
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{ID: "1dbide1"},
			wantErr: false,
			want:    ftpschelr.Connection{ID: "1dbide1", Name: "FTP_TEST1", SerAddr: "my.ftp.serv.net:21", User: "anonymous", Pass: "pass12$*"},
		},
		{
			name:    "test2",
			args:    args{ID: "1dbide2"},
			wantErr: false,
			want:    ftpschelr.Connection{ID: "1dbide2", Name: "FTP_TEST2", SerAddr: "speedtest.tele2.net:21", User: "anonymous", Pass: ""},
		},
		{
			name:    "test3",
			args:    args{ID: "1dbide3"},
			wantErr: true,
			want:    ftpschelr.Connection{},
		},
	}
	tests[0].want.CreateSchedules(".", "./ftpfiles", "10MB.zip", ftpschelr.Download, now, (10 * time.Second))
	tests[1].want.CreateSchedules(".", "./ftpfiles", "10MB.zip", ftpschelr.Download, now, (10 * time.Second))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New()
			if err != nil {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			got, err := d.Restore(tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Restore() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.Restore() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDB_RestoreAll(t *testing.T) {
	type fields struct {
		Driver *scribble.Driver
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ftpschelr.Connection
		wantErr bool
	}{
		{
			name: "test1",
			want: []ftpschelr.Connection{
				{ID: "1dbide1", Name: "FTP_TEST1", SerAddr: "my.ftp.serv.net:21", User: "anonymous", Pass: "pass12$*"},
				{ID: "1dbide2", Name: "FTP_TEST2", SerAddr: "speedtest.tele2.net:21", User: "anonymous", Pass: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want[0].CreateSchedules(".", "./ftpfiles", "10MB.zip", ftpschelr.Download, now, (10 * time.Second))
			tt.want[1].CreateSchedules(".", "./ftpfiles", "10MB.zip", ftpschelr.Download, now, (10 * time.Second))
			d, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			got, err := d.RestoreAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.RestoreAll() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.RestoreAll() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDB_Remove(t *testing.T) {
	type fields struct {
		Driver *scribble.Driver
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{ID: "1dbide1"},
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{ID: "1dbide3"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New()
			if err != nil {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if err := d.Remove(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("DB.Remove() error = %+v, wantErr %+v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_RemoveAll(t *testing.T) {
	type fields struct {
		Driver *scribble.Driver
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New()
			if err != nil {
				t.Errorf("New() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if err := d.RemoveAll(); (err != nil) != tt.wantErr {
				t.Errorf("DB.RemoveAll() error = %+v, wantErr %+v", err, tt.wantErr)
			}
		})
	}
}

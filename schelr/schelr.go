package schelr

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
)

type stream int

const (
	//Download and upload stream type
	Download stream = iota
	Upload
)

type job struct {
	FtpDir    string
	LocalDir  string
	FileName  string
	Direction stream
	StartAt   time.Time
	Interval  time.Duration
	NextRun   time.Time
	PastRuns  []time.Time
	Ticker    *time.Timer
}

// Connection is data type for ftp scheduler
type Connection struct {
	ID      string
	Name    string
	SerAddr string
	User    string
	Pass    string
	Jobs    []job
}

// Connector is interface for Connection Type
type Connector interface {
	ConnAndLogin() (*ftp.ServerConn, error)
	CreateSchedules(fDir, lDir, fName string, d stream, t time.Time, inv time.Duration)
	ScheduleJob(index int)
	CancelJobs(index int)
}

// NewConnection will create and return new Connection with random ID.
func NewConnection(name, serAddr, user, pass string) *Connection {
	return &Connection{ID: strconv.FormatInt(time.Now().Unix(), 32),
		Name: name, SerAddr: serAddr, User: user, Pass: pass}

}

// ConnAndLogin will Connect and login to ftp server.
func (f *Connection) ConnAndLogin() (*ftp.ServerConn, error) {
	srvCon, err := ftp.Dial(f.SerAddr)
	if err != nil {
		return nil, err
	}

	err = srvCon.Login(f.User, f.Pass)
	if err != nil {
		return nil, err
	}

	return srvCon, nil
}

// CreateSchedules creates new Jobs and adds to Connection.
func (f *Connection) CreateSchedules(fDir, lDir, fName string, d stream, t time.Time, intr time.Duration) {
	s := job{
		FtpDir:    fDir,
		LocalDir:  lDir,
		FileName:  fName,
		Direction: d,
		StartAt:   t,
		Interval:  intr,
	}

	s.NextRun = s.StartAt.Add(s.Interval)

	f.Jobs = append(f.Jobs, s)
}

// ScheduleJob will schedule Jobs[index].
func (f *Connection) ScheduleJob(index int) {

	var starDur time.Duration

	//If no future NextRun then then exit.
	if f.Jobs[index].NextRun.Sub(time.Now()) < 0 {
		return
	}

	if f.Jobs[index].StartAt.Sub(time.Now()) < 0 {
		starDur = f.Jobs[index].NextRun.Sub(time.Now())
	} else {
		starDur = f.Jobs[index].StartAt.Sub(time.Now())
	}

	f.Jobs[index].Ticker = time.AfterFunc(starDur, func() {
		if f.Jobs[index].Direction == Download {
			downloader(f, f.Jobs[index])
		} else {
			uploader(f, f.Jobs[index])
		}
		f.Jobs[index].PastRuns = append(f.Jobs[index].PastRuns, f.Jobs[index].NextRun)
		updateSchedule(&f.Jobs[index])
		f.ScheduleJob(index)
	})

}

// CancelJobs already scheduled and future Jobs.
func (f *Connection) CancelJobs(index int) {
	f.Jobs[index].Ticker.Stop()
	f.Jobs[index].Interval = (0 * time.Second)
}

func updateSchedule(s *job) {
	if s.Interval != 0 {
		s.NextRun = s.NextRun.Add(s.Interval)
	}
}

// GetList will get list of all files in path directory.
func getList(f Connector, path string) ([]*ftp.Entry, error) {

	srvCon, err := f.ConnAndLogin()
	if err != nil {
		return nil, fmt.Errorf("error connecting ftp server: %v", err)
	}

	entries, err := srvCon.List(path)
	if err != nil {
		return nil, fmt.Errorf("error accessing dir: %v", err)
	}

	return entries, nil
}

// Downloader will download file from ftp server (FtpDir) and store in to local drive (LocalDir).
func downloader(f Connector, s job) error {

	log.Printf("%s\n", "Downloading a file.."+s.FileName)
	return nil

	srvCon, err := f.ConnAndLogin()
	if err != nil {
		return fmt.Errorf("error connecting ftp server: %v", err)
	}

	resp, err := srvCon.Retr(s.FtpDir + "/" + s.FileName)
	if err != nil {
		return fmt.Errorf("error Retrieving  file %v", err)
	}
	defer resp.Close()

	out, err := os.Create(s.LocalDir + "/" + s.FileName)
	if err != nil {
		return fmt.Errorf("error creating file %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp); err != nil {
		return fmt.Errorf("copying data from conn to file %v", err)
	}

	return nil
}

// Uploader will Upload file from local drive (LocalDir) to ftp server dir (FtpDir).
func uploader(f Connector, s job) error {

	log.Printf("%s\n", "Uploading a file.."+s.FileName)
	return nil

	srvCon, err := f.ConnAndLogin()
	if err != nil {
		return fmt.Errorf("error connecting ftp server: %v", err)
	}

	localF, err := os.Open(s.LocalDir + "/" + s.FileName)
	if err != nil {
		return fmt.Errorf("error opening a file %v", err)
	}

	err = srvCon.Stor(s.FtpDir+"/"+s.FileName, localF)

	if err != nil {
		return fmt.Errorf("error storing a file to ftp %v", err)
	}

	return nil
}

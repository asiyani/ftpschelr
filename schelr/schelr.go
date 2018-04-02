package schelr

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
)

type Stream int

const (
	//Download and upload Stream type
	Download Stream = iota
	Upload
)

type job struct {
	ID        string        `json:"id"`
	FtpDir    string        `json:"ftp_dir"`
	LocalDir  string        `json:"local_dir"`
	FileName  string        `json:"filename"`
	Direction Stream        `json:"direction"`
	StartAt   time.Time     `json:"startAt"`
	Interval  time.Duration `json:"interval"`
	NextRun   time.Time     `json:"next_run"`
	PastRuns  []time.Time   `json:"pastRuns"`
	Ticker    *time.Timer   `json:"-"`
}

// Connection is data type for ftp scheduler
type Connection struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	SerAddr string `json:"server_add"`
	User    string `json:"username"`
	Pass    string `json:"password"`
	Jobs    []job  `json:"jobs"`
}

// Connector is interface for Connection Type
type Connector interface {
	ConnAndLogin() (*ftp.ServerConn, error)
	ScheduleJob(j *job)
	CancelJobs(j *job)
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

// CreateJob creates new Jobs and adds to Connection.
func (f *Connection) CreateJob(fDir, lDir, fName string, d Stream, t time.Time, intr time.Duration) {
	s := job{
		ID:        strconv.FormatInt(time.Now().Unix()+rand.Int63(), 32),
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
func (f *Connection) ScheduleJob(j *job) {
	//var j *job
	var starDur time.Duration
	// for i, job := range f.Jobs {
	// 	if job.ID == id {
	// 		j = &f.Jobs[i]
	// 		break
	// 	}
	// }

	//If no future NextRun then then exit.
	if j.NextRun.Sub(time.Now()) < 0 {
		return
	}

	if j.StartAt.Sub(time.Now()) < 0 {
		starDur = j.NextRun.Sub(time.Now())
	} else {
		starDur = j.StartAt.Sub(time.Now())
	}

	j.Ticker = time.AfterFunc(starDur, func() {
		if j.Direction == Download {
			downloader(f, *j)
		} else {
			uploader(f, *j)
		}
		j.PastRuns = append(j.PastRuns, j.NextRun)
		updateSchedule(j)
		f.ScheduleJob(j)
	})

}

// CancelJobs already scheduled and future Jobs.
func (f *Connection) CancelJobs(j *job) {
	// var j *job
	// for i, job := range f.Jobs {
	// 	if job.ID == id {
	// 		j = &f.Jobs[i]
	// 		break
	// 	}
	// }

	j.Ticker.Stop()
	j.Interval = (0 * time.Second)
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

package ftpschelr

import (
	"fmt"
	"io"
	"log"
	"os"
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
	ftpDir    string
	localDir  string
	fileName  string
	direction stream
	startAt   time.Time
	interval  time.Duration
	nextRun   time.Time
	pastRuns  []time.Time
	ticker    *time.Timer
}

// Schedule is data type for ftp scheduler
type Schedule struct {
	Name    string
	SerAddr string
	User    string
	Pass    string
	Jobs    []job
}

// Scheduler is interface for Schedule Type
type Scheduler interface {
	ConnAndLogin() (*ftp.ServerConn, error)
	CreateSchedules(fDir, lDir, fName string, d stream, t time.Time, inv time.Duration)
	ScheduleJob(index int)
	CancelJobs(index int)
}

// ConnAndLogin will Connect and login to ftp server.
func (f *Schedule) ConnAndLogin() (*ftp.ServerConn, error) {
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

// CreateSchedules creates new Jobs and adds to Schedule.
func (f *Schedule) CreateSchedules(fDir, lDir, fName string, d stream, t time.Time, intr time.Duration) {
	s := job{
		ftpDir:    fDir,
		localDir:  lDir,
		fileName:  fName,
		direction: d,
		startAt:   t,
		interval:  intr,
	}

	s.nextRun = s.startAt.Add(s.interval)

	f.Jobs = append(f.Jobs, s)
}

// ScheduleJob will schedule Jobs[index].
func (f *Schedule) ScheduleJob(index int) {

	var starDur time.Duration

	//If no future nextRun then then exit.
	if f.Jobs[index].nextRun.Sub(time.Now()) < 0 {
		return
	}

	if f.Jobs[index].startAt.Sub(time.Now()) < 0 {
		starDur = f.Jobs[index].nextRun.Sub(time.Now())
	} else {
		starDur = f.Jobs[index].startAt.Sub(time.Now())
	}

	f.Jobs[index].ticker = time.AfterFunc(starDur, func() {
		if f.Jobs[index].direction == Download {
			downloader(f, f.Jobs[index])
		} else {
			uploader(f, f.Jobs[index])
		}
		f.Jobs[index].pastRuns = append(f.Jobs[index].pastRuns, f.Jobs[index].nextRun)
		updateSchedule(&f.Jobs[index])
		f.ScheduleJob(index)
	})

}

// CancelJobs already scheduled and future Jobs.
func (f *Schedule) CancelJobs(index int) {
	f.Jobs[index].ticker.Stop()
	f.Jobs[index].interval = (0 * time.Second)
}

func updateSchedule(s *job) {
	if s.interval != 0 {
		s.nextRun = s.nextRun.Add(s.interval)
	}
}

// GetList will get list of all files in path directory.
func getList(f Scheduler, path string) ([]*ftp.Entry, error) {

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

// Downloader will download file from ftp server (ftpDir) and store in to local drive (localDir).
func downloader(f Scheduler, s job) error {

	log.Printf("%s\n", "Downloading a file.."+s.fileName)
	return nil

	srvCon, err := f.ConnAndLogin()
	if err != nil {
		return fmt.Errorf("error connecting ftp server: %v", err)
	}

	resp, err := srvCon.Retr(s.ftpDir + "/" + s.fileName)
	if err != nil {
		return fmt.Errorf("error Retrieving  file %v", err)
	}
	defer resp.Close()

	out, err := os.Create(s.localDir + "/" + s.fileName)
	if err != nil {
		return fmt.Errorf("error creating file %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp); err != nil {
		return fmt.Errorf("copying data from conn to file %v", err)
	}

	return nil
}

// Uploader will Upload file from local drive (localDir) to ftp server dir (ftpDir).
func uploader(f Scheduler, s job) error {

	log.Printf("%s\n", "Uploading a file.."+s.fileName)
	return nil

	srvCon, err := f.ConnAndLogin()
	if err != nil {
		return fmt.Errorf("error connecting ftp server: %v", err)
	}

	localF, err := os.Open(s.localDir + "/" + s.fileName)
	if err != nil {
		return fmt.Errorf("error opening a file %v", err)
	}

	err = srvCon.Stor(s.ftpDir+"/"+s.fileName, localF)

	if err != nil {
		return fmt.Errorf("error storing a file to ftp %v", err)
	}

	return nil
}

// func main() {

// 	f1 := &Schedule{Name:"f1", SerAddr: "speedtest.tele2.net:21", User: "anonymous", Pass: ""}
// 	//f1 := &Schedule{SerAddr: "test.rebex.net:21", User: "demo", Pass: "Password"}

// 	f1.CreateSchedules(".", "./cmd", "10MB.zip", Download, time.Now(), (10 * time.Second))
// 	//f1.CreateSchedules("./upload", "./cmd", "upload_file.txt", Download, time.Now(), 0)
// 	f1.CreateSchedules("./upload", "./cmd", "local10MB.zip", Upload, time.Now(), (20 * time.Second))

// 	log.Print("Calling shce job func")

// 	for i := range f1.Jobs {
// 		f1.ScheduleJob(i)
// 	}

// 	time.Sleep(29 * time.Second)
// 	log.Print("making interval zero")
// 	f1.CancelJobs(1)

// 	time.Sleep(65 * time.Second)

// 	// log.Printf("%s\n", "Downloading a file......")
// 	// if err := downloader(f1, f1.Jobs[0]); err != nil {
// 	// 	log.Fatal("Error downloading:", err)
// 	// }
// 	// log.Printf("%s\n", "Uploading a file......")
// 	// if err := uploader(f1, f1.Jobs[1]); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// entries, err := getList(f1, "./upload")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// for _, e := range entries {
// 	// 	fmt.Printf("%s - %v\n", e.Name, e.Type)
// 	// }

// }

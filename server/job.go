package webapp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/asiyani/ftpschelr/schelr"

	"github.com/gorilla/mux"
)

func listJobsHandler(w http.ResponseWriter, r *http.Request) {
	cid := mux.Vars(r)["cid"]
	c, err := dB.Restore(cid)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring connection data", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(c.Jobs)
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	cid := mux.Vars(r)["cid"]
	c, err := dB.Restore(cid)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring connection data", http.StatusNotFound)
		return
	}

	type data struct {
		FtpDir    string        `json:"ftp_dir"`
		LocalDir  string        `json:"local_dir"`
		FileName  string        `json:"filename"`
		Direction int           `json:"direction"`
		StartAt   time.Time     `json:"startAt"`
		Interval  time.Duration `json:"interval"`
	}
	var d data
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "Error while decoding request body.", http.StatusBadRequest)
		return
	}

	c.CreateJob(d.FtpDir, d.LocalDir, d.FileName, schelr.Stream(d.Direction), d.StartAt, d.Interval)
	newJob := &c.Jobs[len(c.Jobs)-1]
	c.ScheduleJob(newJob)
	json.NewEncoder(w).Encode(newJob)
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {
	cid := mux.Vars(r)["cid"]
	jid := mux.Vars(r)["jid"]
	c, err := dB.Restore(cid)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring connection data", http.StatusNotFound)
		return
	}
	for _, j := range c.Jobs {
		if j.ID == jid {
			json.NewEncoder(w).Encode(j)
			return
		}
	}

	http.Error(w, fmt.Sprintf("Cant find job with id:%s in connection with id:%s", jid, cid), http.StatusNotFound)
}

// func updateJobHandler(w http.ResponseWriter, r *http.Request) {
// 	//id := mux.Vars(r)["id"]
// }

func deleteJobHandler(w http.ResponseWriter, r *http.Request) {
	cid := mux.Vars(r)["cid"]
	jid := mux.Vars(r)["jid"]
	c, err := dB.Restore(cid)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring connection data", http.StatusNotFound)
		return
	}
	for i, j := range c.Jobs {
		if j.ID == jid {
			c.CancelJobs(&c.Jobs[i])
			c.Jobs[len(c.Jobs)-1], c.Jobs[i] = c.Jobs[i], c.Jobs[len(c.Jobs)-1]
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "OK")
			return
		}
	}

	http.Error(w, fmt.Sprintf("Cant find job with id:%s in connection with id:%s", jid, cid), http.StatusNotFound)
}

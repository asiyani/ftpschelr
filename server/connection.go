package webapp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/asiyani/ftpschelr/db"
	"github.com/asiyani/ftpschelr/schelr"
	"github.com/gorilla/mux"
)

var (
	dB = db.New()
)

// ListConnectionsHandler will send list of all connection in json formate
func listConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := dB.RestoreAll()
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(c)
}

// createConnectionHandler will create new connection and return it back in json formate
func createConnectionHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Name    string `json:"name"`
		SerAddr string `json:"serAddr"`
		User    string `json:"username"`
		Pass    string `json:"password"`
	}
	var d data
	http.Error(w, "Please send a request body", http.StatusBadRequest)
	if r.Body == nil {
		return
	}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "Error while decoding request body.", http.StatusBadRequest)
		return
	}
	f := schelr.NewConnection(d.Name, d.SerAddr, d.User, d.Pass)
	err = dB.Store(*f)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while storing data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(f)
}

// getConnectionHandler will send connection back in json formate
func getConnectionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	c, err := dB.Restore(id)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring data", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(c)
}

// updateConnectionHandler will update existing connection and return it back.
func updateConnectionHandler(w http.ResponseWriter, r *http.Request) {
	var cnew schelr.Connection
	id := mux.Vars(r)["id"]
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&cnew)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "Error while decoding request body.", http.StatusBadRequest)
		return
	}
	if id != cnew.ID {
		http.Error(w, "Connection ID in body and url is different.", http.StatusBadRequest)
		return
	}
	err = dB.Store(cnew)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while storing data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cnew)
}

// deleteConnectionHandler will delete existing connection and return deleted connection.
func deleteConnectionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := dB.Remove(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

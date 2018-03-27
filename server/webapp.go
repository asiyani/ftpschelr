package webapp

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NewServer will run new server for web interface.
func NewServer() {
	router := mux.NewRouter()

	// connecion
	router.HandleFunc("/api/v1/connections", listConnectionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/connection", createConnectionHandler).Methods("POST")
	router.HandleFunc("/api/v1/connection/{id}", getConnectionHandler).Methods("GET")
	router.HandleFunc("/api/v1/connection/{id}", updateConnectionHandler).Methods("PUT")
	router.HandleFunc("/api/v1/connection/{id}", deleteConnectionHandler).Methods("DELETE")

	// jobs

	dir := "./server/static/"
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))

	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}

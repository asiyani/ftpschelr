package webapp

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NewServer will run new server for web interface.
func NewServer() {
	router := mux.NewRouter()
	dir := "./server/static/"

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))

	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}

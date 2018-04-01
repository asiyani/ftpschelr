package webapp

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var funcMap = template.FuncMap{

	"inc": func(i int) int {
		return i + 1
	},

	"formateDate": func(t time.Time) string {
		return t.Format("02-01-2006 15:04:05")
	},
}

// compile all templates and cache them
var templt = template.Must(template.New("mytemplate").Funcs(funcMap).ParseGlob("./server/static/template/*"))

func loadConnectionHandler(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	c, err := dB.Restore(id)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "DB Error while restoring data", http.StatusNotFound)
		return
	}

	err = templt.ExecuteTemplate(w, "connection", c)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "Error Executing template file.", http.StatusInternalServerError)
		return
	}
}

func loadIndexHandler(w http.ResponseWriter, r *http.Request) {

	err := templt.ExecuteTemplate(w, "index", nil)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, "Error Executing template file.", http.StatusInternalServerError)
		return
	}
}

// NewServer will run new server for web interface.
func NewServer() {
	router := mux.NewRouter()

	//webpage routes
	router.HandleFunc("/connection/{id}", loadConnectionHandler)

	// connecion
	router.HandleFunc("/api/v1/connections", listConnectionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/connection", createConnectionHandler).Methods("POST")
	router.HandleFunc("/api/v1/connection/{id}", getConnectionHandler).Methods("GET")
	router.HandleFunc("/api/v1/connection/{id}", updateConnectionHandler).Methods("PUT")
	router.HandleFunc("/api/v1/connection/{id}", deleteConnectionHandler).Methods("DELETE")

	// jobs

	router.HandleFunc("/", loadIndexHandler)

	dir := "./server/static/"
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))

	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}

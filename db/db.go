package db

import (
	"encoding/json"
	"fmt"

	"github.com/asiyani/ftpschelr/schelr"
	scribble "github.com/nanobox-io/golang-scribble"
)

// DB is a wrapper to scribble.Driver
type DB struct {
	*scribble.Driver
}

// New creates a new database.
func New() (*DB, error) {
	db, err := scribble.New("./jsonDB", nil)
	if err != nil {
		return nil, fmt.Errorf("error initialising database: %v", err)
	}

	return &DB{db}, nil
}

// Store will write ftp schedule to json file on disk.
func (d DB) Store(f schelr.Connection) error {

	err := d.Write("ftpschelr", f.ID, f)
	return err

}

// Restore will write ftp schedule to json file on disk.
func (d DB) Restore(ID string) (schelr.Connection, error) {

	var temp schelr.Connection
	err := d.Read("ftpschelr", ID, &temp)
	if err != nil {
		return schelr.Connection{}, fmt.Errorf("error restoring data from db %v", err)
	}
	return temp, nil

}

// RestoreAll will restore all schedules from db (json file)
func (d DB) RestoreAll() ([]schelr.Connection, error) {

	var fs []schelr.Connection
	records, err := d.ReadAll("ftpschelr")
	if err != nil {
		return nil, fmt.Errorf("error reading data from db %v", err)
	}

	for _, f := range records {
		fFound := schelr.Connection{}
		if err := json.Unmarshal([]byte(f), &fFound); err != nil {
			return nil, fmt.Errorf("error decoding data from db %v", err)
		}
		fs = append(fs, fFound)
	}

	return fs, nil
}

// Remove will delete entry for given ID
func (d DB) Remove(ID string) error {
	err := d.Delete("ftpschelr", ID)
	return err
}

// RemoveAll will delete will entries from db.
func (d DB) RemoveAll() error {

	err := d.Delete("ftpschelr", "")
	return err
}

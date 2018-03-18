package db

import (
	"encoding/json"
	"fmt"

	"github.com/asiyani/ftpschelr/ftpschelr"
	scribble "github.com/nanobox-io/golang-scribble"
)

// DB is a wrapper to scribble.Driver
type DB struct {
	*scribble.Driver
}

// New creates a new database.
func New() (*DB, error) {
	db, err := scribble.New("./jsondb", nil)
	if err != nil {
		return nil, fmt.Errorf("error initialising database: %v", err)
	}

	return &DB{db}, nil
}

// Store will write ftp schedule to json file on disk.
func (d DB) Store(f ftpschelr.Schedule) error {

	err := d.Write("ftpschelr", f.Name, f)
	return err

}

// Restore will write ftp schedule to json file on disk.
func (d DB) Restore(name string) (ftpschelr.Schedule, error) {

	var temp ftpschelr.Schedule
	err := d.Read("ftpschelr", name, &temp)
	if err != nil {
		return ftpschelr.Schedule{}, fmt.Errorf("error restoring data from db %v", err)
	}
	fmt.Printf("%v\n", temp)
	return temp, nil

}

// RestoreAll will restore all schedules from db (json file)
func (d DB) RestoreAll() ([]ftpschelr.Schedule, error) {

	var fs []ftpschelr.Schedule
	records, err := d.ReadAll("ftpschelr")
	if err != nil {
		return nil, fmt.Errorf("error reading data from db %v", err)
	}

	for _, f := range records {
		fFound := ftpschelr.Schedule{}
		if err := json.Unmarshal([]byte(f), &fFound); err != nil {
			return nil, fmt.Errorf("error decoding data from db %v", err)
		}
		fs = append(fs, fFound)
	}

	return fs, nil
}

// Remove will delete entry for name
func (d DB) Remove(name string) error {
	err := d.Delete("ftpschelr", name)
	return err
}

// RemoveAll will delete will entries from db.
func (d DB) RemoveAll() error {

	err := d.Delete("ftpschelr", "")
	return err
}

package storage

//go:generate sh -c "echo 'package storage\n\n// Do not edit. Generated code.\n\nconst sqlSchema=`' > schema.go; cat schema.sql >> schema.go; echo '`' >> schema.go" --
//go:generate sh -c "echo 'package storage\n\n// Do not edit. Generated code.\n\nconst sqlTestRecords=`' > testdb.go; cat testdb.sql >> testdb.go; echo '`' >> testdb.go" --

import (
	"database/sql"
	"log"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/config"
	_ "github.com/mattn/go-sqlite3"
)

// Storage encapsulates the state of the database.
// All CRUD-operations are defined for this type, which makes the code more readable.
type Storage struct {
	db *sql.DB
}

// Open connects to the SQLite database and initializes the schema if not done yet.
func Open(conf *config.StorageConfiguration) *Storage {

	db, err := sql.Open("sqlite3", conf.DBFilename)
	if err != nil {
		log.Fatalf("Connecting to database '%s' has failed. %s", conf.DBFilename, err)
	}

	_, err = db.Exec(sqlSchema)
	if err != nil {
		log.Fatalf("Setting up database failed. %s", err)
	}

	log.Printf("Connecting to database '%s' was successfull.", conf.DBFilename)

	return &Storage{db: db}
}

// Close writes all pending transactions to the database and disconnects from it.
func (st *Storage) Close() {

	if err := st.db.Close(); err != nil {
		log.Printf("Closing database failed. %s", err)
	}

	log.Print("Closing database was successfull.")
}

package docdb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Creates a new connection to a local database, initializing the database and
// the directory to store files if they do not already exist.
func New(dbPath string, filesPath string) IDocDB {
	return new(dbPath, filesPath)
}

func new(dbPath string, filesPath string) *DocDB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	docDb := DocDB{
		db:        db,
		filesPath: filesPath,
	}

	sqlStmt := `
 CREATE TABLE IF NOT EXISTS docs (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL
 );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}

	if _, err := os.Stat(filesPath); os.IsNotExist(err) {
		err = os.Mkdir(filesPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating files folder: %q\n", err)
		}
	}

	return &docDb
}

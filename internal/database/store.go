package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID    int
	TEXT  string
	STATE int
}

type Database struct {
	conn *sql.DB
}

func InitDB(dbName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        state INTEGER NOT NULL
    )`

	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &Database{conn: db}, nil
}

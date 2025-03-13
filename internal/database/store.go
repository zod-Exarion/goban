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

	// NOTE: Ensure the database is reachable
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// NOTE: Create table if it doesn't exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        state INTEGER NOT NULL
    )`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	// NOTE: Return the database instance
	return &Database{conn: db}, nil
}

func (db *Database) SaveTask(task *Task, updateOnConflict bool) error {
	insertCmd := `INSERT INTO tasks (text, state) VALUES (?, ?)`

	if updateOnConflict {
		insertCmd += ` ON CONFLICT(id) DO UPDATE 
                       SET text = EXCLUDED.text, 
                           state = EXCLUDED.state`
	} else {
		insertCmd += ` ON CONFLICT(id) DO NOTHING`
	}

	if _, err := db.conn.Exec(insertCmd, task.TEXT, task.STATE); err != nil {
		return fmt.Errorf("error saving task to database: %v", err)
	}

	return nil
}

func (db *Database) DeleteTask(id int) error {
	deleteCmd := `DELETE FROM tasks WHERE id = ?`

	if _, err := db.conn.Exec(deleteCmd, id); err != nil {
		return fmt.Errorf("error DELETING task from database: %v", err)
	}

	return nil
}

func (db *Database) EditTask(id int, newText string) error {
	updateQuery := `UPDATE tasks SET text = ? WHERE id = ?`

	result, err := db.conn.Exec(updateQuery, newText, id)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return nil
}

func (db *Database) MarkTask(id int) error {
	updateQuery := `UPDATE tasks SET state = state + 1 WHERE id = ?`

	result, err := db.conn.Exec(updateQuery, id)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	var newState int
	err = db.conn.QueryRow("SELECT state FROM tasks WHERE id = ?", id).Scan(&newState)
	if err != nil {
		return fmt.Errorf("error fetching updated task state: %v", err)
	}

	if newState > 2 {
		return db.DeleteTask(id)
	}

	return nil
}

func (db *Database) GetAllTasks() ([]Task, error) {
	fetchQuery := `SELECT id,text,state FROM tasks`
	rows, err := db.conn.Query(fetchQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks from database: %v", err)
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.TEXT, &task.STATE); err != nil {
			return nil, fmt.Errorf("error scanning task from database: %v", err)
		}
		tasks = append(tasks, task)
	}

	// NOTE: Extra Care
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tasks from database: %v", err)
	}

	return tasks, nil
}

func (db *Database) GetTaskByState(taskState int) ([]Task, error) {
	fetchQuery := `SELECT id,text,state FROM tasks WHERE state = ?`
	rows, err := db.conn.Query(fetchQuery, taskState)
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks from database: %v", err)
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.TEXT, &task.STATE); err != nil {
			return nil, fmt.Errorf("error scanning task from database: %v", err)
		}
		tasks = append(tasks, task)
	}

	// NOTE: Extra Care
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tasks from database: %v", err)
	}

	return tasks, nil
}

func (db *Database) NukeDB() error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("error BEGINNING transaction: %v", err)
	}

	nukeCmd := `DELETE FROM tasks`
	if _, err := db.conn.Exec(nukeCmd); err != nil {
		return fmt.Errorf("error NUKING database: %v", err)
	}

	resetSequenceCmd := `DELETE FROM sqlite_sequence WHERE name = 'tasks'`
	if _, err := db.conn.Exec(resetSequenceCmd); err != nil {
		return fmt.Errorf("error RESETTING SEQUENCE: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error COMMITING transaction: %v", err)
	}
	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

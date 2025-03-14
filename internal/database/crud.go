package database

import "fmt"

func (db *Database) SaveTask(task *Task, updateOnConflict bool) error {
	insertCmd := `INSERT INTO tasks (text, state) VALUES (?, ?)`

	if updateOnConflict {
		insertCmd += ` ON CONFLICT(id) DO UPDATE SET text = EXCLUDED.text, state = EXCLUDED.state`
	} else {
		insertCmd += ` ON CONFLICT(id) DO NOTHING`
	}

	_, err := db.conn.Exec(insertCmd, task.TEXT, task.STATE)
	return err
}

func (db *Database) DeleteTask(id int) error {
	_, err := db.conn.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func (db *Database) EditTask(id int, newText string) error {
	result, err := db.conn.Exec(`UPDATE tasks SET text = ? WHERE id = ?`, newText, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	return nil
}

func (db *Database) MarkTask(id int) error {
	result, err := db.conn.Exec(`UPDATE tasks SET state = state + 1 WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	var newState int
	if err := db.conn.QueryRow("SELECT state FROM tasks WHERE id = ?", id).Scan(&newState); err != nil {
		return err
	}

	if newState > 2 {
		return db.DeleteTask(id)
	}
	return nil
}

func (db *Database) GetAllTasks() ([]Task, error) {
	return db.queryTasks(`SELECT id, text, state FROM tasks`)
}

func (db *Database) GetTaskByState(taskState int) ([]Task, error) {
	return db.queryTasks(`SELECT id, text, state FROM tasks WHERE state = ?`, taskState)
}

func (db *Database) GetAllTasksSorted(sortOption int) ([]Task, error) {
	query := `SELECT id, text, state FROM tasks ORDER BY `
	query += getSortClause(sortOption)
	return db.queryTasks(query)
}

func (db *Database) GetTaskByStateSorted(taskState int, sortOption int) ([]Task, error) {
	query := `SELECT id, text, state FROM tasks WHERE state = ? ORDER BY `
	query += getSortClause(sortOption)
	return db.queryTasks(query, taskState)
}

func getSortClause(sortOption int) string {
	switch sortOption {
	case 1:
		return `id DESC`
	case 2:
		return `LENGTH(text) ASC`
	case 3:
		return `LENGTH(text) DESC`
	default:
		return `id ASC`
	}
}

func (db *Database) queryTasks(query string, args ...interface{}) ([]Task, error) {
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.TEXT, &task.STATE); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (db *Database) NukeDB() error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	if _, err := db.conn.Exec(`DELETE FROM tasks`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := db.conn.Exec(`DELETE FROM sqlite_sequence WHERE name = 'tasks'`); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *Database) Close() error {
	return db.conn.Close()
}

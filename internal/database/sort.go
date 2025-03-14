package database

import "fmt"

func (db *Database) SortByOldestFirst() error {
	sortQuery := `CREATE TEMP TABLE sorted AS SELECT * FROM tasks ORDER BY id ASC;
				  DELETE FROM tasks;
				  INSERT INTO tasks SELECT * FROM sorted;
				  DROP TABLE sorted;`

	_, err := db.conn.Exec(sortQuery)
	if err != nil {
		return fmt.Errorf("error sorting tasks by oldest first: %v", err)
	}
	return nil
}

func (db *Database) SortByNewestFirst() error {
	sortQuery := `CREATE TEMP TABLE sorted AS SELECT * FROM tasks ORDER BY id DESC;
				  DELETE FROM tasks;
				  INSERT INTO tasks SELECT * FROM sorted;
				  DROP TABLE sorted;`

	_, err := db.conn.Exec(sortQuery)
	if err != nil {
		return fmt.Errorf("error sorting tasks by newest first: %v", err)
	}
	return nil
}

func (db *Database) SortByTextLengthAsc() error {
	sortQuery := `CREATE TEMP TABLE sorted AS SELECT * FROM tasks ORDER BY LENGTH(text) ASC;
				  DELETE FROM tasks;
				  INSERT INTO tasks SELECT * FROM sorted;
				  DROP TABLE sorted;`

	_, err := db.conn.Exec(sortQuery)
	if err != nil {
		return fmt.Errorf("error sorting tasks by text length (shortest first): %v", err)
	}
	return nil
}

func (db *Database) SortByTextLengthDesc() error {
	sortQuery := `CREATE TEMP TABLE sorted AS SELECT * FROM tasks ORDER BY LENGTH(text) DESC;
				  DELETE FROM tasks;
				  INSERT INTO tasks SELECT * FROM sorted;
				  DROP TABLE sorted;`

	_, err := db.conn.Exec(sortQuery)
	if err != nil {
		return fmt.Errorf("error sorting tasks by text length (longest first): %v", err)
	}
	return nil
}

package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const sqliteDBFile = "store.db"

func openSQLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sqliteDBFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			state TEXT NOT NULL DEFAULT 'pending'
		)
	`)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func writeTaskToSQLite(task Task) error {
	db, err := openSQLiteDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO tasks (id, name, description, state) VALUES (?, ?, ?, ?)`,
		task.ID,
		task.Name,
		task.Description,
		task.State,
	)
	if err != nil {
		return err
	}

	return nil
}

func readTasksFromSQLite() ([]Task, error) {
	db, err := openSQLiteDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, name, description, state FROM tasks ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.State)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

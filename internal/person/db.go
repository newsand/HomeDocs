package person

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func Open(p_Path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(p_Path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", p_Path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func Migrate(p_DB *sql.DB) error {
	names := columnNames()
	parts := make([]string, len(names))
	for i, name := range names {
		if name == "id" {
			parts[i] = "id TEXT PRIMARY KEY"
			continue
		}
		parts[i] = name + " TEXT NOT NULL DEFAULT ''"
	}
	people := "CREATE TABLE IF NOT EXISTS people (" + strings.Join(parts, ", ") + ")"
	if _, err := p_DB.Exec(people); err != nil {
		return err
	}
	_, err := p_DB.Exec(`CREATE TABLE IF NOT EXISTS person_scan (
		person_id TEXT NOT NULL,
		position INTEGER NOT NULL,
		path TEXT NOT NULL,
		PRIMARY KEY (person_id, position),
		FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE
	)`)
	return err
}

func insertPerson(p_Ctx context.Context, p_Tx *sql.Tx, p_Person Person) error {
	names := columnNames()
	marks := make([]string, len(names))
	for i := range names {
		marks[i] = "?"
	}
	query := "INSERT INTO people (" + strings.Join(names, ", ") + ") VALUES (" + strings.Join(marks, ", ") + ")"
	_, err := p_Tx.ExecContext(p_Ctx, query, columnValues(p_Person)...)
	return err
}

func updatePersonRow(p_Ctx context.Context, p_Tx *sql.Tx, p_Person Person) error {
	names := columnNames()[1:]
	sets := make([]string, len(names))
	for i, name := range names {
		sets[i] = name + " = ?"
	}
	values := columnValues(p_Person)[1:]
	values = append(values, p_Person.ID)
	query := "UPDATE people SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	result, err := p_Tx.ExecContext(p_Ctx, query, values...)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func replaceScans(p_Ctx context.Context, p_Tx *sql.Tx, p_ID string, p_Paths []string) error {
	if _, err := p_Tx.ExecContext(p_Ctx, `DELETE FROM person_scan WHERE person_id = ?`, p_ID); err != nil {
		return err
	}
	for i, path := range p_Paths {
		if _, err := p_Tx.ExecContext(p_Ctx, `INSERT INTO person_scan (person_id, position, path) VALUES (?, ?, ?)`, p_ID, i, path); err != nil {
			return err
		}
	}
	return nil
}

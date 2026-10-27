package person

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

func ShowPerson(p_Ctx context.Context, p_DB *sql.DB, p_ID string) (Person, error) {
	row := p_DB.QueryRowContext(p_Ctx, "SELECT "+strings.Join(columnNames(), ", ")+" FROM people WHERE id = ?", p_ID)
	person, err := scanPerson(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return Person{}, ErrNotFound
	}
	if err != nil {
		return Person{}, err
	}
	scans, err := loadScans(p_Ctx, p_DB, person.ID)
	if err != nil {
		return Person{}, err
	}
	person.DocumentScans = scans
	return person, nil
}

func scanPerson(p_Scan func(p_Dest ...any) error) (Person, error) {
	var person Person
	if err := p_Scan(scanDest(&person)...); err != nil {
		return Person{}, err
	}
	return person, nil
}

func loadScans(p_Ctx context.Context, p_DB *sql.DB, p_ID string) ([]string, error) {
	rows, err := p_DB.QueryContext(p_Ctx, `SELECT path FROM person_scan WHERE person_id = ? ORDER BY position`, p_ID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	paths := []string{}
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return paths, nil
}

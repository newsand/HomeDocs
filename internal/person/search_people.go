package person

import (
	"context"
	"database/sql"
	"strings"
)

func SearchPeople(p_Ctx context.Context, p_DB *sql.DB, p_Query string) ([]Person, error) {
	query := "SELECT " + strings.Join(columnNames(), ", ") + " FROM people"
	text := strings.TrimSpace(p_Query)
	args := []any{}
	if text != "" {
		needle := likeNeedle(text)
		cols := searchColumns()
		parts := make([]string, len(cols))
		for i, col := range cols {
			parts[i] = "lower(" + col + ") LIKE lower(?) ESCAPE '\\'"
			args = append(args, needle)
		}
		query += " WHERE " + strings.Join(parts, " OR ")
	}
	query += " ORDER BY name COLLATE NOCASE, id"
	rows, err := p_DB.QueryContext(p_Ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	people := []Person{}
	for rows.Next() {
		person, err := scanPerson(rows.Scan)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range people {
		scans, err := loadScans(p_Ctx, p_DB, people[i].ID)
		if err != nil {
			return nil, err
		}
		people[i].DocumentScans = scans
	}
	return people, nil
}

func likeNeedle(p_Query string) string {
	var b strings.Builder
	b.WriteByte('%')
	for _, char := range p_Query {
		if char == '%' || char == '_' || char == '\\' {
			b.WriteByte('\\')
		}
		b.WriteRune(char)
	}
	b.WriteByte('%')
	return b.String()
}

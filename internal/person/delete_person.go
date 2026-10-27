package person

import (
	"context"
	"database/sql"
)

func DeletePerson(p_Ctx context.Context, p_DB *sql.DB, p_Dir, p_ID string) error {
	current, err := ShowPerson(p_Ctx, p_DB, p_ID)
	if err != nil {
		return err
	}
	tx, err := p_DB.BeginTx(p_Ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	result, err := tx.ExecContext(p_Ctx, `DELETE FROM people WHERE id = ?`, p_ID)
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
	if err := tx.Commit(); err != nil {
		return err
	}
	names := []string{}
	if current.ProfilePhoto != "" {
		names = append(names, current.ProfilePhoto)
	}
	names = append(names, current.DocumentScans...)
	return removeNames(p_Dir, names)
}

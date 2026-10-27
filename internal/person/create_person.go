package person

import (
	"context"
	"database/sql"
	"strings"
)

func CreatePerson(p_Ctx context.Context, p_DB *sql.DB, p_Dir string, p_Input Input) (person Person, err error) {
	if !ValidName(p_Input.Name) {
		return Person{}, ErrNameRequired
	}
	var written []string
	photo := ""
	if p_Input.ProfilePhoto != nil {
		photo, err = writeUpload(p_Dir, *p_Input.ProfilePhoto)
		if err != nil {
			return Person{}, err
		}
		written = append(written, photo)
	}
	scans, err := writeUploads(p_Dir, p_Input.DocumentScans)
	if err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	written = append(written, scans...)
	id, err := newID()
	if err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	record := Person{
		ID:           id,
		Name:         strings.TrimSpace(p_Input.Name),
		ProfilePhoto: photo,
		CNH:          p_Input.CNH,
		CPF:          p_Input.CPF,
		RG:           p_Input.RG,
		VoterID:      p_Input.VoterID,
		Passport:     p_Input.Passport,
	}
	tx, err := p_DB.BeginTx(p_Ctx, nil)
	if err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if err = insertPerson(p_Ctx, tx, record); err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	if err = replaceScans(p_Ctx, tx, record.ID, scans); err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	if err = tx.Commit(); err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	return ShowPerson(p_Ctx, p_DB, record.ID)
}

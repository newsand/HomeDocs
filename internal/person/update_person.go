package person

import (
	"context"
	"database/sql"
	"strings"
)

func UpdatePerson(p_Ctx context.Context, p_DB *sql.DB, p_Dir, p_ID string, p_Input Input) (Person, error) {
	if !ValidName(p_Input.Name) {
		return Person{}, ErrNameRequired
	}
	current, err := ShowPerson(p_Ctx, p_DB, p_ID)
	if err != nil {
		return Person{}, err
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
	var scans []string
	if p_Input.DocumentScans != nil {
		scans, err = writeUploads(p_Dir, p_Input.DocumentScans)
		if err != nil {
			if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
				return Person{}, cleanErr
			}
			return Person{}, err
		}
		written = append(written, scans...)
	}
	record := Person{
		ID:           current.ID,
		Name:         strings.TrimSpace(p_Input.Name),
		ProfilePhoto: current.ProfilePhoto,
		CNH:          p_Input.CNH,
		CPF:          p_Input.CPF,
		RG:           p_Input.RG,
		VoterID:      p_Input.VoterID,
		Passport:     p_Input.Passport,
	}
	if p_Input.ProfilePhoto != nil {
		record.ProfilePhoto = photo
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
	if err = updatePersonRow(p_Ctx, tx, record); err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	if p_Input.DocumentScans != nil {
		if err = replaceScans(p_Ctx, tx, record.ID, scans); err != nil {
			if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
				return Person{}, cleanErr
			}
			return Person{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		if cleanErr := removeNames(p_Dir, written); cleanErr != nil {
			return Person{}, cleanErr
		}
		return Person{}, err
	}
	var old []string
	if p_Input.ProfilePhoto != nil && current.ProfilePhoto != "" {
		old = append(old, current.ProfilePhoto)
	}
	if p_Input.DocumentScans != nil {
		old = append(old, current.DocumentScans...)
	}
	if cleanErr := removeNames(p_Dir, old); cleanErr != nil {
		return Person{}, cleanErr
	}
	return ShowPerson(p_Ctx, p_DB, record.ID)
}

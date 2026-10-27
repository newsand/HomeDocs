package person

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestValidName(t *testing.T) {
	t.Run("should accept text", func(t *testing.T) {
		if !ValidName("Ana") {
			t.Fatal("expected valid name")
		}
	})
	t.Run("should accept text with surrounding spaces", func(t *testing.T) {
		if !ValidName("  Ana  ") {
			t.Fatal("expected valid name")
		}
	})
	t.Run("should reject empty", func(t *testing.T) {
		if ValidName("") {
			t.Fatal("expected invalid name")
		}
	})
	t.Run("should reject spaces", func(t *testing.T) {
		if ValidName("   ") {
			t.Fatal("expected invalid name")
		}
	})
}

func TestCreatePerson(t *testing.T) {
	t.Run("should persist a person with only the name", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{Name: "  Ana  "})
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "Ana" {
			t.Fatalf("name %q", got.Name)
		}
		if got.ID == "" {
			t.Fatal("expected id")
		}
		shown, err := ShowPerson(context.Background(), db, got.ID)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Name != "Ana" {
			t.Fatalf("shown %q", shown.Name)
		}
		if len(shown.DocumentScans) != 0 || shown.ProfilePhoto != "" {
			t.Fatalf("expected no files, photo %q scans %d", shown.ProfilePhoto, len(shown.DocumentScans))
		}
	})

	t.Run("should persist one field of each document", func(t *testing.T) {
		db, dir := openTestDB(t)
		in := Input{
			Name:     "Ana",
			CNH:      CNH{RegisterNumber: "CNH-1", Category: "B"},
			CPF:      CPF{Number: "12345678900", Status: "regular"},
			RG:       RG{Number: "MG-1", SocialName: "An"},
			VoterID:  VoterID{Number: "T-1", Zone: "12"},
			Passport: Passport{Number: "P-1", IssuingCountry: "BRA"},
		}
		got, err := CreatePerson(context.Background(), db, dir, in)
		if err != nil {
			t.Fatal(err)
		}
		shown, err := ShowPerson(context.Background(), db, got.ID)
		if err != nil {
			t.Fatal(err)
		}
		if shown.CNH.RegisterNumber != "CNH-1" || shown.CNH.Category != "B" {
			t.Fatalf("cnh %+v", shown.CNH)
		}
		if shown.CPF.Number != "12345678900" || shown.CPF.Status != "regular" {
			t.Fatalf("cpf %+v", shown.CPF)
		}
		if shown.RG.Number != "MG-1" || shown.RG.SocialName != "An" {
			t.Fatalf("rg %+v", shown.RG)
		}
		if shown.VoterID.Number != "T-1" || shown.VoterID.Zone != "12" {
			t.Fatalf("voter %+v", shown.VoterID)
		}
		if shown.Passport.Number != "P-1" || shown.Passport.IssuingCountry != "BRA" {
			t.Fatalf("passport %+v", shown.Passport)
		}
	})

	t.Run("should allow two people with the same name", func(t *testing.T) {
		db, dir := openTestDB(t)
		first, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana"})
		if err != nil {
			t.Fatal(err)
		}
		second, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana"})
		if err != nil {
			t.Fatal(err)
		}
		if first.ID == second.ID {
			t.Fatal("expected distinct ids")
		}
	})

	t.Run("should store a sql shaped name as data", func(t *testing.T) {
		db, dir := openTestDB(t)
		name := "Robert'); DROP TABLE people;--"
		got, err := CreatePerson(context.Background(), db, dir, Input{Name: name})
		if err != nil {
			t.Fatal(err)
		}
		shown, err := ShowPerson(context.Background(), db, got.ID)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Name != name {
			t.Fatalf("name %q", shown.Name)
		}
		var n int
		if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM people`).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != 1 {
			t.Fatalf("rows %d", n)
		}
	})

	t.Run("should reject a blank name without writing a row", func(t *testing.T) {
		db, dir := openTestDB(t)
		_, err := CreatePerson(context.Background(), db, dir, Input{Name: "   "})
		if !errors.Is(err, ErrNameRequired) {
			t.Fatalf("err %v", err)
		}
		people, err := SearchPeople(context.Background(), db, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(people) != 0 {
			t.Fatalf("rows %d", len(people))
		}
	})

	t.Run("should keep an upload filename inside the storage dir", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{
			Name:         "Ana",
			ProfilePhoto: &Upload{Name: "../evil.png", Bytes: []byte("png-bytes")},
		})
		if err != nil {
			t.Fatal(err)
		}
		if got.ProfilePhoto == "" || strings.Contains(got.ProfilePhoto, "..") || strings.ContainsAny(got.ProfilePhoto, `/\`) {
			t.Fatalf("photo %q", got.ProfilePhoto)
		}
		full := filepath.Join(dir, got.ProfilePhoto)
		body, err := os.ReadFile(full)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "png-bytes" {
			t.Fatalf("body %q", body)
		}
		outside := filepath.Join(filepath.Dir(dir), "evil.png")
		if _, err := os.Stat(outside); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("file escaped the dir: %v", err)
		}
	})

	t.Run("should remove the file when the insert fails", func(t *testing.T) {
		db, dir := openTestDB(t)
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
		ro, err := sql.Open("sqlite", "file:"+filepath.Join(dir, "id.sqlite")+"?mode=ro")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = ro.Close() })
		files := t.TempDir()
		_, err = CreatePerson(context.Background(), ro, files, Input{
			Name:         "Ana",
			ProfilePhoto: &Upload{Name: "a.png", Bytes: []byte("png")},
		})
		if err == nil {
			t.Fatal("expected insert error")
		}
		entries, err := os.ReadDir(files)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("orphan files %d", len(entries))
		}
	})
}

func TestUpdatePerson(t *testing.T) {
	t.Run("should replace text fields", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{
			Name: "Ana",
			CPF:  CPF{Number: "1"},
			CNH:  CNH{Category: "B"},
		})
		if err != nil {
			t.Fatal(err)
		}
		updated, err := UpdatePerson(context.Background(), db, dir, got.ID, Input{
			Name: "Bia",
			CPF:  CPF{Number: "2"},
		})
		if err != nil {
			t.Fatal(err)
		}
		if updated.Name != "Bia" || updated.CPF.Number != "2" || updated.CNH.Category != "" {
			t.Fatalf("%+v", updated)
		}
	})

	t.Run("should keep the row when the name is blank", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana", CPF: CPF{Number: "1"}})
		if err != nil {
			t.Fatal(err)
		}
		_, err = UpdatePerson(context.Background(), db, dir, got.ID, Input{
			Name:         "   ",
			ProfilePhoto: &Upload{Name: "n.png", Bytes: []byte("new")},
		})
		if !errors.Is(err, ErrNameRequired) {
			t.Fatalf("err %v", err)
		}
		shown, err := ShowPerson(context.Background(), db, got.ID)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Name != "Ana" || shown.CPF.Number != "1" || shown.ProfilePhoto != "" {
			t.Fatalf("%+v", shown)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if entry.Name() != "id.sqlite" {
				t.Fatalf("unexpected file %s", entry.Name())
			}
		}
	})

	t.Run("should keep files when the upload is empty", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{
			Name:          "Ana",
			ProfilePhoto:  &Upload{Name: "face.png", Bytes: []byte("face")},
			DocumentScans: []Upload{{Name: "doc.png", Bytes: []byte("scan")}},
		})
		if err != nil {
			t.Fatal(err)
		}
		updated, err := UpdatePerson(context.Background(), db, dir, got.ID, Input{Name: "Ana"})
		if err != nil {
			t.Fatal(err)
		}
		if updated.ProfilePhoto != got.ProfilePhoto || len(updated.DocumentScans) != 1 || updated.DocumentScans[0] != got.DocumentScans[0] {
			t.Fatalf("files changed %+v", updated)
		}
	})

	t.Run("should replace scans and the profile photo", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{
			Name:          "Ana",
			ProfilePhoto:  &Upload{Name: "face.png", Bytes: []byte("face")},
			DocumentScans: []Upload{{Name: "doc.png", Bytes: []byte("scan")}},
		})
		if err != nil {
			t.Fatal(err)
		}
		updated, err := UpdatePerson(context.Background(), db, dir, got.ID, Input{
			Name:          "Ana",
			ProfilePhoto:  &Upload{Name: "face2.png", Bytes: []byte("face2")},
			DocumentScans: []Upload{{Name: "doc2.png", Bytes: []byte("scan2")}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if updated.ProfilePhoto == got.ProfilePhoto || updated.DocumentScans[0] == got.DocumentScans[0] {
			t.Fatal("expected new files")
		}
		if _, err := os.Stat(filepath.Join(dir, got.ProfilePhoto)); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("old photo remains")
		}
		if _, err := os.Stat(filepath.Join(dir, got.DocumentScans[0])); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("old scan remains")
		}
		body, err := os.ReadFile(filepath.Join(dir, updated.ProfilePhoto))
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "face2" {
			t.Fatalf("body %q", body)
		}
	})

	t.Run("should keep the previous row when the update fails", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{
			Name:         "Ana",
			ProfilePhoto: &Upload{Name: "face.png", Bytes: []byte("face")},
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.ExecContext(context.Background(), `CREATE TRIGGER fail_update BEFORE UPDATE ON people BEGIN SELECT RAISE(ABORT, 'fail'); END`)
		if err != nil {
			t.Fatal(err)
		}
		_, err = UpdatePerson(context.Background(), db, dir, got.ID, Input{
			Name:         "Bia",
			ProfilePhoto: &Upload{Name: "face2.png", Bytes: []byte("face2")},
		})
		if err == nil {
			t.Fatal("expected update error")
		}
		shown, err := ShowPerson(context.Background(), db, got.ID)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Name != "Ana" || shown.ProfilePhoto != got.ProfilePhoto {
			t.Fatalf("%+v", shown)
		}
		if _, err := os.Stat(filepath.Join(dir, got.ProfilePhoto)); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if entry.Name() != "id.sqlite" && entry.Name() != got.ProfilePhoto {
				t.Fatalf("orphan %s", entry.Name())
			}
		}
	})

	t.Run("should return not found for a missing id", func(t *testing.T) {
		db, dir := openTestDB(t)
		_, err := UpdatePerson(context.Background(), db, dir, "missing", Input{Name: "Ana"})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err %v", err)
		}
	})
}

func TestDeletePerson(t *testing.T) {
	t.Run("should remove the row and the files", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{
			Name:          "Ana",
			ProfilePhoto:  &Upload{Name: "face.png", Bytes: []byte("face")},
			DocumentScans: []Upload{{Name: "doc.png", Bytes: []byte("scan")}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := DeletePerson(context.Background(), db, dir, got.ID); err != nil {
			t.Fatal(err)
		}
		_, err = ShowPerson(context.Background(), db, got.ID)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err %v", err)
		}
		if _, err := os.Stat(filepath.Join(dir, got.ProfilePhoto)); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("photo remains")
		}
		if _, err := os.Stat(filepath.Join(dir, got.DocumentScans[0])); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("scan remains")
		}
	})

	t.Run("should return not found on the second delete", func(t *testing.T) {
		db, dir := openTestDB(t)
		got, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana"})
		if err != nil {
			t.Fatal(err)
		}
		if err := DeletePerson(context.Background(), db, dir, got.ID); err != nil {
			t.Fatal(err)
		}
		err = DeletePerson(context.Background(), db, dir, got.ID)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err %v", err)
		}
	})

	t.Run("should treat an injection shaped id as missing", func(t *testing.T) {
		db, dir := openTestDB(t)
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana"}); err != nil {
			t.Fatal(err)
		}
		err := DeletePerson(context.Background(), db, dir, "1 OR 1=1")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err %v", err)
		}
		people, err := SearchPeople(context.Background(), db, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(people) != 1 {
			t.Fatalf("rows %d", len(people))
		}
	})
}

func TestShowPerson(t *testing.T) {
	t.Run("should return not found", func(t *testing.T) {
		db, _ := openTestDB(t)
		_, err := ShowPerson(context.Background(), db, "missing")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err %v", err)
		}
	})
}

func TestSearchPeople(t *testing.T) {
	t.Run("should list everyone ordered by name then id", func(t *testing.T) {
		db, dir := openTestDB(t)
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Maria"}); err != nil {
			t.Fatal(err)
		}
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "ana"}); err != nil {
			t.Fatal(err)
		}
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Bia"}); err != nil {
			t.Fatal(err)
		}
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Bia"}); err != nil {
			t.Fatal(err)
		}
		people, err := SearchPeople(context.Background(), db, "  ")
		if err != nil {
			t.Fatal(err)
		}
		if len(people) != 4 {
			t.Fatalf("len %d", len(people))
		}
		if people[0].Name != "ana" || people[1].Name != "Bia" || people[2].Name != "Bia" || people[3].Name != "Maria" {
			t.Fatalf("order %+v", names(people))
		}
		if people[1].ID > people[2].ID {
			t.Fatalf("id order %s %s", people[1].ID, people[2].ID)
		}
	})

	t.Run("should match name and document fields without case", func(t *testing.T) {
		db, dir := openTestDB(t)
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana"}); err != nil {
			t.Fatal(err)
		}
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Leo", CNH: CNH{Notes: "Observação XYZ"}}); err != nil {
			t.Fatal(err)
		}
		byName, err := SearchPeople(context.Background(), db, "ANA")
		if err != nil {
			t.Fatal(err)
		}
		if len(byName) != 1 || byName[0].Name != "Ana" {
			t.Fatalf("%+v", names(byName))
		}
		byNote, err := SearchPeople(context.Background(), db, "xyz")
		if err != nil {
			t.Fatal(err)
		}
		if len(byNote) != 1 || byNote[0].Name != "Leo" {
			t.Fatalf("%+v", names(byNote))
		}
	})

	t.Run("should treat percent as text", func(t *testing.T) {
		db, dir := openTestDB(t)
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "Ana"}); err != nil {
			t.Fatal(err)
		}
		if _, err := CreatePerson(context.Background(), db, dir, Input{Name: "100%"}); err != nil {
			t.Fatal(err)
		}
		people, err := SearchPeople(context.Background(), db, "%")
		if err != nil {
			t.Fatal(err)
		}
		if len(people) != 1 || people[0].Name != "100%" {
			t.Fatalf("%+v", names(people))
		}
	})

	t.Run("should return an empty list", func(t *testing.T) {
		db, _ := openTestDB(t)
		people, err := SearchPeople(context.Background(), db, "ninguem")
		if err != nil {
			t.Fatal(err)
		}
		if len(people) != 0 {
			t.Fatalf("len %d", len(people))
		}
	})
}

func openTestDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	dir := t.TempDir()
	db, err := Open(filepath.Join(dir, "id.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	return db, dir
}

func names(p_People []Person) []string {
	out := make([]string, len(p_People))
	for i := range p_People {
		out[i] = p_People[i].Name
	}
	return out
}

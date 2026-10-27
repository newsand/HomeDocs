package person

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func newID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func extensionOf(p_Name string) string {
	ext := filepath.Ext(filepath.Base(p_Name))
	if len(ext) < 2 || len(ext) > 6 || ext[0] != '.' {
		return ""
	}
	for _, char := range ext[1:] {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
			continue
		}
		return ""
	}
	return ext
}

func writeUpload(p_Dir string, p_Upload Upload) (string, error) {
	if err := os.MkdirAll(p_Dir, 0o755); err != nil {
		return "", err
	}
	id, err := newID()
	if err != nil {
		return "", err
	}
	name := id + extensionOf(p_Upload.Name)
	full := filepath.Join(p_Dir, name)
	if err := os.WriteFile(full, p_Upload.Bytes, 0o644); err != nil {
		return "", err
	}
	return name, nil
}

func writeUploads(p_Dir string, p_Uploads []Upload) ([]string, error) {
	paths := make([]string, 0, len(p_Uploads))
	for _, upload := range p_Uploads {
		name, err := writeUpload(p_Dir, upload)
		if err != nil {
			if cleanErr := removeNames(p_Dir, paths); cleanErr != nil {
				return nil, cleanErr
			}
			return nil, err
		}
		paths = append(paths, name)
	}
	return paths, nil
}

func removeNames(p_Dir string, p_Names []string) error {
	var first error
	for _, name := range p_Names {
		if err := removeStored(p_Dir, name); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func removeStored(p_Dir, p_Name string) error {
	if p_Name == "" {
		return nil
	}
	full := filepath.Join(p_Dir, filepath.Base(p_Name))
	err := os.Remove(full)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func FilePath(p_Dir, p_Name string) (string, error) {
	if p_Name == "" || p_Name != filepath.Base(p_Name) || strings.Contains(p_Name, "..") {
		return "", ErrNotFound
	}
	dir, err := filepath.Abs(p_Dir)
	if err != nil {
		return "", err
	}
	full := filepath.Join(dir, p_Name)
	rel, err := filepath.Rel(dir, full)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", ErrNotFound
	}
	info, err := os.Stat(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrNotFound
		}
		return "", err
	}
	if info.IsDir() {
		return "", ErrNotFound
	}
	return full, nil
}

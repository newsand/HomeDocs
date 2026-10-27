package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"idcenter/internal/person"
)

func TestPages(t *testing.T) {
	t.Run("should redirect the root to the list", func(t *testing.T) {
		router := newRouter(t)
		rec := serve(router, http.MethodGet, "/", nil, "")
		if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/people" {
			t.Fatalf("status %d loc %s", rec.Code, rec.Header().Get("Location"))
		}
	})

	t.Run("should create with only the name and show copy buttons", func(t *testing.T) {
		router := newRouter(t)
		rec := serve(router, http.MethodPost, "/people", strings.NewReader("name=Ana"), "application/x-www-form-urlencoded")
		if rec.Code != http.StatusSeeOther {
			t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
		}
		loc := rec.Header().Get("Location")
		page := serve(router, http.MethodGet, loc, nil, "")
		if page.Code != http.StatusOK {
			t.Fatalf("status %d", page.Code)
		}
		body := page.Body.String()
		if !strings.Contains(body, "Ana") || !strings.Contains(body, "data-copy=") || !strings.Contains(body, "navigator.clipboard") {
			t.Fatalf("body %s", body)
		}
		if strings.Count(body, "data-copy=") < 56 {
			t.Fatalf("copy buttons %d", strings.Count(body, "data-copy="))
		}
		for _, label := range []string{"Doc. identidade", "Situação cadastral", "Nome social", "Zona", "País emissor"} {
			if !strings.Contains(body, label) {
				t.Fatalf("missing %s", label)
			}
		}
	})

	t.Run("should reject a blank name", func(t *testing.T) {
		router := newRouter(t)
		rec := serve(router, http.MethodPost, "/people", strings.NewReader("name=%20%20"), "application/x-www-form-urlencoded")
		if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "Informe o nome.") {
			t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
		}
		list := serve(router, http.MethodGet, "/people", nil, "")
		if !strings.Contains(list.Body.String(), "Nenhuma pessoa.") {
			t.Fatalf("body %s", list.Body.String())
		}
	})

	t.Run("should escape a script shaped name", func(t *testing.T) {
		router := newRouter(t)
		form := url.Values{}
		form.Set("name", "<script>alert(1)</script>")
		rec := serve(router, http.MethodPost, "/people", strings.NewReader(form.Encode()), "application/x-www-form-urlencoded")
		page := serve(router, http.MethodGet, rec.Header().Get("Location"), nil, "")
		body := page.Body.String()
		if strings.Contains(body, "<script>alert(1)</script>") || !strings.Contains(body, "&lt;script&gt;") {
			t.Fatalf("body %s", body)
		}
	})

	t.Run("should return not found without a sql error", func(t *testing.T) {
		router := newRouter(t)
		path := "/people/" + url.PathEscape("1 OR 1=1")
		rec := serve(router, http.MethodGet, path, nil, "")
		if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Pessoa não encontrada.") {
			t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
		}
		if strings.Contains(strings.ToLower(rec.Body.String()), "sql") {
			t.Fatalf("leaked sql: %s", rec.Body.String())
		}
	})

	t.Run("should update and delete", func(t *testing.T) {
		router := newRouter(t)
		created := serve(router, http.MethodPost, "/people", strings.NewReader("name=Ana&cpf_number=999"), "application/x-www-form-urlencoded")
		loc := created.Header().Get("Location")
		id := strings.TrimPrefix(loc, "/people/")
		updated := serve(router, http.MethodPost, "/people/"+id, strings.NewReader("name=Bia"), "application/x-www-form-urlencoded")
		if updated.Code != http.StatusSeeOther {
			t.Fatalf("status %d body %s", updated.Code, updated.Body.String())
		}
		page := serve(router, http.MethodGet, loc, nil, "")
		if !strings.Contains(page.Body.String(), "Bia") || strings.Contains(page.Body.String(), "999") {
			t.Fatalf("body %s", page.Body.String())
		}
		deleted := serve(router, http.MethodPost, "/people/"+id+"/delete", nil, "")
		if deleted.Code != http.StatusSeeOther || deleted.Header().Get("Location") != "/people" {
			t.Fatalf("status %d loc %s", deleted.Code, deleted.Header().Get("Location"))
		}
		missing := serve(router, http.MethodGet, loc, nil, "")
		if missing.Code != http.StatusNotFound {
			t.Fatalf("status %d", missing.Code)
		}
	})

	t.Run("should filter the list", func(t *testing.T) {
		router := newRouter(t)
		serve(router, http.MethodPost, "/people", strings.NewReader("name=Ana"), "application/x-www-form-urlencoded")
		serve(router, http.MethodPost, "/people", strings.NewReader("name="+url.QueryEscape("100%")), "application/x-www-form-urlencoded")
		all := serve(router, http.MethodGet, "/people", nil, "")
		if !strings.Contains(all.Body.String(), "Ana") || !strings.Contains(all.Body.String(), "100%") {
			t.Fatalf("body %s", all.Body.String())
		}
		filtered := serve(router, http.MethodGet, "/people?q=%25", nil, "")
		body := filtered.Body.String()
		if !strings.Contains(body, "100%") || strings.Contains(body, ">Ana<") {
			t.Fatalf("body %s", body)
		}
	})

	t.Run("should store a scan and a profile photo", func(t *testing.T) {
		router := newRouter(t)
		body, contentType := multipartBody(t, map[string]string{"name": "Ana"}, map[string][]byte{
			"profile_photo":  []byte("face-bytes"),
			"document_scans": []byte("scan-bytes"),
		})
		rec := serve(router, http.MethodPost, "/people", body, contentType)
		if rec.Code != http.StatusSeeOther {
			t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
		}
		page := serve(router, http.MethodGet, rec.Header().Get("Location"), nil, "")
		html := page.Body.String()
		paths := filePaths(html)
		if len(paths) < 2 {
			t.Fatalf("files %v body %s", paths, html)
		}
		got := map[string]bool{}
		for _, path := range paths {
			rec := serve(router, http.MethodGet, path, nil, "")
			if rec.Code != http.StatusOK {
				t.Fatalf("status %d path %s", rec.Code, path)
			}
			got[rec.Body.String()] = true
		}
		if !got["face-bytes"] || !got["scan-bytes"] {
			t.Fatalf("bodies %v", got)
		}
		escaped := serve(router, http.MethodGet, "/files/../secret", nil, "")
		if escaped.Code != http.StatusNotFound {
			t.Fatalf("status %d", escaped.Code)
		}
	})
}

func newRouter(t *testing.T) http.Handler {
	t.Helper()
	dir := t.TempDir()
	db, err := person.Open(filepath.Join(dir, "id.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := person.Migrate(db); err != nil {
		t.Fatal(err)
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	templates := filepath.Join(filepath.Dir(file), "..", "..", "web", "templates")
	router, err := New(db, filepath.Join(dir, "files"), templates)
	if err != nil {
		t.Fatal(err)
	}
	return router
}

func serve(p_Router http.Handler, p_Method, p_Path string, p_Body io.Reader, p_ContentType string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(p_Method, p_Path, p_Body)
	if p_ContentType != "" {
		req.Header.Set("Content-Type", p_ContentType)
	}
	rec := httptest.NewRecorder()
	p_Router.ServeHTTP(rec, req)
	return rec
}

func multipartBody(t *testing.T, p_Fields map[string]string, p_Files map[string][]byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for key, value := range p_Fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	for key, value := range p_Files {
		name := "face.png"
		if key == "document_scans" {
			name = "scan.png"
		}
		part, err := writer.CreateFormFile(key, name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(value); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf, writer.FormDataContentType()
}

func filePaths(p_Body string) []string {
	var out []string
	rest := p_Body
	for {
		i := strings.Index(rest, "/files/")
		if i < 0 {
			return out
		}
		rest = rest[i:]
		j := strings.IndexAny(rest, "\"' <")
		if j < 0 {
			out = append(out, rest)
			return out
		}
		out = append(out, rest[:j])
		rest = rest[j:]
	}
}

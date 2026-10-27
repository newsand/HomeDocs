package http

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"

	"idcenter/internal/person"
)

type app struct {
	db        *sql.DB
	files     string
	templates map[string]*pongo2.Template
}

func New(p_DB *sql.DB, p_Files, p_Templates string) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	loader, err := pongo2.NewLocalFileSystemLoader(p_Templates)
	if err != nil {
		return nil, err
	}
	set := pongo2.NewSet("pages", loader)
	loaded := map[string]*pongo2.Template{}
	for _, name := range []string{"list.html", "form.html", "show.html", "error.html"} {
		tpl, err := set.FromFile(name)
		if err != nil {
			return nil, err
		}
		loaded[name] = tpl
	}
	application := &app{db: p_DB, files: p_Files, templates: loaded}
	router := gin.New()
	router.Use(requestLog)
	webDir := filepath.Dir(p_Templates)
	router.StaticFile("/tokens.css", filepath.Join(webDir, "tokens.css"))
	router.StaticFile("/app.css", filepath.Join(webDir, "app.css"))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/people")
	})
	router.GET("/people", application.list)
	router.GET("/people/new", application.newForm)
	router.POST("/people", application.create)
	router.GET("/people/:id", application.show)
	router.GET("/people/:id/edit", application.editForm)
	router.POST("/people/:id", application.update)
	router.POST("/people/:id/delete", application.delete)
	router.GET("/files/:name", application.file)
	return router, nil
}

func requestLog(c *gin.Context) {
	id, err := requestID()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("X-Request-ID", id)
	start := time.Now()
	c.Next()
	log.Printf("request_id=%s method=%s path=%s status=%d ms=%d", id, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start).Milliseconds())
}

func requestID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (a *app) render(c *gin.Context, p_Name string, p_Status int, p_Data pongo2.Context) {
	c.Status(p_Status)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := a.templates[p_Name].ExecuteWriter(p_Data, c.Writer); err != nil {
		log.Printf("request_id=%s template=%s failed", c.Writer.Header().Get("X-Request-ID"), p_Name)
	}
}

func (a *app) fail(c *gin.Context) {
	a.render(c, "error.html", http.StatusInternalServerError, pongo2.Context{
		"message": "Não foi possível concluir.",
	})
}

func (a *app) missing(c *gin.Context) {
	a.render(c, "error.html", http.StatusNotFound, pongo2.Context{
		"message": "Pessoa não encontrada.",
	})
}

func (a *app) list(c *gin.Context) {
	people, err := person.SearchPeople(c.Request.Context(), a.db, c.Query("q"))
	if err != nil {
		a.fail(c)
		return
	}
	items := make([]ListItem, len(people))
	for i, item := range people {
		items[i] = ListItem{ID: item.ID, Name: item.Name}
	}
	a.render(c, "list.html", http.StatusOK, pongo2.Context{
		"people": items,
		"query":  c.Query("q"),
	})
}

func (a *app) newForm(c *gin.Context) {
	a.render(c, "form.html", http.StatusOK, pongo2.Context{
		"action": "/people",
		"groups": groupsFrom(person.Person{}),
		"error":  "",
	})
}

func (a *app) editForm(c *gin.Context) {
	item, err := person.ShowPerson(c.Request.Context(), a.db, c.Param("id"))
	if errors.Is(err, person.ErrNotFound) {
		a.missing(c)
		return
	}
	if err != nil {
		a.fail(c)
		return
	}
	a.render(c, "form.html", http.StatusOK, pongo2.Context{
		"action": "/people/" + item.ID,
		"groups": groupsFrom(item),
		"error":  "",
	})
}

func (a *app) create(c *gin.Context) {
	input, err := a.readInput(c)
	if err != nil {
		a.fail(c)
		return
	}
	created, err := person.CreatePerson(c.Request.Context(), a.db, a.files, input)
	if errors.Is(err, person.ErrNameRequired) {
		a.render(c, "form.html", http.StatusUnprocessableEntity, pongo2.Context{
			"action": "/people",
			"groups": groupsFrom(personFromInput(input)),
			"error":  "Informe o nome.",
		})
		return
	}
	if err != nil {
		a.fail(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/people/"+created.ID)
}

func (a *app) update(c *gin.Context) {
	input, err := a.readInput(c)
	if err != nil {
		a.fail(c)
		return
	}
	id := c.Param("id")
	updated, err := person.UpdatePerson(c.Request.Context(), a.db, a.files, id, input)
	if errors.Is(err, person.ErrNameRequired) {
		current, showErr := person.ShowPerson(c.Request.Context(), a.db, id)
		if showErr != nil {
			a.fail(c)
			return
		}
		current.Name = input.Name
		a.render(c, "form.html", http.StatusUnprocessableEntity, pongo2.Context{
			"action": "/people/" + id,
			"groups": groupsFrom(current),
			"error":  "Informe o nome.",
		})
		return
	}
	if errors.Is(err, person.ErrNotFound) {
		a.missing(c)
		return
	}
	if err != nil {
		a.fail(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/people/"+updated.ID)
}

func (a *app) show(c *gin.Context) {
	item, err := person.ShowPerson(c.Request.Context(), a.db, c.Param("id"))
	if errors.Is(err, person.ErrNotFound) {
		a.missing(c)
		return
	}
	if err != nil {
		a.fail(c)
		return
	}
	a.render(c, "show.html", http.StatusOK, pongo2.Context{
		"name":   item.Name,
		"id":     item.ID,
		"photo":  item.ProfilePhoto,
		"scans":  item.DocumentScans,
		"groups": groupsFrom(item),
	})
}

func (a *app) delete(c *gin.Context) {
	err := person.DeletePerson(c.Request.Context(), a.db, a.files, c.Param("id"))
	if errors.Is(err, person.ErrNotFound) {
		a.missing(c)
		return
	}
	if err != nil {
		a.fail(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/people")
}

func (a *app) file(c *gin.Context) {
	full, err := person.FilePath(a.files, c.Param("name"))
	if errors.Is(err, person.ErrNotFound) {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		a.fail(c)
		return
	}
	c.File(full)
}

func (a *app) readInput(c *gin.Context) (person.Input, error) {
	var input person.Input
	for _, group := range catalog() {
		for _, field := range group.fields {
			field.apply(&input, c.PostForm(field.name))
		}
	}
	if err := readFiles(c, &input); err != nil {
		return person.Input{}, err
	}
	return input, nil
}

func readFiles(c *gin.Context, p_Input *person.Input) error {
	form, err := c.MultipartForm()
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) {
			return nil
		}
		return err
	}
	photo, err := readParts(form.File["profile_photo"])
	if err != nil {
		return err
	}
	if len(photo) == 1 {
		p_Input.ProfilePhoto = &photo[0]
	}
	scans, err := readParts(form.File["document_scans"])
	if err != nil {
		return err
	}
	if len(scans) > 0 {
		p_Input.DocumentScans = scans
	}
	return nil
}

func readParts(p_Headers []*multipart.FileHeader) ([]person.Upload, error) {
	out := []person.Upload{}
	for _, header := range p_Headers {
		if header == nil || header.Filename == "" || header.Size == 0 {
			continue
		}
		opened, err := header.Open()
		if err != nil {
			return nil, err
		}
		body, readErr := io.ReadAll(opened)
		closeErr := opened.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		out = append(out, person.Upload{Name: header.Filename, Bytes: body})
	}
	return out, nil
}

func personFromInput(p_Input person.Input) person.Person {
	return person.Person{
		Name:     p_Input.Name,
		CNH:      p_Input.CNH,
		CPF:      p_Input.CPF,
		RG:       p_Input.RG,
		VoterID:  p_Input.VoterID,
		Passport: p_Input.Passport,
	}
}

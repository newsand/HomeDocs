package main

import (
	"idcenter/internal/http"
	"idcenter/internal/person"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	db, err := person.Open("data/idcenter.sqlite")
	if err != nil {
		return err
	}
	if err := person.Migrate(db); err != nil {
		return err
	}
	router, err := http.New(db, "data/files", "web/templates")
	if err != nil {
		return err
	}
	return router.Run(":8080")
}

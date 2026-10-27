package http

import (
	"strings"
	"testing"
)

func TestCatalogKeepsCPFOnce(t *testing.T) {
	form := serve(newRouter(t), "GET", "/people/new", nil, "")
	body := form.Body.String()
	if strings.Count(body, `name="cpf_number"`) != 1 {
		t.Fatalf("cpf inputs %d", strings.Count(body, `name="cpf_number"`))
	}
	for _, gone := range []string{"cnh_cpf", "rg_cpf", "passport_cpf", "cnh_name", "rg_name", "voter_name", "passport_surname"} {
		if strings.Contains(body, `name="`+gone+`"`) {
			t.Fatalf("still present %s", gone)
		}
	}
}

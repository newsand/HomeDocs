# Tests

| suite | purpose | path | command |
|---|---|---|---|
| person | nome, CRUD SQLite, arquivos, busca | internal/person/person_test.go | go test ./internal/person/ |
| http | páginas HTML, 303/422/404, upload, escape | internal/http/server_test.go | go test ./internal/http/ |
